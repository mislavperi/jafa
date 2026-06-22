package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

// Ensure stubExpenseQuerier satisfies the interface at compile time.
var _ ExpenseQuerier = (*stubExpenseQuerier)(nil)

// stubExpenseQuerier satisfies ExpenseQuerier with configurable per-method stubs.
type stubExpenseQuerier struct {
	softDeleteFn       func(ctx context.Context, arg psql.SoftDeleteExpenseParams) (int64, error)
	updateExpenseFn    func(ctx context.Context, arg psql.UpdateExpenseParams) (psql.Expense, error)
	getExpenseByIdFn   func(ctx context.Context, arg psql.GetExpenseByIdParams) (psql.Expense, error)
	getAllExpensesFn   func(ctx context.Context, userID int64) ([]psql.Expense, error)
	getAllEntriesFn    func(ctx context.Context, userID int64) ([]psql.Expense, error)
	createExpenseFn    func(ctx context.Context, arg psql.CreateExpenseParams) (psql.Expense, error)
	getTotalFn         func(ctx context.Context, userID int64) (pgtype.Numeric, error)
	getIncomeFn        func(ctx context.Context, userID int64) (pgtype.Numeric, error)
	getDailySpendFn    func(ctx context.Context, arg psql.GetDailySpendParams) ([]psql.GetDailySpendRow, error)
	getByMonthFn       func(ctx context.Context, arg psql.GetExpensesByMonthParams) ([]psql.Expense, error)
	getFirstDateFn     func(ctx context.Context, userID int64) (any, error)
	getDailyForMonthFn func(ctx context.Context, arg psql.GetDailySpendForMonthParams) ([]psql.GetDailySpendForMonthRow, error)
	upsertTagFn        func(ctx context.Context, arg psql.UpsertTagParams) (psql.Tag, error)
	addTagFn           func(ctx context.Context, arg psql.AddTagToExpenseParams) error
}

func (s *stubExpenseQuerier) SoftDeleteExpense(ctx context.Context, arg psql.SoftDeleteExpenseParams) (int64, error) {
	return s.softDeleteFn(ctx, arg)
}
func (s *stubExpenseQuerier) UpdateExpense(ctx context.Context, arg psql.UpdateExpenseParams) (psql.Expense, error) {
	return s.updateExpenseFn(ctx, arg)
}
func (s *stubExpenseQuerier) GetExpenseById(ctx context.Context, arg psql.GetExpenseByIdParams) (psql.Expense, error) {
	return s.getExpenseByIdFn(ctx, arg)
}
func (s *stubExpenseQuerier) GetAllExpenses(ctx context.Context, userID int64) ([]psql.Expense, error) {
	return s.getAllExpensesFn(ctx, userID)
}
func (s *stubExpenseQuerier) CreateExpense(ctx context.Context, arg psql.CreateExpenseParams) (psql.Expense, error) {
	return s.createExpenseFn(ctx, arg)
}
func (s *stubExpenseQuerier) GetAllEntries(ctx context.Context, userID int64) ([]psql.Expense, error) {
	if s.getAllEntriesFn != nil {
		return s.getAllEntriesFn(ctx, userID)
	}
	return nil, nil
}
func (s *stubExpenseQuerier) GetTotalSpendThisMonth(ctx context.Context, userID int64) (pgtype.Numeric, error) {
	return s.getTotalFn(ctx, userID)
}
func (s *stubExpenseQuerier) GetTotalIncomeThisMonth(ctx context.Context, userID int64) (pgtype.Numeric, error) {
	if s.getIncomeFn != nil {
		return s.getIncomeFn(ctx, userID)
	}
	return pgtype.Numeric{}, nil
}
func (s *stubExpenseQuerier) GetDailySpend(ctx context.Context, arg psql.GetDailySpendParams) ([]psql.GetDailySpendRow, error) {
	return s.getDailySpendFn(ctx, arg)
}
func (s *stubExpenseQuerier) GetExpensesByMonth(ctx context.Context, arg psql.GetExpensesByMonthParams) ([]psql.Expense, error) {
	return s.getByMonthFn(ctx, arg)
}
func (s *stubExpenseQuerier) GetFirstExpenseDate(ctx context.Context, userID int64) (interface{}, error) {
	return s.getFirstDateFn(ctx, userID)
}
func (s *stubExpenseQuerier) GetDailySpendForMonth(ctx context.Context, arg psql.GetDailySpendForMonthParams) ([]psql.GetDailySpendForMonthRow, error) {
	return s.getDailyForMonthFn(ctx, arg)
}
func (s *stubExpenseQuerier) UpsertTag(ctx context.Context, arg psql.UpsertTagParams) (psql.Tag, error) {
	return s.upsertTagFn(ctx, arg)
}
func (s *stubExpenseQuerier) AddTagToExpense(ctx context.Context, arg psql.AddTagToExpenseParams) error {
	return s.addTagFn(ctx, arg)
}
func (s *stubExpenseQuerier) WithTx(_ pgx.Tx) ExpenseQuerier { return s }

