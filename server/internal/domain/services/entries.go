package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mislavperi/jafa/server/internal/domain/apperr"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

// Read-side queries for expenses: fetching single entries and listing them.

func (es *ExpenseService) GetById(ctx context.Context, userID, id int64) (models.Expense, error) {
	expense, err := es.Queries.GetExpenseById(ctx, psql.GetExpenseByIdParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Expense{}, apperr.ErrExpenseNotFound
		}
		return models.Expense{}, err
	}
	return es.Mapper.MapToDomain(expense)
}

func (es *ExpenseService) GetAllExpenses(ctx context.Context, userID int64) ([]models.Expense, error) {
	expenses, err := es.Queries.GetAllExpenses(ctx, userID)
	if err != nil {
		return nil, err
	}
	return es.Mapper.MapManyToDomain(expenses)
}

// GetAllEntries returns both expenses and income for the user, newest first.
// Used by the transactions table where the two kinds are shown together.
func (es *ExpenseService) GetAllEntries(ctx context.Context, userID int64) ([]models.Expense, error) {
	entries, err := es.Queries.GetAllEntries(ctx, userID)
	if err != nil {
		return nil, err
	}
	return es.Mapper.MapManyToDomain(entries)
}

func (es *ExpenseService) GetExpensesByMonth(ctx context.Context, userID int64, year, month int32) ([]models.Expense, error) {
	expenses, err := es.Queries.GetExpensesByMonth(ctx, psql.GetExpensesByMonthParams{
		UserID: userID,
		Year:   year,
		Month:  month,
	})
	if err != nil {
		return nil, err
	}
	return es.Mapper.MapManyToDomain(expenses)
}

func (es *ExpenseService) GetFirstExpenseDate(ctx context.Context, userID int64) (string, error) {
	result, err := es.Queries.GetFirstExpenseDate(ctx, userID)
	if err != nil {
		return "", err
	}
	s, ok := result.(string)
	if !ok {
		return "", nil
	}
	return s, nil
}
