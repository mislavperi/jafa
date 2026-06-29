// Package mappers converts SQLC-generated database rows into domain models.
//
// Every mapper exposes MapToDomain (one row) and MapManyToDomain (a slice).
// The slice variants all share the same shape, so they delegate to mapSlice
// rather than repeating the loop.
package mappers

// mapSlice applies a row-to-domain mapping over a slice, stopping at the first
// error. It returns a non-nil empty slice for empty input so JSON encodes [].
func mapSlice[Row, Domain any](rows []Row, mapOne func(Row) (Domain, error)) ([]Domain, error) {
	mapped := make([]Domain, 0, len(rows))
	for _, row := range rows {
		domain, err := mapOne(row)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, domain)
	}
	return mapped, nil
}
