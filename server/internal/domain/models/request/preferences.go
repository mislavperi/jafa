package request

// UpsertPreferencesRequest is the body for updating a user's preferences.
// Optional fields use pointers so an omitted field falls back to the stored
// value rather than resetting it.
type UpsertPreferencesRequest struct {
	AccentID             string   `json:"accentId" binding:"required"`
	FontSize             string   `json:"fontSize" binding:"required"`
	DarkMode             bool     `json:"darkMode"`
	Currency             string   `json:"currency"`
	WeekStart            string   `json:"weekStart"`
	MonthlyBudget        *float32 `json:"monthlyBudget"`
	NotifyWeeklySummary  *bool    `json:"notifyWeeklySummary"`
	NotifyBudgetAlerts   *bool    `json:"notifyBudgetAlerts"`
	NotifyProductUpdates *bool    `json:"notifyProductUpdates"`
}
