package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/utils"
)

type PreferencesService struct {
	Queries *psql.Queries
}

func NewPreferencesService(queries *psql.Queries) *PreferencesService {
	return &PreferencesService{Queries: queries}
}

func mapPreferencesRow(row psql.UserPreference) (models.UserPreferences, error) {
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

func (ps *PreferencesService) Get(ctx context.Context, userID int64) (models.UserPreferences, error) {
	row, err := ps.Queries.GetUserPreferences(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.DefaultPreferences(userID), nil
		}
		return models.UserPreferences{}, err
	}
	return mapPreferencesRow(row)
}

func (ps *PreferencesService) Upsert(ctx context.Context, input models.UpsertPreferencesInput) (models.UserPreferences, error) {
	budget, err := utils.FloatToNumeric(input.MonthlyBudget)
	if err != nil {
		return models.UserPreferences{}, err
	}
	row, err := ps.Queries.UpsertUserPreferences(ctx, psql.UpsertUserPreferencesParams{
		UserID:               input.UserID,
		AccentID:             input.AccentID,
		FontSize:             input.FontSize,
		DarkMode:             input.DarkMode,
		Currency:             input.Currency,
		WeekStart:            input.WeekStart,
		MonthlyBudget:        budget,
		NotifyWeeklySummary:  input.NotifyWeeklySummary,
		NotifyBudgetAlerts:   input.NotifyBudgetAlerts,
		NotifyProductUpdates: input.NotifyProductUpdates,
	})
	if err != nil {
		return models.UserPreferences{}, err
	}
	return mapPreferencesRow(row)
}
