package services

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/apperr"
	"github.com/mislavperi/jafa/server/internal/domain/dto"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	"github.com/pashagolub/pgxmock/v4"
)

// Tests for the write operations in expense.go: Create, Update, Delete.

// ---- CreateExpense (installment plan) ----

func TestCreateExpense_WithInstallments(t *testing.T) {
	svc, mock := newExpenseService(t)
	// Pin the installment-count param (7th createExpense arg); the rest are free.
	mock.ExpectQuery("CreateExpense").
		WithArgs(
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgtype.Int4{Int32: 4, Valid: true},
			pgxmock.AnyArg(), pgxmock.AnyArg(),
		).
		WillReturnRows(pgxmock.NewRows(expenseColumns).AddRow(
			returnedExpense(t, 1, "Phone", "1.000", "200.000", pgtype.Int4{Int32: 4, Valid: true}, "expense")...,
		))

	got, err := svc.CreateExpense(context.Background(), dto.CreateExpenseInput{
		UserID: 1, Name: "Phone", Amount: 1, Cost: 200, InstallmentCount: intPtr(4),
	})
	if err != nil {
		t.Fatalf("CreateExpense: %v", err)
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
	svc, mock := newExpenseService(t)
	// Installment-count param must be the zero (invalid) Int4 when no split.
	mock.ExpectQuery("CreateExpense").
		WithArgs(
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgtype.Int4{},
			pgxmock.AnyArg(), pgxmock.AnyArg(),
		).
		WillReturnRows(pgxmock.NewRows(expenseColumns).AddRow(
			returnedExpense(t, 1, "Coffee", "10.000", "10.000", pgtype.Int4{}, "expense")...,
		))

	got, err := svc.CreateExpense(context.Background(), dto.CreateExpenseInput{UserID: 1, Name: "Coffee", Amount: 1, Cost: 10})
	if err != nil {
		t.Fatalf("CreateExpense: %v", err)
	}
	if got.InstallmentPlan != nil {
		t.Errorf("InstallmentPlan = %v, want nil", got.InstallmentPlan)
	}
}

func TestCreateExpense_InvalidInstallmentCount(t *testing.T) {
	svc, _ := newExpenseService(t) // no query expected: validation fails first
	_, err := svc.CreateExpense(context.Background(), dto.CreateExpenseInput{
		UserID: 1, Name: "Phone", Amount: 1, Cost: 200, InstallmentCount: intPtr(1),
	})
	if !errors.Is(err, apperr.ErrInvalidInstallmentCount) {
		t.Errorf("CreateExpense(count=1) = %v, want apperr.ErrInvalidInstallmentCount", err)
	}
}

func TestUpdateExpense_InvalidInstallmentCount(t *testing.T) {
	svc, _ := newExpenseService(t)
	_, err := svc.UpdateExpense(context.Background(), dto.UpdateExpenseInput{
		ID: 1, UserID: 1, Name: "Phone", Amount: 1, Cost: 200, InstallmentCount: intPtr(0),
	})
	if !errors.Is(err, apperr.ErrInvalidInstallmentCount) {
		t.Errorf("UpdateExpense(count=0) = %v, want apperr.ErrInvalidInstallmentCount", err)
	}
}

// ---- Kind (expense vs income) ----

func TestCreateExpense_DefaultsToExpenseKind(t *testing.T) {
	svc, mock := newExpenseService(t)
	// Kind is the 9th createExpense arg; it must default to "expense".
	mock.ExpectQuery("CreateExpense").
		WithArgs(
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), "expense",
		).
		WillReturnRows(pgxmock.NewRows(expenseColumns).AddRow(
			returnedExpense(t, 1, "Coffee", "10.000", "10.000", pgtype.Int4{}, "expense")...,
		))

	got, err := svc.CreateExpense(context.Background(), dto.CreateExpenseInput{UserID: 1, Name: "Coffee", Amount: 1, Cost: 10})
	if err != nil {
		t.Fatalf("CreateExpense: %v", err)
	}
	if got.Kind != models.ExpenseKindExpense {
		t.Errorf("got.Kind = %q, want expense", got.Kind)
	}
}

