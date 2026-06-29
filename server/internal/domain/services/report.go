package services

import (
	"context"
	"math"

	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type ReportService struct {
	Queries        *psql.Queries
	ExpenseMapper  *mappers.ExpenseMapper
	CategoryMapper *mappers.CategoryMapper
}

func NewReportService(pool psql.Pool) *ReportService {
	return &ReportService{
		Queries:        psql.New(pool),
		ExpenseMapper:  mappers.NewExpenseMapper(),
		CategoryMapper: mappers.NewCategoryMapper(),
	}
}

// CategoryBreakdown buckets the user's expenses into categories (by keyword
// match) and returns the spend, remaining budget and percentage used for each.
func (rs *ReportService) CategoryBreakdown(ctx context.Context, userID int64) ([]models.CategoryBreakdown, error) {
	rawCategories, err := rs.Queries.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	categories, err := rs.CategoryMapper.MapManyToDomain(rawCategories)
	if err != nil {
		return nil, err
	}

	rawExpenses, err := rs.Queries.GetAllExpenses(ctx, userID)
	if err != nil {
		return nil, err
	}
	expenses, err := rs.ExpenseMapper.MapManyToDomain(rawExpenses)
	if err != nil {
		return nil, err
	}

	totals := make(map[string]float32, len(categories))
	for _, expense := range expenses {
		name := Categorize(expense.Name, categories)
		totals[name] += expense.Cost
	}

	breakdown := make([]models.CategoryBreakdown, 0, len(categories))
	for _, category := range categories {
		spent := totals[category.Name]
		pct := 0
		if category.Budget > 0 {
			pct = int(math.Min(100, math.Round(float64(spent/category.Budget*100))))
		}
		breakdown = append(breakdown, models.CategoryBreakdown{
			Name:      category.Name,
			Icon:      category.Icon,
			Color:     category.Color,
			Budget:    category.Budget,
			Spent:     spent,
			Remaining: category.Budget - spent,
			Pct:       pct,
		})
	}
	return breakdown, nil
}

// MonthlySpend returns the total spend per calendar month, oldest first.
func (rs *ReportService) MonthlySpend(ctx context.Context, userID int64) ([]models.MonthlySpend, error) {
	rows, err := rs.Queries.GetMonthlySpend(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]models.MonthlySpend, 0, len(rows))
	for _, row := range rows {
		f, err := row.Total.Float64Value()
		if err != nil {
			return nil, err
		}
		if !f.Valid {
			continue
		}
		result = append(result, models.MonthlySpend{
			Month: row.Month,
			Total: float32(f.Float64),
		})
	}
	return result, nil
}
