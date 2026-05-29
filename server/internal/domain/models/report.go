package models

// CategoryBreakdown is the spend for a single category against its budget,
// computed on the backend so the frontend doesn't have to categorize.
type CategoryBreakdown struct {
	Name      string  `json:"name"`
	Icon      string  `json:"icon"`
	Color     string  `json:"color"`
	Budget    float32 `json:"budget"`
	Spent     float32 `json:"spent"`
	Remaining float32 `json:"remaining"`
	Pct       int     `json:"pct"`
}

// MonthlySpend is the summed spend for a single calendar month (YYYY-MM).
type MonthlySpend struct {
	Month string  `json:"month"`
	Total float32 `json:"total"`
}
