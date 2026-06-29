package services

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/dto"
	"github.com/pashagolub/pgxmock/v4"
)

// Tests for the transactional bulk import in receipts.go.

func TestBulkCreateExpenses_CommitsAll(t *testing.T) {
	svc, mock := newExpenseService(t)
	mock.ExpectBegin()
	mock.ExpectQuery("CreateExpense").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnRows(pgxmock.NewRows(expenseColumns).AddRow(
		returnedExpense(t, 1, "Coffee", "3.000", "3.000", pgtype.Int4{}, "expense")...,
	))
	mock.ExpectQuery("CreateExpense").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnRows(pgxmock.NewRows(expenseColumns).AddRow(
		returnedExpense(t, 2, "Bagel", "2.000", "2.000", pgtype.Int4{}, "expense")...,
	))
	mock.ExpectCommit()

	got, err := svc.BulkCreateExpenses(context.Background(), 1, []dto.BulkExpenseItem{
		{Name: "Coffee", Amount: 1, Cost: 3},
		{Name: "Bagel", Amount: 1, Cost: 2},
	})
	if err != nil {
		t.Fatalf("BulkCreateExpenses: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("created = %d, want 2", len(got))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// A failure mid-import rolls the whole transaction back, so an earlier insert
// does not leak.
func TestBulkCreateExpenses_RollsBackOnError(t *testing.T) {
	svc, mock := newExpenseService(t)
	dbErr := errors.New("insert failed")
	mock.ExpectBegin()
	mock.ExpectQuery("CreateExpense").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnRows(pgxmock.NewRows(expenseColumns).AddRow(
		returnedExpense(t, 1, "Coffee", "3.000", "3.000", pgtype.Int4{}, "expense")...,
	))
	mock.ExpectQuery("CreateExpense").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(dbErr)
	mock.ExpectRollback()

	_, err := svc.BulkCreateExpenses(context.Background(), 1, []dto.BulkExpenseItem{
		{Name: "Coffee", Amount: 1, Cost: 3},
		{Name: "Bagel", Amount: 1, Cost: 2},
	})
	if !errors.Is(err, dbErr) {
		t.Errorf("BulkCreateExpenses = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations (rollback not issued?): %v", err)
	}
}
