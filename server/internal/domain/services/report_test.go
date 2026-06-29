package services

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/pashagolub/pgxmock/v4"
)

// newReportService builds a ReportService backed by the mocked pool. The caller
// owns the mock and sets query expectations before exercising the service.
func newReportService(t *testing.T) (*ReportService, pgxmock.PgxPoolIface) {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool: %v", err)
	}
	t.Cleanup(mock.Close)
	svc := &ReportService{
		Queries:        psql.New(mock),
		ExpenseMapper:  mappers.NewExpenseMapper(),
		CategoryMapper: mappers.NewCategoryMapper(),
	}
	return svc, mock
}

// categoryColumns mirrors the SELECT column order SQLC scans for a Category row.
// (expenseColumns lives in setup_test.go, shared with the expense tests.)
var categoryColumns = []string{"id", "name", "icon", "color", "budget", "keywords", "sort_order"}

// expenseRow returns one Expense row's values in scan order. Amount is set to a
// sentinel distinct from cost so the breakdown's reliance on cost (not amount)
// is pinned: summing amount instead would change the result.
func expenseRow(t *testing.T, name, cost string) []any {
	t.Helper()
	return []any{
		int64(0), name, scanNumeric(t, "0.001"), scanNumeric(t, cost),
		pgtype.Int8{}, false, pgtype.Timestamptz{}, pgtype.Timestamptz{},
		pgtype.Text{}, pgtype.Int4{}, pgtype.Date{}, int64(1), pgtype.Int4{}, "expense",
	}
}

// ---- CategoryBreakdown ----

func TestCategoryBreakdown_SpendAndPct(t *testing.T) {
	svc, mock := newReportService(t)
	mock.ExpectQuery("ListCategories").WillReturnRows(
		pgxmock.NewRows(categoryColumns).
			AddRow(int64(1), "Groceries", "", "", scanNumeric(t, "200.000"), []string{"food", "market"}, int32(0)).
			AddRow(int64(2), "Other", "", "", scanNumeric(t, "0"), []string(nil), int32(1)),
	)
	mock.ExpectQuery("GetAllExpenses").WithArgs(int64(1)).WillReturnRows(
		pgxmock.NewRows(expenseColumns).
			AddRow(expenseRow(t, "Whole Foods Market", "100.000")...).
			AddRow(expenseRow(t, "FOOD truck", "50.000")...),
	)

	breakdown, err := svc.CategoryBreakdown(context.Background(), 1)
	if err != nil {
		t.Fatalf("CategoryBreakdown: %v", err)
	}

	var groceries *models.CategoryBreakdown
	for i := range breakdown {
		if breakdown[i].Name == "Groceries" {
			groceries = &breakdown[i]
		}
	}
	if groceries == nil {
		t.Fatal("Groceries not found in breakdown")
	}
	if groceries.Spent < 149.9 || groceries.Spent > 150.1 {
		t.Errorf("Groceries.Spent = %v, want ~150", groceries.Spent)
	}
	if groceries.Pct != 75 {
		t.Errorf("Groceries.Pct = %d, want 75", groceries.Pct)
	}
	if groceries.Remaining < 49.9 || groceries.Remaining > 50.1 {
		t.Errorf("Groceries.Remaining = %v, want ~50", groceries.Remaining)
	}
}

func TestCategoryBreakdown_PctCappedAt100(t *testing.T) {
	svc, mock := newReportService(t)
	mock.ExpectQuery("ListCategories").WillReturnRows(
		pgxmock.NewRows(categoryColumns).
			AddRow(int64(1), "Dining", "", "", scanNumeric(t, "100.000"), []string{"restaurant"}, int32(0)),
	)
	mock.ExpectQuery("GetAllExpenses").WithArgs(int64(1)).WillReturnRows(
		pgxmock.NewRows(expenseColumns).
			AddRow(expenseRow(t, "Nice restaurant", "200.000")...), // overspent
	)

	breakdown, err := svc.CategoryBreakdown(context.Background(), 1)
	if err != nil {
		t.Fatalf("CategoryBreakdown: %v", err)
	}
	if len(breakdown) == 0 {
		t.Fatal("empty breakdown")
	}
	if breakdown[0].Pct != 100 {
		t.Errorf("Pct = %d, want 100 (capped)", breakdown[0].Pct)
	}
}

func TestCategoryBreakdown_ZeroBudget(t *testing.T) {
	svc, mock := newReportService(t)
	mock.ExpectQuery("ListCategories").WillReturnRows(
		pgxmock.NewRows(categoryColumns).
			AddRow(int64(1), "Other", "", "", scanNumeric(t, "0"), []string(nil), int32(0)),
	)
	mock.ExpectQuery("GetAllExpenses").WithArgs(int64(1)).WillReturnRows(pgxmock.NewRows(expenseColumns))

	breakdown, err := svc.CategoryBreakdown(context.Background(), 1)
	if err != nil {
		t.Fatalf("CategoryBreakdown: %v", err)
	}
	if len(breakdown) == 0 {
		t.Fatal("empty breakdown")
	}
	if breakdown[0].Pct != 0 {
		t.Errorf("Pct = %d, want 0 when budget is zero", breakdown[0].Pct)
	}
}

// ---- MonthlySpend ----

func TestMonthlySpend_SkipsNullTotals(t *testing.T) {
	svc, mock := newReportService(t)
	mock.ExpectQuery("GetMonthlySpend").WithArgs(int64(1)).WillReturnRows(
		pgxmock.NewRows([]string{"month", "total"}).
			AddRow("2024-01", scanNumeric(t, "150.000")).
			AddRow("2024-02", pgtype.Numeric{}). // NULL
			AddRow("2024-03", scanNumeric(t, "150.000")),
	)

	result, err := svc.MonthlySpend(context.Background(), 1)
	if err != nil {
		t.Fatalf("MonthlySpend: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("len = %d, want 2 (NULL row skipped)", len(result))
	}
	if result[0].Month != "2024-01" || result[1].Month != "2024-03" {
		t.Errorf("months = [%s, %s]", result[0].Month, result[1].Month)
	}
}
