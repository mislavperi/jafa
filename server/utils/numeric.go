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
