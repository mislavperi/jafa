package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/utils"
)

type PreferencesMapper struct {
}

func NewPreferencesMapper() *PreferencesMapper {
	return &PreferencesMapper{}
}

func (pm *PreferencesMapper) MapToDomain(row psql.UserPreference) (models.UserPreferences, error) {
	budget, err := utils.NumericToFloat(row.MonthlyBudget)
	if err != nil {
		return models.UserPreferences{}, err
	}
	return models.UserPreferences{
		UserID:               row.UserID,
		AccentID:             row.AccentID,
		FontSize:             row.FontSize,
		DarkMode:             row.DarkMode,
		Currency:             row.Currency,
		WeekStart:            row.WeekStart,
		MonthlyBudget:        budget,
		NotifyWeeklySummary:  row.NotifyWeeklySummary,
		NotifyBudgetAlerts:   row.NotifyBudgetAlerts,
		NotifyProductUpdates: row.NotifyProductUpdates,
		CreatedAt:            utils.FormatRFC3339(row.CreatedAt),
		UpdatedAt:            utils.FormatRFC3339(row.UpdatedAt),
	}, nil
}
