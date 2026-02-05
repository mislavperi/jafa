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
