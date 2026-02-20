package models

type DailySpend struct {
	Day   string  `json:"day"`
	Total float32 `json:"total"`
}
