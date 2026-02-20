package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mislavperi/jafa/server/internal/infrastructure/psql"
)

func main() {
	pool, err := psql.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	// Realistic financial categories
	categories := []string{
		"Groceries",
		"Dining Out",
		"Transport",
		"Utilities",
		"Entertainment",
		"Health",
		"Housing",
	}

	// Realistic items with base prices
	type itemDef struct {
		name  string
		price float64
	}
	items := []itemDef{
		{name: "Supermarket Shop", price: 85.00},
		{name: "Bread & Bakery", price: 6.50},
		{name: "Eggs & Dairy", price: 9.80},
		{name: "Coffee", price: 4.50},
		{name: "Restaurant Dinner", price: 42.00},
		{name: "Takeaway", price: 18.00},
		{name: "Bus Pass", price: 50.00},
		{name: "Petrol", price: 65.00},
		{name: "Electricity Bill", price: 90.00},
		{name: "Internet Bill", price: 40.00},
		{name: "Water Bill", price: 30.00},
		{name: "Netflix", price: 15.99},
		{name: "Spotify", price: 9.99},
		{name: "Cinema Tickets", price: 24.00},
		{name: "Gym Membership", price: 35.00},
		{name: "Pharmacy", price: 22.00},
		{name: "Rent", price: 1200.00},
		{name: "Home Insurance", price: 55.00},
	}

	// Seed items and their prices
	itemIDs := make([]int64, 0, len(items))
	for i, it := range items {
		var id int64
		err := pool.QueryRow(
			ctx,
			"INSERT INTO item (name) VALUES ($1) RETURNING id",
			it.name,
		).Scan(&id)
		if err != nil {
			log.Fatalf("failed to insert item: %v", err)
		}
		itemIDs = append(itemIDs, id)

		_, err = pool.Exec(
			ctx,
			"INSERT INTO item_price (item_id, price) VALUES ($1, $2)",
			id,
			it.price,
		)
		if err != nil {
			log.Fatalf("failed to insert item price: %v", err)
		}
		fmt.Printf("inserted item %d/%d: name=%q price=%.2f\n", i+1, len(items), it.name, it.price)
	}
	fmt.Printf("successfully seeded %d items\n", len(items))

	// Seed tags (categories)
	tagIDs := make([]int64, 0, len(categories))
	for i, cat := range categories {
		var id int64
		err := pool.QueryRow(
			ctx,
			"INSERT INTO tags (name) VALUES ($1) RETURNING id",
			cat,
		).Scan(&id)
		if err != nil {
			log.Fatalf("failed to insert tag: %v", err)
		}
		tagIDs = append(tagIDs, id)
		fmt.Printf("inserted tag %d/%d: name=%q\n", i+1, len(categories), cat)
	}
	fmt.Printf("successfully seeded %d tags\n", len(categories))

	// Tag index constants for readability
	const (
		tagGroceries    = 0
		tagDiningOut    = 1
		tagTransport    = 2
		tagUtilities    = 3
		tagEntertainment = 4
		tagHealth       = 5
		tagHousing      = 6
	)

	// Expense templates: monthly recurring and occasional expenses
	// itemIdx maps to the items slice above
	type expenseTemplate struct {
		name    string
		amount  float32
		itemIdx int
		tagIdx  int
		// everyMonth = true means this appears each month; false = random chance
		everyMonth bool
		// chance is the probability (0.0–1.0) of appearing in a given month if not everyMonth
		chance float32
	}

	templates := []expenseTemplate{
		// Groceries
		{name: "Weekly supermarket shop", amount: 1, itemIdx: 0, tagIdx: tagGroceries, everyMonth: true},
		{name: "Bread and bakery", amount: 2, itemIdx: 1, tagIdx: tagGroceries, everyMonth: true},
		{name: "Eggs and dairy", amount: 1, itemIdx: 2, tagIdx: tagGroceries, everyMonth: true},
		// Dining Out
		{name: "Morning coffee", amount: 8, itemIdx: 3, tagIdx: tagDiningOut, everyMonth: true},
		{name: "Restaurant dinner", amount: 1, itemIdx: 4, tagIdx: tagDiningOut, everyMonth: false, chance: 0.7},
		{name: "Takeaway", amount: 2, itemIdx: 5, tagIdx: tagDiningOut, everyMonth: true},
		// Transport
		{name: "Monthly bus pass", amount: 1, itemIdx: 6, tagIdx: tagTransport, everyMonth: true},
		{name: "Petrol top-up", amount: 1, itemIdx: 7, tagIdx: tagTransport, everyMonth: false, chance: 0.6},
		// Utilities
		{name: "Electricity bill", amount: 1, itemIdx: 8, tagIdx: tagUtilities, everyMonth: true},
		{name: "Internet bill", amount: 1, itemIdx: 9, tagIdx: tagUtilities, everyMonth: true},
		{name: "Water bill", amount: 1, itemIdx: 10, tagIdx: tagUtilities, everyMonth: true},
		// Entertainment
		{name: "Netflix subscription", amount: 1, itemIdx: 11, tagIdx: tagEntertainment, everyMonth: true},
		{name: "Spotify subscription", amount: 1, itemIdx: 12, tagIdx: tagEntertainment, everyMonth: true},
		{name: "Cinema", amount: 2, itemIdx: 13, tagIdx: tagEntertainment, everyMonth: false, chance: 0.5},
		// Health
		{name: "Gym membership", amount: 1, itemIdx: 14, tagIdx: tagHealth, everyMonth: true},
		{name: "Pharmacy", amount: 1, itemIdx: 15, tagIdx: tagHealth, everyMonth: false, chance: 0.4},
		// Housing
		{name: "Rent", amount: 1, itemIdx: 16, tagIdx: tagHousing, everyMonth: true},
		{name: "Home insurance", amount: 1, itemIdx: 17, tagIdx: tagHousing, everyMonth: false, chance: 0.25},
	}

	// Seed expenses spread over the last 12 months
	now := time.Now()
	expenseIDs := make([]int64, 0)
	expenseTagIndices := make([]int, 0)

	for monthsAgo := 11; monthsAgo >= 0; monthsAgo-- {
		year, month, _ := now.Date()
		targetMonth := time.Month(int(month) - monthsAgo)
		targetYear := year
		for targetMonth <= 0 {
			targetMonth += 12
			targetYear--
		}
		monthStart := time.Date(targetYear, targetMonth, 1, 0, 0, 0, 0, time.UTC)
		daysInMonth := monthStart.AddDate(0, 1, 0).Add(-time.Second).Day()

		for _, tmpl := range templates {
			if !tmpl.everyMonth && rand.Float32() > tmpl.chance {
				continue
			}

			// Randomize the day within the month
			day := rand.Intn(daysInMonth) + 1
			createdAt := time.Date(targetYear, targetMonth, day, rand.Intn(14)+8, rand.Intn(60), 0, 0, time.UTC)

			// Apply ±10% variation to item price as the actual cost
			basePrice := items[tmpl.itemIdx].price
			variation := (rand.Float64() - 0.5) * 0.2 * basePrice
			cost := basePrice + variation

			var id int64
			err := pool.QueryRow(
				ctx,
				"INSERT INTO expenses (name, amount, cost, item_id, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
				tmpl.name,
				tmpl.amount,
				cost,
				itemIDs[tmpl.itemIdx],
				createdAt,
			).Scan(&id)
			if err != nil {
				log.Fatalf("failed to insert expense: %v", err)
			}
			expenseIDs = append(expenseIDs, id)
			expenseTagIndices = append(expenseTagIndices, tmpl.tagIdx)
			fmt.Printf("inserted expense: name=%q cost=%.2f date=%s\n", tmpl.name, cost, createdAt.Format("2006-01-02"))
		}
	}
	fmt.Printf("successfully seeded %d expenses\n", len(expenseIDs))

	// Seed expenses_tags
	for i, expenseID := range expenseIDs {
		tagID := tagIDs[expenseTagIndices[i]]
		_, err := pool.Exec(
			ctx,
			"INSERT INTO expenses_tags (expense_id, tag_id) VALUES ($1, $2)",
			expenseID,
			tagID,
		)
		if err != nil {
			log.Fatalf("failed to insert expenses_tag: %v", err)
		}
	}
	fmt.Printf("successfully seeded %d expenses_tags\n", len(expenseIDs))

	fmt.Println("seeding complete!")
}
