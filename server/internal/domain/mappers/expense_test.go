package mappers

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

func makeNumeric(t *testing.T, s string) pgtype.Numeric {
	t.Helper()
	var n pgtype.Numeric
	if err := n.Scan(s); err != nil {
		t.Fatalf("makeNumeric(%q): %v", s, err)
	}
	return n
}

func makeTimestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func TestExpenseMapper_MapToDomain_Basic(t *testing.T) {
	mapper := NewExpenseMapper()
	now := time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC)

	exp := psql.Expense{
		ID:        42,
		Name:      "Coffee",
		Amount:    makeNumeric(t, "3.500"),
		Cost:      makeNumeric(t, "1.000"),
		IsDeleted: false,
		CreatedAt: makeTimestamp(now),
		UpdatedAt: makeTimestamp(now),
	}

	got, err := mapper.MapToDomain(exp)
	if err != nil {
		t.Fatalf("MapToDomain: %v", err)
	}

	if got.Id != 42 {
		t.Errorf("Id = %d, want 42", got.Id)
	}
	if got.Name != "Coffee" {
		t.Errorf("Name = %q, want Coffee", got.Name)
	}
	if got.Amount < 3.49 || got.Amount > 3.51 {
		t.Errorf("Amount = %v, want ~3.5", got.Amount)
	}
	if got.Cost < 0.99 || got.Cost > 1.01 {
		t.Errorf("Cost = %v, want ~1.0", got.Cost)
	}
	if got.RecurringSchedule != nil {
		t.Errorf("RecurringSchedule = %v, want nil", got.RecurringSchedule)
	}
	if got.CreatedAt != "2024-05-01T10:00:00Z" {
		t.Errorf("CreatedAt = %q, want 2024-05-01T10:00:00Z", got.CreatedAt)
	}
}

func TestExpenseMapper_MapToDomain_WithRecurringSchedule(t *testing.T) {
	mapper := NewExpenseMapper()
	now := time.Now().UTC()

	exp := psql.Expense{
		ID:     10,
		Name:   "Netflix",
		Amount: makeNumeric(t, "15.990"),
		Cost:   makeNumeric(t, "15.990"),
		RecurrenceInterval: pgtype.Text{
			String: "monthly",
			Valid:  true,
		},
		RecurrenceDay: pgtype.Int4{Int32: 15, Valid: true},
		RecurrenceStartDate: pgtype.Date{
			Time:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		CreatedAt: makeTimestamp(now),
		UpdatedAt: makeTimestamp(now),
	}

	got, err := mapper.MapToDomain(exp)
	if err != nil {
		t.Fatalf("MapToDomain: %v", err)
	}

	if got.RecurringSchedule == nil {
		t.Fatal("RecurringSchedule is nil, want non-nil")
	}
	if got.RecurringSchedule.Interval != "monthly" {
		t.Errorf("Interval = %q, want monthly", got.RecurringSchedule.Interval)
	}
	if got.RecurringSchedule.DayOfMonth != 15 {
		t.Errorf("DayOfMonth = %d, want 15", got.RecurringSchedule.DayOfMonth)
	}
	if got.RecurringSchedule.StartDate != "2024-01-15" {
		t.Errorf("StartDate = %q, want 2024-01-15", got.RecurringSchedule.StartDate)
	}
}

func TestExpenseMapper_MapManyToDomain_Empty(t *testing.T) {
	mapper := NewExpenseMapper()
	got, err := mapper.MapManyToDomain(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

func TestExpenseMapper_MapManyToDomain_Multiple(t *testing.T) {
	mapper := NewExpenseMapper()
	now := time.Now().UTC()
	ts := makeTimestamp(now)

	expenses := []psql.Expense{
		{ID: 1, Name: "A", Amount: makeNumeric(t, "10.000"), Cost: makeNumeric(t, "10.000"), CreatedAt: ts, UpdatedAt: ts},
		{ID: 2, Name: "B", Amount: makeNumeric(t, "20.000"), Cost: makeNumeric(t, "20.000"), CreatedAt: ts, UpdatedAt: ts},
	}

	got, err := mapper.MapManyToDomain(expenses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Name != "A" || got[1].Name != "B" {
		t.Errorf("names = [%q, %q], want [A, B]", got[0].Name, got[1].Name)
	}
}
