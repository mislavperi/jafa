package dto

// UpsertPreferencesInput is the command PreferencesService.Upsert needs. The
// controller resolves optional/omitted fields against stored values before
// building this, so every field here is the final value to persist.
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
