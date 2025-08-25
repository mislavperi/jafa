package services

import "github.com/mislavperi/jafa/internal/domain/models"

type ExpenseService struct{}

func NewExpenseService() *ExpenseService {
	return &ExpenseService{}
}

func (es *ExpenseService) GetById() (models.Expense, error) {
	return models.Expense{}, nil
}

func (es *ExpenseService) GetByType(expenseType models.ExpenseType) ([]models.Expense, error) {
	return []models.Expense{}, nil
}
