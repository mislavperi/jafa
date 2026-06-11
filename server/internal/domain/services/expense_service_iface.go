package services

import (
	"context"

	"github.com/mislavperi/jafa/server/internal/domain/models"
)

// ExpenseServicer is the contract controllers depend on.
// Defined here so tests can substitute a lightweight stub.
type ExpenseServicer interface {
	GetAllExpenses(ctx context.Context, userID int64) ([]models.Expense, error)
	GetById(ctx context.Context, userID, id int64) (models.Expense, error)
	GetTotalSpendThisMonth(ctx context.Context, userID int64) (models.MonthlyTotal, error)
	GetDailySpend(ctx context.Context, userID int64, months int32) ([]models.DailySpend, error)
	GetFirstExpenseDate(ctx context.Context, userID int64) (string, error)
	GetDailySpendForMonth(ctx context.Context, userID int64, year, month int32) ([]models.DailySpend, error)
	GetExpensesByMonth(ctx context.Context, userID int64, year, month int32) ([]models.Expense, error)
	CreateExpense(ctx context.Context, input CreateExpenseInput) (models.Expense, error)
	BulkCreateExpenses(ctx context.Context, userID int64, items []BulkExpenseItem) ([]models.Expense, error)
	UpdateExpense(ctx context.Context, input UpdateExpenseInput) (models.Expense, error)
	DeleteExpense(ctx context.Context, userID, id int64) error
}
