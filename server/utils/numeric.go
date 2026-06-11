package utils

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func FloatToNumeric(f float32) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(fmt.Sprintf("%.3f", f))
	return n, err
}

// NumericToFloat converts a pgtype.Numeric to float32, treating NULL as 0.
func NumericToFloat(n pgtype.Numeric) (float32, error) {
	f, err := n.Float64Value()
	if err != nil {
		return 0, err
	}
	if !f.Valid {
		return 0, nil
	}
	return float32(f.Float64), nil
}
