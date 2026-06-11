package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestParseDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		got, err := ParseDate("2024-06-15")
		if err != nil {
			t.Fatalf("ParseDate returned error: %v", err)
		}
		if got.Year() != 2024 || int(got.Month()) != 6 || got.Day() != 15 {
			t.Errorf("ParseDate = %v, want 2024-06-15", got)
		}
	})

	t.Run("invalid date returns error", func(t *testing.T) {
		if _, err := ParseDate("not-a-date"); err == nil {
			t.Error("ParseDate(invalid) = nil error, want error")
		}
	})

	t.Run("empty string returns error", func(t *testing.T) {
		if _, err := ParseDate(""); err == nil {
			t.Error("ParseDate(\"\") = nil error, want error")
		}
	})
}

func TestFormatRFC3339_Valid(t *testing.T) {
	ts := pgtype.Timestamptz{
		Time:  time.Date(2024, 3, 10, 14, 30, 0, 0, time.UTC),
		Valid: true,
	}
	got := FormatRFC3339(ts)
	if got != "2024-03-10T14:30:00Z" {
		t.Errorf("FormatRFC3339 = %q, want %q", got, "2024-03-10T14:30:00Z")
	}
}

func TestFormatRFC3339_ValidNonUTC(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	ts := pgtype.Timestamptz{
		Time:  time.Date(2024, 3, 10, 9, 0, 0, 0, loc),
		Valid: true,
	}
	got := FormatRFC3339(ts)
	if !strings.HasSuffix(got, "Z") {
		t.Errorf("FormatRFC3339 = %q, expected UTC Z suffix", got)
	}
}

func TestFormatRFC3339_Null(t *testing.T) {
	var ts pgtype.Timestamptz // zero value = NULL
	got := FormatRFC3339(ts)
	if got != "" {
		t.Errorf("FormatRFC3339(null) = %q, want empty string", got)
	}
}
