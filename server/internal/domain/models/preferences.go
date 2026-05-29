package models

type UserPreferences struct {
	UserID    int64  `json:"userId"`
	AccentID  string `json:"accentId"`
	FontSize  string `json:"fontSize"`
	DarkMode  bool   `json:"darkMode"`
	Currency  string `json:"currency"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpsertPreferencesInput struct {
	UserID   int64
	AccentID string
	FontSize string
	DarkMode bool
	Currency string
}