func newExpenseServiceWithStub(q ExpenseQuerier) *ExpenseService {
	return &ExpenseService{Queries: q, Mapper: mappers.NewExpenseMapper()}
}

// ---- CreateExpense (installment plan) ----

func intPtr(v int) *int { return &v }

func TestCreateExpense_WithInstallments(t *testing.T) {
	now := pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
	var captured psql.CreateExpenseParams
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, arg psql.CreateExpenseParams) (psql.Expense, error) {
			captured = arg
			var amount, cost pgtype.Numeric
			amount.Scan("1.000")
			cost.Scan("200.000")
			return psql.Expense{
				ID:               1,
				Name:             arg.Name,
				Amount:           amount,
				Cost:             cost,
				InstallmentCount: arg.InstallmentCount,
				CreatedAt:        now,
				UpdatedAt:        now,
			}, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	got, err := svc.CreateExpense(context.Background(), CreateExpenseInput{
		UserID:           1,
		Name:             "Phone",
		Amount:           1,
		Cost:             200,
		InstallmentCount: intPtr(4),
	})
	if err != nil {
		t.Fatalf("CreateExpense: %v", err)
	}
	if !captured.InstallmentCount.Valid || captured.InstallmentCount.Int32 != 4 {
		t.Errorf("InstallmentCount param = %+v, want {4 true}", captured.InstallmentCount)
	}
	if got.InstallmentPlan == nil {
		t.Fatal("InstallmentPlan is nil, want non-nil")
	}
	if got.InstallmentPlan.Count != 4 {
		t.Errorf("Count = %d, want 4", got.InstallmentPlan.Count)
	}
	if got.InstallmentPlan.PaymentAmount < 49.99 || got.InstallmentPlan.PaymentAmount > 50.01 {
		t.Errorf("PaymentAmount = %v, want ~50", got.InstallmentPlan.PaymentAmount)
	}
}

