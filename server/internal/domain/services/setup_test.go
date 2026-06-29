package services

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/pashagolub/pgxmock/v4"
)

// Shared test scaffolding for the services package.

// newExpenseService builds an ExpenseService backed by the mocked pool. The
// caller owns the mock and sets query expectations before exercising it.
func newExpenseService(t *testing.T) (*ExpenseService, pgxmock.PgxPoolIface) {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool: %v", err)
	}
	t.Cleanup(mock.Close)
	return &ExpenseService{Queries: psql.New(mock), Pool: mock, Mapper: mappers.NewExpenseMapper()}, mock
}

func intPtr(v int) *int { return &v }

// scanNumeric builds a pgtype.Numeric from a decimal string, like the DB returns.
func scanNumeric(t *testing.T, s string) pgtype.Numeric {
	t.Helper()
	var n pgtype.Numeric
	if err := n.Scan(s); err != nil {
		t.Fatalf("Scan(%q): %v", s, err)
	}
	return n
}

// expenseColumns mirrors the SELECT column order SQLC scans for an Expense row.
// Shared by the expense and report service tests.
var expenseColumns = []string{
	"id", "name", "amount", "cost", "item_id", "is_deleted", "created_at", "updated_at",
	"recurrence_interval", "recurrence_day", "recurrence_start_date", "user_id", "installment_count", "kind",
}

// returnedExpense is one CreateExpense/Update/GetById RETURNING row in scan
// order, with only the fields the mapper reads (amount, cost, installment, kind)
// set meaningfully.
func returnedExpense(t *testing.T, id int64, name, amount, cost string, installment pgtype.Int4, kind string) []any {
	t.Helper()
	return []any{
		id, name, scanNumeric(t, amount), scanNumeric(t, cost),
		pgtype.Int8{}, false, pgtype.Timestamptz{}, pgtype.Timestamptz{},
		pgtype.Text{}, pgtype.Int4{}, pgtype.Date{}, int64(1), installment, kind,
	}
}
