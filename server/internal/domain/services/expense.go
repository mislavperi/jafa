package services

import (
	"context"

	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type ExpenseService struct {
	Queries *psql.Queries
	Mapper  *mappers.ExpenseMapper
}

func NewExpenseService(queries *psql.Queries) *ExpenseService {
	return &ExpenseService{
		Queries: queries,
		Mapper:  mappers.NewExpenseMapper(),
	}
}

func (es *ExpenseService) GetAllExpenses() ([]models.Expense, error) {
	expenses, err := es.Queries.GetAllExpenses(context.Background())
	if err != nil {
		return nil, err
	}
	mappedExpenses, err := es.Mapper.MapManyToDomain(expenses)
	if err != nil {
		return nil, err
	}
	return mappedExpenses, nil
}

func (es *ExpenseService) GetTotalSpendThisMonth() (models.MonthlyTotal, error) {
	total, err := es.Queries.GetTotalSpendThisMonth(context.Background())
	if err != nil {
		return models.MonthlyTotal{}, err
	}
	f, err := total.Float64Value()
	if err != nil || !f.Valid {
		return models.MonthlyTotal{}, err
	}
	return models.MonthlyTotal{
		Total: float32(f.Float64),
	}, nil
}

func (es *ExpenseService) GetDailySpend(months int32) ([]models.DailySpend, error) {
	rows, err := es.Queries.GetDailySpend(context.Background(), months)
	if err != nil {
		return nil, err
	}
	result := make([]models.DailySpend, 0, len(rows))
	for _, row := range rows {
		f, err := row.Total.Float64Value()
		if err != nil || !f.Valid {
			return nil, err
		}
		result = append(result, models.DailySpend{
			Day:   row.Day.Time.Format("2006-01-02"),
			Total: float32(f.Float64),
		})
	}
	return result, nil
}

func (es *ExpenseService) GetExpensesByMonth(year, month int32) ([]models.Expense, error) {
	expenses, err := es.Queries.GetExpensesByMonth(context.Background(), psql.GetExpensesByMonthParams{
		Year:  year,
		Month: month,
	})
	if err != nil {
		return nil, err
	}
	return es.Mapper.MapManyToDomain(expenses)
}

func (es *ExpenseService) GetFirstExpenseDate() (string, error) {
	result, err := es.Queries.GetFirstExpenseDate(context.Background())
	if err != nil {
		return "", err
	}
	s, ok := result.(string)
	if !ok {
		return "", nil
	}
	return s, nil
}

func (es *ExpenseService) GetDailySpendForMonth(year, month int32) ([]models.DailySpend, error) {
	rows, err := es.Queries.GetDailySpendForMonth(context.Background(), psql.GetDailySpendForMonthParams{
		Year:  year,
		Month: month,
	})
	if err != nil {
		return nil, err
	}
	result := make([]models.DailySpend, 0, len(rows))
	for _, row := range rows {
		f, err := row.Total.Float64Value()
		if err != nil || !f.Valid {
			return nil, err
		}
		result = append(result, models.DailySpend{
			Day:   row.Day.Time.Format("2006-01-02"),
			Total: float32(f.Float64),
		})
	}
	return result, nil
}

func (es *ExpenseService) GetById(id int64) (models.Expense, error) {
	expense, err := es.Queries.GetExpenseById(context.Background(), id)
	if err != nil {
		return models.Expense{}, err
	}
	mappedExpense, err := es.Mapper.MapToDomain(expense)
	if err != nil {
		return models.Expense{}, err
	}
	return mappedExpense, nil
}
