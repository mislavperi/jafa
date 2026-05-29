package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

func formatRFC3339(t pgtype.Timestamptz) string {
	if !t.Valid {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}

type PreferencesService struct {
	Queries *psql.Queries
}

func NewPreferencesService(queries *psql.Queries) *PreferencesService {
	return &PreferencesService{Queries: queries}
}

func (ps *PreferencesService) Get(userID int64) (models.UserPreferences, error) {
	row, err := ps.Queries.GetUserPreferences(context.Background(), userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.UserPreferences{
				UserID:   userID,
				AccentID: "amber",
				FontSize: "normal",
				DarkMode: true,
				Currency: "EUR",
			}, nil
		}
		return models.UserPreferences{}, err
	}
	return models.UserPreferences{
		UserID:    row.UserID,
		AccentID:  row.AccentID,
		FontSize:  row.FontSize,
		DarkMode:  row.DarkMode,
		Currency:  row.Currency,
		CreatedAt: formatRFC3339(row.CreatedAt),
		UpdatedAt: formatRFC3339(row.UpdatedAt),
	}, nil
}

func (ps *PreferencesService) Upsert(input models.UpsertPreferencesInput) (models.UserPreferences, error) {
	row, err := ps.Queries.UpsertUserPreferences(context.Background(), psql.UpsertUserPreferencesParams{
		UserID:   input.UserID,
		AccentID: input.AccentID,
		FontSize: input.FontSize,
		DarkMode: input.DarkMode,
		Currency: input.Currency,
	})
	if err != nil {
		return models.UserPreferences{}, err
	}
	return models.UserPreferences{
		UserID:    row.UserID,
		AccentID:  row.AccentID,
		FontSize:  row.FontSize,
		DarkMode:  row.DarkMode,
		Currency:  row.Currency,
		CreatedAt: formatRFC3339(row.CreatedAt),
		UpdatedAt: formatRFC3339(row.UpdatedAt),
	}, nil
}
