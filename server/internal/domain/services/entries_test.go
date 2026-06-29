package services

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/apperr"
	"github.com/pashagolub/pgxmock/v4"
)

// Tests for the read queries in entries.go.

func TestGetById_NotFound(t *testing.T) {
	svc, mock := newExpenseService(t)
	mock.ExpectQuery("GetExpenseById").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(pgx.ErrNoRows)

	_, err := svc.GetById(context.Background(), 1, 99)
	if !errors.Is(err, apperr.ErrExpenseNotFound) {
		t.Errorf("GetById(no rows) = %v, want apperr.ErrExpenseNotFound", err)
	}
}

func TestGetById_Success(t *testing.T) {
	svc, mock := newExpenseService(t)
	// getExpenseById args are (id, userID).
	mock.ExpectQuery("GetExpenseById").WithArgs(int64(7), int64(1)).WillReturnRows(
		pgxmock.NewRows(expenseColumns).AddRow(
			returnedExpense(t, 7, "Lunch", "25.000", "25.000", pgtype.Int4{}, "expense")...,
		),
	)

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
