package services

import "github.com/mislavperi/jafa/server/internal/domain/models"

// ExpenseServicer is the contract controllers depend on.
// Defined here so tests can substitute a lightweight stub.
type ExpenseServicer interface {
	GetAllExpenses(userID int64) ([]models.Expense, error)
	GetById(userID, id int64) (models.Expense, error)
	GetTotalSpendThisMonth(userID int64) (models.MonthlyTotal, error)
	GetDailySpend(userID int64, months int32) ([]models.DailySpend, error)
	GetFirstExpenseDate(userID int64) (string, error)
	GetDailySpendForMonth(userID int64, year, month int32) ([]models.DailySpend, error)
	GetExpensesByMonth(userID int64, year, month int32) ([]models.Expense, error)
	CreateExpense(input CreateExpenseInput) (models.Expense, error)
	BulkCreateExpenses(userID int64, items []BulkExpenseItem) ([]models.Expense, error)
	UpdateExpense(input UpdateExpenseInput) (models.Expense, error)
	DeleteExpense(userID, id int64) error
}
