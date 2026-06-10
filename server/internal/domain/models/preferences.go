package models

// Default preferences, used when a user has no stored row yet. Kept in one place
// so the service fallback and any future caller agree. (The DB column DEFAULTs in
// the migrations mirror these.)
const (
	DefaultAccentID             = "amber"
	DefaultFontSize             = "normal"
	DefaultDarkMode             = true
	DefaultCurrency             = "EUR"
	DefaultWeekStart            = "Monday"
	DefaultMonthlyBudget        = float32(0)
	DefaultNotifyWeeklySummary  = true
	DefaultNotifyBudgetAlerts   = true
	DefaultNotifyProductUpdates = false
)

// ValidFontSizes / ValidCurrencies / ValidWeekStarts whitelist the values Upsert
// accepts, so a client cannot persist arbitrary strings.
var (
	ValidFontSizes  = map[string]bool{"small": true, "normal": true, "large": true}
	ValidCurrencies = map[string]bool{"EUR": true, "USD": true}
	ValidWeekStarts = map[string]bool{"Monday": true, "Sunday": true, "Saturday": true}
)

// DefaultPreferences returns the fallback preferences for a user with no row.
func DefaultPreferences(userID int64) UserPreferences {
	return UserPreferences{
		UserID:               userID,
		AccentID:             DefaultAccentID,
		FontSize:             DefaultFontSize,
		DarkMode:             DefaultDarkMode,
		Currency:             DefaultCurrency,
		WeekStart:            DefaultWeekStart,
		MonthlyBudget:        DefaultMonthlyBudget,
		NotifyWeeklySummary:  DefaultNotifyWeeklySummary,
		NotifyBudgetAlerts:   DefaultNotifyBudgetAlerts,
		NotifyProductUpdates: DefaultNotifyProductUpdates,
	}
}

type UserPreferences struct {
	UserID               int64   `json:"userId"`
	AccentID             string  `json:"accentId"`
	FontSize             string  `json:"fontSize"`
	DarkMode             bool    `json:"darkMode"`
	Currency             string  `json:"currency"`
	WeekStart            string  `json:"weekStart"`
	MonthlyBudget        float32 `json:"monthlyBudget"`
	NotifyWeeklySummary  bool    `json:"notifyWeeklySummary"`
	NotifyBudgetAlerts   bool    `json:"notifyBudgetAlerts"`
	NotifyProductUpdates bool    `json:"notifyProductUpdates"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

type UpsertPreferencesInput struct {
	UserID               int64
	AccentID             string
	FontSize             string
	DarkMode             bool
	Currency             string
	WeekStart            string
	MonthlyBudget        float32
	NotifyWeeklySummary  bool
	NotifyBudgetAlerts   bool
	NotifyProductUpdates bool
}
