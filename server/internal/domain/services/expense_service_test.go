package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

// Ensure stubExpenseQuerier satisfies the interface at compile time.
var _ ExpenseQuerier = (*stubExpenseQuerier)(nil)

// stubExpenseQuerier satisfies ExpenseQuerier with configurable per-method stubs.
type stubExpenseQuerier struct {
	softDeleteFn       func(ctx context.Context, arg psql.SoftDeleteExpenseParams) (int64, error)
	updateExpenseFn    func(ctx context.Context, arg psql.UpdateExpenseParams) (psql.Expense, error)
	getExpenseByIdFn   func(ctx context.Context, arg psql.GetExpenseByIdParams) (psql.Expense, error)
	getAllExpensesFn    func(ctx context.Context, userID int64) ([]psql.Expense, error)
	createExpenseFn    func(ctx context.Context, arg psql.CreateExpenseParams) (psql.Expense, error)
	getTotalFn         func(ctx context.Context, userID int64) (pgtype.Numeric, error)
	getDailySpendFn    func(ctx context.Context, arg psql.GetDailySpendParams) ([]psql.GetDailySpendRow, error)
	getByMonthFn       func(ctx context.Context, arg psql.GetExpensesByMonthParams) ([]psql.Expense, error)
	getFirstDateFn     func(ctx context.Context, userID int64) (interface{}, error)
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
func (s *stubExpenseQuerier) GetTotalSpendThisMonth(ctx context.Context, userID int64) (pgtype.Numeric, error) {
	return s.getTotalFn(ctx, userID)
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


// ---- DeleteExpense ----

func TestDeleteExpense_NotFound(t *testing.T) {
	stub := &stubExpenseQuerier{
		softDeleteFn: func(_ context.Context, _ psql.SoftDeleteExpenseParams) (int64, error) {
			return 0, nil // 0 rows affected = not found
		},
	}
	svc := newExpenseServiceWithStub(stub)
	err := svc.DeleteExpense(1, 99)
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
	err := svc.DeleteExpense(1, 99)
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
	if err := svc.DeleteExpense(3, 5); err != nil {
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
	_, err := svc.UpdateExpense(UpdateExpenseInput{ID: 1, UserID: 1, Name: "X", Amount: 1, Cost: 1})
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
	_, err := svc.UpdateExpense(UpdateExpenseInput{ID: 1, UserID: 1, Name: "X", Amount: 1, Cost: 1})
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
	_, err := svc.GetById(1, 99)
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
	got, err := svc.GetById(1, 7)
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
