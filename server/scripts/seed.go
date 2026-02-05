package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	"github.com/mislavperi/jafa/server/internal/infrastructure/psql"
)

func main() {
	pool, err := psql.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	const (
		itemCount       = 20
		tagCount        = 10
		itemPriceCount  = 30
		expenseCount    = 50
		expenseTagCount = 40
	)

	ctx := context.Background()

	// Seed Items
	itemIDs := make([]int64, 0, itemCount)
	for i := range itemCount {
		var item models.Item
		if err := faker.FakeData(&item); err != nil {
			log.Fatalf("failed to generate fake item data: %v", err)
		}

		var id int64
		err := pool.QueryRow(
			ctx,
			"INSERT INTO item (name) VALUES ($1) RETURNING id",
			item.Name,
		).Scan(&id)
		if err != nil {
			log.Fatalf("failed to insert item: %v", err)
		}
		itemIDs = append(itemIDs, id)
		fmt.Printf("inserted item %d/%d: name=%q\n", i+1, itemCount, item.Name)
	}
	fmt.Printf("successfully seeded %d items\n", itemCount)

	// Seed Tags
	tagIDs := make([]int64, 0, tagCount)
	for i := range tagCount {
		var tag models.Tag
		if err := faker.FakeData(&tag); err != nil {
			log.Fatalf("failed to generate fake tag data: %v", err)
		}

		var id int64
		err := pool.QueryRow(
			ctx,
			"INSERT INTO tags (name) VALUES ($1) RETURNING id",
			tag.Name,
		).Scan(&id)
		if err != nil {
			log.Fatalf("failed to insert tag: %v", err)
		}
		tagIDs = append(tagIDs, id)
		fmt.Printf("inserted tag %d/%d: name=%q\n", i+1, tagCount, tag.Name)
	}
	fmt.Printf("successfully seeded %d tags\n", tagCount)

	// Seed ItemPrices
	for i := range itemPriceCount {
		var itemPrice models.ItemPrice
		if err := faker.FakeData(&itemPrice); err != nil {
			log.Fatalf("failed to generate fake item price data: %v", err)
		}

		itemID := itemIDs[rand.Intn(len(itemIDs))]
		_, err := pool.Exec(
			ctx,
			"INSERT INTO item_price (item_id, price) VALUES ($1, $2)",
			itemID,
			itemPrice.Price,
		)
		if err != nil {
			log.Fatalf("failed to insert item price: %v", err)
		}
		fmt.Printf("inserted item_price %d/%d: item_id=%d price=%.2f\n", i+1, itemPriceCount, itemID, itemPrice.Price)
	}
	fmt.Printf("successfully seeded %d item prices\n", itemPriceCount)

	// Seed Expenses
	expenseIDs := make([]int64, 0, expenseCount)
	for i := range expenseCount {
		var expense models.Expense
		if err := faker.FakeData(&expense); err != nil {
			log.Fatalf("failed to generate fake expense data: %v", err)
		}

		itemID := itemIDs[rand.Intn(len(itemIDs))]
		var id int64
		err := pool.QueryRow(
			ctx,
			"INSERT INTO expenses (name, amount, cost, item_id) VALUES ($1, $2, $3, $4) RETURNING id",
			expense.Name,
			expense.Amount,
			expense.Cost,
			itemID,
		).Scan(&id)
		if err != nil {
			log.Fatalf("failed to insert expense: %v", err)
		}
		expenseIDs = append(expenseIDs, id)
		fmt.Printf("inserted expense %d/%d: name=%q amount=%.2f cost=%.2f\n", i+1, expenseCount, expense.Name, expense.Amount, expense.Cost)
	}
	fmt.Printf("successfully seeded %d expenses\n", expenseCount)

	// Seed ExpensesTags
	usedPairs := make(map[string]bool)
	inserted := 0
	for inserted < expenseTagCount {
		expenseID := expenseIDs[rand.Intn(len(expenseIDs))]
		tagID := tagIDs[rand.Intn(len(tagIDs))]
		pairKey := fmt.Sprintf("%d-%d", expenseID, tagID)

		if usedPairs[pairKey] {
			continue
		}
		usedPairs[pairKey] = true

		_, err := pool.Exec(
			ctx,
			"INSERT INTO expenses_tags (expense_id, tag_id) VALUES ($1, $2)",
			expenseID,
			tagID,
		)
		if err != nil {
			log.Fatalf("failed to insert expenses_tag: %v", err)
		}
		inserted++
		fmt.Printf("inserted expenses_tag %d/%d: expense_id=%d tag_id=%d\n", inserted, expenseTagCount, expenseID, tagID)
	}
	fmt.Printf("successfully seeded %d expenses_tags\n", expenseTagCount)

	fmt.Println("seeding complete!")
}