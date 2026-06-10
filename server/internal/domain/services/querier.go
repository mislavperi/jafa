package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

// ExpenseQuerier is the subset of psql.Queries used by ExpenseService.
// Extracted as an interface to allow test doubles.
type ExpenseQuerier interface {
	GetAllExpenses(ctx context.Context, userID int64) ([]psql.Expense, error)
	GetExpenseById(ctx context.Context, arg psql.GetExpenseByIdParams) (psql.Expense, error)
	CreateExpense(ctx context.Context, arg psql.CreateExpenseParams) (psql.Expense, error)
	UpdateExpense(ctx context.Context, arg psql.UpdateExpenseParams) (psql.Expense, error)
	SoftDeleteExpense(ctx context.Context, arg psql.SoftDeleteExpenseParams) (int64, error)
	GetTotalSpendThisMonth(ctx context.Context, userID int64) (pgtype.Numeric, error)
	GetDailySpend(ctx context.Context, arg psql.GetDailySpendParams) ([]psql.GetDailySpendRow, error)
	GetExpensesByMonth(ctx context.Context, arg psql.GetExpensesByMonthParams) ([]psql.Expense, error)
	GetFirstExpenseDate(ctx context.Context, userID int64) (interface{}, error)
	GetDailySpendForMonth(ctx context.Context, arg psql.GetDailySpendForMonthParams) ([]psql.GetDailySpendForMonthRow, error)
	UpsertTag(ctx context.Context, arg psql.UpsertTagParams) (psql.Tag, error)
	AddTagToExpense(ctx context.Context, arg psql.AddTagToExpenseParams) error
	WithTx(tx pgx.Tx) ExpenseQuerier
}

// psqlQueriesWrapper wraps *psql.Queries so its WithTx satisfies ExpenseQuerier.
type psqlQueriesWrapper struct {
	*psql.Queries
}

func (w *psqlQueriesWrapper) WithTx(tx pgx.Tx) ExpenseQuerier {
	return &psqlQueriesWrapper{w.Queries.WithTx(tx)}
}

func wrapExpenseQueries(q *psql.Queries) ExpenseQuerier {
	return &psqlQueriesWrapper{q}
}

// ReportQuerier is the subset of psql.Queries used by ReportService.
type ReportQuerier interface {
	ListCategories(ctx context.Context) ([]psql.Category, error)
	GetAllExpenses(ctx context.Context, userID int64) ([]psql.Expense, error)
	GetMonthlySpend(ctx context.Context, userID int64) ([]psql.GetMonthlySpendRow, error)
}
