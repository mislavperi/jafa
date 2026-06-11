package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ParseDate parses a YYYY-MM-DD date string.
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatRFC3339 renders a pgtype.Timestamptz as a UTC RFC3339 string, or "" when
// the value is NULL/invalid. Used for API responses so timestamps parse cleanly
// in JS `new Date(...)`.
func FormatRFC3339(t pgtype.Timestamptz) string {
	if !t.Valid {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}
