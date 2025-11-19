package services

import (
	"context"

	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type ExpenseService struct {
	Queries *psql.Queries
}

func NewExpenseService(queries *psql.Queries) *ExpenseService {
	return &ExpenseService{
		Queries: queries,
	}
}

func (es *ExpenseService) GetById(id int64) (models.Expense, error) {
	es.Queries.GetExpenseById(context.Background(), id)
	return models.Expense{}, nil
}

func (es *ExpenseService) GetByType(expenseType models.ExpenseType) ([]models.Expense, error) {
	return []models.Expense{}, nil
}
