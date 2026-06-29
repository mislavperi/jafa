package services

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

// Tests for the aggregate queries in spending.go.

func TestGetTotalIncomeThisMonth(t *testing.T) {
	svc, mock := newExpenseService(t)
	mock.ExpectQuery("GetTotalIncomeThisMonth").WithArgs(int64(1)).WillReturnRows(
		pgxmock.NewRows([]string{"total"}).AddRow(scanNumeric(t, "2500.000")),
	)

	got, err := svc.GetTotalIncomeThisMonth(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTotalIncomeThisMonth: %v", err)
	}
	if got.Total != 2500 {
		t.Errorf("Total = %v, want 2500", got.Total)
	}
}
