package models

// Default preferences, used when a user has no stored row yet. Kept in one place
// so the service fallback and any future caller agree. (The DB column DEFAULTs in
// the migrations mirror these.)
const (
	DefaultAccentID = "amber"
	DefaultFontSize = "normal"
	DefaultDarkMode = true
	DefaultCurrency = "EUR"
)

// ValidFontSizes / ValidCurrencies whitelist the values Upsert accepts, so a
// client cannot persist arbitrary strings.
var (
	ValidFontSizes  = map[string]bool{"small": true, "normal": true, "large": true}
	ValidCurrencies = map[string]bool{"EUR": true, "USD": true}
)

// DefaultPreferences returns the fallback preferences for a user with no row.
func DefaultPreferences(userID int64) UserPreferences {
	return UserPreferences{
		UserID:   userID,
		AccentID: DefaultAccentID,
		FontSize: DefaultFontSize,
		DarkMode: DefaultDarkMode,
		Currency: DefaultCurrency,
	}
}

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