func TestCreateExpense_IncomeKind(t *testing.T) {
	svc, mock := newExpenseService(t)
	mock.ExpectQuery("CreateExpense").
		WithArgs(
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), "income",
		).
		WillReturnRows(pgxmock.NewRows(expenseColumns).AddRow(
			returnedExpense(t, 1, "Salary", "1.000", "2500.000", pgtype.Int4{}, "income")...,
		))

	got, err := svc.CreateExpense(context.Background(), dto.CreateExpenseInput{UserID: 1, Kind: "income", Name: "Salary", Amount: 1, Cost: 2500})
	if err != nil {
		t.Fatalf("CreateExpense: %v", err)
	}
	if got.Kind != models.ExpenseKindIncome {
		t.Errorf("got.Kind = %q, want income", got.Kind)
	}
}

func TestCreateExpense_InvalidKind(t *testing.T) {
	svc, _ := newExpenseService(t) // no query expected: validation fails first
	_, err := svc.CreateExpense(context.Background(), dto.CreateExpenseInput{UserID: 1, Kind: "refund", Name: "x", Amount: 1, Cost: 10})
	if !errors.Is(err, apperr.ErrInvalidKind) {
		t.Errorf("CreateExpense(kind=refund) = %v, want apperr.ErrInvalidKind", err)
	}
}

// ---- UpdateExpense (ErrNoRows mapping) ----

func TestUpdateExpense_NotFound(t *testing.T) {
	svc, mock := newExpenseService(t)
	mock.ExpectQuery("UpdateExpense").WithArgs(
		pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
		pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
	).WillReturnError(pgx.ErrNoRows)

	_, err := svc.UpdateExpense(context.Background(), dto.UpdateExpenseInput{ID: 1, UserID: 1, Name: "X", Amount: 1, Cost: 1})
	if !errors.Is(err, apperr.ErrExpenseNotFound) {
		t.Errorf("UpdateExpense(no rows) = %v, want apperr.ErrExpenseNotFound", err)
	}
}

func TestUpdateExpense_DBError(t *testing.T) {
	svc, mock := newExpenseService(t)
	dbErr := errors.New("timeout")
	mock.ExpectQuery("UpdateExpense").WithArgs(
		pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
		pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
	).WillReturnError(dbErr)

	_, err := svc.UpdateExpense(context.Background(), dto.UpdateExpenseInput{ID: 1, UserID: 1, Name: "X", Amount: 1, Cost: 1})
	if !errors.Is(err, dbErr) {
		t.Errorf("UpdateExpense(db error) = %v, want %v", err, dbErr)
	}
}

// ---- DeleteExpense ----

func TestDeleteExpense_NotFound(t *testing.T) {
	svc, mock := newExpenseService(t)
	mock.ExpectExec("SoftDeleteExpense").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	if err := svc.DeleteExpense(context.Background(), 1, 99); !errors.Is(err, apperr.ErrExpenseNotFound) {
		t.Errorf("DeleteExpense(missing) = %v, want apperr.ErrExpenseNotFound", err)
	}
}

func TestDeleteExpense_DBError(t *testing.T) {
	svc, mock := newExpenseService(t)
	dbErr := errors.New("connection reset")
	mock.ExpectExec("SoftDeleteExpense").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(dbErr)

	if err := svc.DeleteExpense(context.Background(), 1, 99); !errors.Is(err, dbErr) {
		t.Errorf("DeleteExpense(db error) = %v, want %v", err, dbErr)
	}
}

func TestDeleteExpense_Success(t *testing.T) {
	svc, mock := newExpenseService(t)
	// softDeleteExpense args are (id, userID).
	mock.ExpectExec("SoftDeleteExpense").WithArgs(int64(5), int64(3)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	if err := svc.DeleteExpense(context.Background(), 3, 5); err != nil {
		t.Errorf("DeleteExpense: unexpected error: %v", err)
	}
}
