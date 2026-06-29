package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/utils"
)

// Aggregate/dashboard queries for expenses: monthly totals and daily spend.

// numericToMonthlyTotal converts a pgtype.Numeric sum into a MonthlyTotal,
// treating a NULL/invalid value as a zero total.
func numericToMonthlyTotal(n pgtype.Numeric) (models.MonthlyTotal, error) {
	total, err := utils.NumericToFloat(n)
	if err != nil {
		return models.MonthlyTotal{}, err
	}
	return models.MonthlyTotal{Total: total}, nil
}

// dailySpendFromRow converts a (day, total) pair from a daily-spend query row
// into a DailySpend. A NULL/invalid total yields a zero-value DailySpend.
func dailySpendFromRow(day pgtype.Date, total pgtype.Numeric) (models.DailySpend, error) {
	f, err := utils.NumericToFloat(total)
	if err != nil {
		return models.DailySpend{}, err
	}
	return models.DailySpend{
		Day:   day.Time.Format("2006-01-02"),
		Total: f,
	}, nil
}

func (es *ExpenseService) GetTotalSpendThisMonth(ctx context.Context, userID int64) (models.MonthlyTotal, error) {
	total, err := es.Queries.GetTotalSpendThisMonth(ctx, userID)
	if err != nil {
		return models.MonthlyTotal{}, err
	}
	return numericToMonthlyTotal(total)
}

func (es *ExpenseService) GetTotalIncomeThisMonth(ctx context.Context, userID int64) (models.MonthlyTotal, error) {
	total, err := es.Queries.GetTotalIncomeThisMonth(ctx, userID)
	if err != nil {
		return models.MonthlyTotal{}, err
	}
	return numericToMonthlyTotal(total)
}

func (es *ExpenseService) GetDailySpend(ctx context.Context, userID int64, months int32) ([]models.DailySpend, error) {
	rows, err := es.Queries.GetDailySpend(ctx, psql.GetDailySpendParams{
		UserID: userID,
		Months: months,
	})
	if err != nil {
		return nil, err
	}
	result := make([]models.DailySpend, 0, len(rows))
	for _, row := range rows {
		spend, err := dailySpendFromRow(row.Day, row.Total)
		if err != nil {
			return nil, err
		}
		result = append(result, spend)
	}
	return result, nil
}

func (es *ExpenseService) GetDailySpendForMonth(ctx context.Context, userID int64, year, month int32) ([]models.DailySpend, error) {
	rows, err := es.Queries.GetDailySpendForMonth(ctx, psql.GetDailySpendForMonthParams{
		UserID: userID,
		Year:   year,
		Month:  month,
	})
	if err != nil {
		return nil, err
	}
	result := make([]models.DailySpend, 0, len(rows))
	for _, row := range rows {
		spend, err := dailySpendFromRow(row.Day, row.Total)
		if err != nil {
			return nil, err
		}
		result = append(result, spend)
	}
	return result, nil
}
