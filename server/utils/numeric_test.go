package utils

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestFloatToNumeric(t *testing.T) {
	cases := []struct {
		input float32
		want  string
	}{
		{10.5, "10.500"},
		{0, "0.000"},
		{-3.14159, "-3.142"},
		{100, "100.000"},
	}
	for _, tc := range cases {
		n, err := FloatToNumeric(tc.input)
		if err != nil {
			t.Errorf("FloatToNumeric(%v): unexpected error: %v", tc.input, err)
			continue
		}
		if !n.Valid {
			t.Errorf("FloatToNumeric(%v): got invalid numeric", tc.input)
		}
	}
}

func TestNumericToFloat_Valid(t *testing.T) {
	var n pgtype.Numeric
	if err := n.Scan("42.500"); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	got, err := NumericToFloat(n)
	if err != nil {
		t.Fatalf("NumericToFloat: %v", err)
	}
	if got < 42.49 || got > 42.51 {
		t.Errorf("NumericToFloat = %v, want ~42.5", got)
	}
}

func TestNumericToFloat_Null(t *testing.T) {
	var n pgtype.Numeric // zero value = NULL
	got, err := NumericToFloat(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0 {
		t.Errorf("NumericToFloat(null) = %v, want 0", got)
	}
}

func TestParseDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		got := ParseDate("2024-06-15")
		if got.Year() != 2024 || int(got.Month()) != 6 || got.Day() != 15 {
			t.Errorf("ParseDate = %v, want 2024-06-15", got)
		}
	})

	t.Run("invalid date returns zero time", func(t *testing.T) {
		got := ParseDate("not-a-date")
		if !got.IsZero() {
			t.Errorf("ParseDate(invalid) = %v, want zero time", got)
		}
	})

	t.Run("empty string returns zero time", func(t *testing.T) {
		got := ParseDate("")
		if !got.IsZero() {
			t.Errorf("ParseDate(\"\") = %v, want zero time", got)
		}
	})
}
