package services

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type stubReportQuerier struct {
	listCategoriesFn  func(ctx context.Context) ([]psql.Category, error)
	getAllExpensesFn   func(ctx context.Context, userID int64) ([]psql.Expense, error)
	getMonthlySpendFn func(ctx context.Context, userID int64) ([]psql.GetMonthlySpendRow, error)
}

func (s *stubReportQuerier) ListCategories(ctx context.Context) ([]psql.Category, error) {
	return s.listCategoriesFn(ctx)
}
func (s *stubReportQuerier) GetAllExpenses(ctx context.Context, userID int64) ([]psql.Expense, error) {
	return s.getAllExpensesFn(ctx, userID)
}
func (s *stubReportQuerier) GetMonthlySpend(ctx context.Context, userID int64) ([]psql.GetMonthlySpendRow, error) {
	return s.getMonthlySpendFn(ctx, userID)
}

func newReportServiceWithStub(q ReportQuerier) *ReportService {
	return &ReportService{
		Queries:        q,
		ExpenseMapper:  mappers.NewExpenseMapper(),
		CategoryMapper: mappers.NewCategoryMapper(),
	}
}

func scanNumeric(t *testing.T, s string) pgtype.Numeric {
	t.Helper()
	var n pgtype.Numeric
	if err := n.Scan(s); err != nil {
		t.Fatalf("Scan(%q): %v", s, err)
	}
	return n
}

func makeExpense(t *testing.T, name string, amount string) psql.Expense {
	t.Helper()
	now := pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
	return psql.Expense{
		Name:      name,
		Amount:    scanNumeric(t, amount),
		Cost:      scanNumeric(t, amount),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ---- CategoryBreakdown ----

func TestCategoryBreakdown_SpendAndPct(t *testing.T) {
	stub := &stubReportQuerier{
		listCategoriesFn: func(_ context.Context) ([]psql.Category, error) {
			return []psql.Category{
				{Name: "Groceries", Budget: scanNumeric(t, "200.000"), Keywords: []string{"food", "market"}},
				{Name: "Other", Budget: scanNumeric(t, "0"), Keywords: nil},
			}, nil
		},
		getAllExpensesFn: func(_ context.Context, _ int64) ([]psql.Expense, error) {
			return []psql.Expense{
				makeExpense(t, "Whole Foods Market", "100.000"),
				makeExpense(t, "FOOD truck", "50.000"),
			}, nil
		},
	}
	svc := newReportServiceWithStub(stub)
	breakdown, err := svc.CategoryBreakdown(1)
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
	stub := &stubReportQuerier{
		listCategoriesFn: func(_ context.Context) ([]psql.Category, error) {
			return []psql.Category{
				{Name: "Dining", Budget: scanNumeric(t, "100.000"), Keywords: []string{"restaurant"}},
			}, nil
		},
		getAllExpensesFn: func(_ context.Context, _ int64) ([]psql.Expense, error) {
			return []psql.Expense{
				makeExpense(t, "Nice restaurant", "200.000"), // overspent
			}, nil
		},
	}
	svc := newReportServiceWithStub(stub)
	breakdown, err := svc.CategoryBreakdown(1)
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
	stub := &stubReportQuerier{
		listCategoriesFn: func(_ context.Context) ([]psql.Category, error) {
			return []psql.Category{
				{Name: "Other", Budget: scanNumeric(t, "0"), Keywords: nil},
			}, nil
		},
		getAllExpensesFn: func(_ context.Context, _ int64) ([]psql.Expense, error) {
			return nil, nil
		},
	}
	svc := newReportServiceWithStub(stub)
	breakdown, err := svc.CategoryBreakdown(1)
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
	var validTotal pgtype.Numeric
	validTotal.Scan("150.000")

	stub := &stubReportQuerier{
		getMonthlySpendFn: func(_ context.Context, _ int64) ([]psql.GetMonthlySpendRow, error) {
			return []psql.GetMonthlySpendRow{
				{Month: "2024-01", Total: validTotal},
				{Month: "2024-02", Total: pgtype.Numeric{}}, // NULL
				{Month: "2024-03", Total: validTotal},
			}, nil
		},
	}
	svc := newReportServiceWithStub(stub)
	result, err := svc.MonthlySpend(1)
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
