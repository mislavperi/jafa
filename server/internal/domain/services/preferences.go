package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mislavperi/jafa/server/internal/domain/dto"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/utils"
)

type PreferencesService struct {
	Queries *psql.Queries
	Mapper  *mappers.PreferencesMapper
}

func NewPreferencesService(pool psql.Pool) *PreferencesService {
	return &PreferencesService{Queries: psql.New(pool), Mapper: mappers.NewPreferencesMapper()}
}

func (ps *PreferencesService) Get(ctx context.Context, userID int64) (models.UserPreferences, error) {
	row, err := ps.Queries.GetUserPreferences(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.DefaultPreferences(userID), nil
		}
		return models.UserPreferences{}, err
	}
	return ps.Mapper.MapToDomain(row)
}

func (ps *PreferencesService) Upsert(ctx context.Context, input dto.UpsertPreferencesInput) (models.UserPreferences, error) {
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
	return ps.Mapper.MapToDomain(row)
}