func TestCreateExpense_NoInstallments(t *testing.T) {
	now := pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
	var captured psql.CreateExpenseParams
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, arg psql.CreateExpenseParams) (psql.Expense, error) {
			captured = arg
			var n pgtype.Numeric
			n.Scan("10.000")
			return psql.Expense{ID: 1, Name: arg.Name, Amount: n, Cost: n, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	got, err := svc.CreateExpense(context.Background(), CreateExpenseInput{UserID: 1, Name: "Coffee", Amount: 1, Cost: 10})
	if err != nil {
		t.Fatalf("CreateExpense: %v", err)
	}
	if captured.InstallmentCount.Valid {
		t.Errorf("InstallmentCount param = %+v, want invalid (no split)", captured.InstallmentCount)
	}
	if got.InstallmentPlan != nil {
		t.Errorf("InstallmentPlan = %v, want nil", got.InstallmentPlan)
	}
}

func TestCreateExpense_InvalidInstallmentCount(t *testing.T) {
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, _ psql.CreateExpenseParams) (psql.Expense, error) {
			t.Fatal("query should not be reached for invalid installment count")
			return psql.Expense{}, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	_, err := svc.CreateExpense(context.Background(), CreateExpenseInput{
		UserID:           1,
		Name:             "Phone",
		Amount:           1,
		Cost:             200,
		InstallmentCount: intPtr(1),
	})
	if !errors.Is(err, ErrInvalidInstallmentCount) {
		t.Errorf("CreateExpense(count=1) = %v, want ErrInvalidInstallmentCount", err)
	}
}

func TestUpdateExpense_InvalidInstallmentCount(t *testing.T) {
	stub := &stubExpenseQuerier{
		updateExpenseFn: func(_ context.Context, _ psql.UpdateExpenseParams) (psql.Expense, error) {
			t.Fatal("query should not be reached for invalid installment count")
			return psql.Expense{}, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	_, err := svc.UpdateExpense(context.Background(), UpdateExpenseInput{
		ID: 1, UserID: 1, Name: "Phone", Amount: 1, Cost: 200, InstallmentCount: intPtr(0),
	})
	if !errors.Is(err, ErrInvalidInstallmentCount) {
		t.Errorf("UpdateExpense(count=0) = %v, want ErrInvalidInstallmentCount", err)
	}
}

// ---- Kind (expense vs income) ----

func TestCreateExpense_DefaultsToExpenseKind(t *testing.T) {
	now := pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
	var captured psql.CreateExpenseParams
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, arg psql.CreateExpenseParams) (psql.Expense, error) {
			captured = arg
			var n pgtype.Numeric
			n.Scan("10.000")
			return psql.Expense{ID: 1, Name: arg.Name, Kind: arg.Kind, Amount: n, Cost: n, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	got, err := svc.CreateExpense(context.Background(), CreateExpenseInput{UserID: 1, Name: "Coffee", Amount: 1, Cost: 10})
	if err != nil {
		t.Fatalf("CreateExpense: %v", err)
	}
	if captured.Kind != "expense" {
		t.Errorf("Kind param = %q, want \"expense\"", captured.Kind)
	}
	if got.Kind != models.ExpenseKindExpense {
		t.Errorf("got.Kind = %q, want expense", got.Kind)
	}
}

func TestCreateExpense_IncomeKind(t *testing.T) {
	now := pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
	var captured psql.CreateExpenseParams
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, arg psql.CreateExpenseParams) (psql.Expense, error) {
			captured = arg
			var n pgtype.Numeric
			n.Scan("2500.000")
			return psql.Expense{ID: 1, Name: arg.Name, Kind: arg.Kind, Amount: n, Cost: n, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	got, err := svc.CreateExpense(context.Background(), CreateExpenseInput{UserID: 1, Kind: "income", Name: "Salary", Amount: 1, Cost: 2500})
	if err != nil {
		t.Fatalf("CreateExpense: %v", err)
	}
	if captured.Kind != "income" {
		t.Errorf("Kind param = %q, want \"income\"", captured.Kind)
	}
	if got.Kind != models.ExpenseKindIncome {
		t.Errorf("got.Kind = %q, want income", got.Kind)
	}
}

func TestCreateExpense_InvalidKind(t *testing.T) {
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, _ psql.CreateExpenseParams) (psql.Expense, error) {
			t.Fatal("query should not be reached for invalid kind")
			return psql.Expense{}, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	_, err := svc.CreateExpense(context.Background(), CreateExpenseInput{UserID: 1, Kind: "refund", Name: "x", Amount: 1, Cost: 10})
	if !errors.Is(err, ErrInvalidKind) {
		t.Errorf("CreateExpense(kind=refund) = %v, want ErrInvalidKind", err)
	}
}

func TestGetTotalIncomeThisMonth(t *testing.T) {
	var income pgtype.Numeric
	income.Scan("2500.000")
	stub := &stubExpenseQuerier{
		getIncomeFn: func(_ context.Context, _ int64) (pgtype.Numeric, error) {
			return income, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	got, err := svc.GetTotalIncomeThisMonth(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTotalIncomeThisMonth: %v", err)
	}
	if got.Total != 2500 {
		t.Errorf("Total = %v, want 2500", got.Total)
	}
}

// ---- DeleteExpense ----

func TestDeleteExpense_NotFound(t *testing.T) {
	stub := &stubExpenseQuerier{
		softDeleteFn: func(_ context.Context, _ psql.SoftDeleteExpenseParams) (int64, error) {
			return 0, nil // 0 rows affected = not found
		},
	}
	svc := newExpenseServiceWithStub(stub)
	err := svc.DeleteExpense(context.Background(), 1, 99)
	if !errors.Is(err, ErrExpenseNotFound) {
		t.Errorf("DeleteExpense(missing) = %v, want ErrExpenseNotFound", err)
	}
}

func TestDeleteExpense_DBError(t *testing.T) {
	dbErr := errors.New("connection reset")
	stub := &stubExpenseQuerier{
		softDeleteFn: func(_ context.Context, _ psql.SoftDeleteExpenseParams) (int64, error) {
			return 0, dbErr
		},
	}
	svc := newExpenseServiceWithStub(stub)
	err := svc.DeleteExpense(context.Background(), 1, 99)
	if !errors.Is(err, dbErr) {
		t.Errorf("DeleteExpense(db error) = %v, want %v", err, dbErr)
	}
}

func TestDeleteExpense_Success(t *testing.T) {
	stub := &stubExpenseQuerier{
		softDeleteFn: func(_ context.Context, arg psql.SoftDeleteExpenseParams) (int64, error) {
			if arg.ID != 5 || arg.UserID != 3 {
				t.Errorf("wrong params: %+v", arg)
			}
			return 1, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	if err := svc.DeleteExpense(context.Background(), 3, 5); err != nil {
		t.Errorf("DeleteExpense: unexpected error: %v", err)
	}
}

// ---- UpdateExpense (ErrNoRows mapping) ----

func TestUpdateExpense_NotFound(t *testing.T) {
	stub := &stubExpenseQuerier{
		updateExpenseFn: func(_ context.Context, _ psql.UpdateExpenseParams) (psql.Expense, error) {
			return psql.Expense{}, pgx.ErrNoRows
		},
	}
	svc := newExpenseServiceWithStub(stub)
	_, err := svc.UpdateExpense(context.Background(), UpdateExpenseInput{ID: 1, UserID: 1, Name: "X", Amount: 1, Cost: 1})
	if !errors.Is(err, ErrExpenseNotFound) {
		t.Errorf("UpdateExpense(no rows) = %v, want ErrExpenseNotFound", err)
	}
}

func TestUpdateExpense_DBError(t *testing.T) {
	dbErr := errors.New("timeout")
	stub := &stubExpenseQuerier{
		updateExpenseFn: func(_ context.Context, _ psql.UpdateExpenseParams) (psql.Expense, error) {
			return psql.Expense{}, dbErr
		},
	}
	svc := newExpenseServiceWithStub(stub)
	_, err := svc.UpdateExpense(context.Background(), UpdateExpenseInput{ID: 1, UserID: 1, Name: "X", Amount: 1, Cost: 1})
	if !errors.Is(err, dbErr) {
		t.Errorf("UpdateExpense(db error) = %v, want %v", err, dbErr)
	}
}

// ---- GetById (ErrNoRows mapping) ----

func TestGetById_NotFound(t *testing.T) {
	stub := &stubExpenseQuerier{
		getExpenseByIdFn: func(_ context.Context, _ psql.GetExpenseByIdParams) (psql.Expense, error) {
			return psql.Expense{}, pgx.ErrNoRows
		},
	}
	svc := newExpenseServiceWithStub(stub)
	_, err := svc.GetById(context.Background(), 1, 99)
	if !errors.Is(err, ErrExpenseNotFound) {
		t.Errorf("GetById(no rows) = %v, want ErrExpenseNotFound", err)
	}
}

func TestGetById_Success(t *testing.T) {
	now := pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
	var amount pgtype.Numeric
	amount.Scan("25.000")
	var cost pgtype.Numeric
	cost.Scan("25.000")

	stub := &stubExpenseQuerier{
		getExpenseByIdFn: func(_ context.Context, arg psql.GetExpenseByIdParams) (psql.Expense, error) {
			return psql.Expense{
				ID:        arg.ID,
				Name:      "Lunch",
				Amount:    amount,
				Cost:      cost,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}
	svc := newExpenseServiceWithStub(stub)
	got, err := svc.GetById(context.Background(), 1, 7)
	if err != nil {
		t.Fatalf("GetById: %v", err)
	}
	if got.Id != 7 {
		t.Errorf("Id = %d, want 7", got.Id)
	}
	if got.Name != "Lunch" {
		t.Errorf("Name = %q, want Lunch", got.Name)
	}
}
