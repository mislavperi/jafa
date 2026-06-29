package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/apperr"
	"github.com/mislavperi/jafa/server/internal/domain/dto"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/utils"
)

type ExpenseService struct {
	Queries *psql.Queries
	Pool    psql.Pool
	Mapper  *mappers.ExpenseMapper
}

func NewExpenseService(pool psql.Pool) *ExpenseService {
	return &ExpenseService{
		Queries: psql.New(pool),
		Pool:    pool,
		Mapper:  mappers.NewExpenseMapper(),
	}
}

// resolveKind defaults an empty kind to expense and validates the value.
func resolveKind(kind string) (string, error) {
	if kind == "" {
		return string(models.ExpenseKindExpense), nil
	}
	switch models.ExpenseKind(kind) {
	case models.ExpenseKindExpense, models.ExpenseKindIncome:
		return kind, nil
	default:
		return "", apperr.ErrInvalidKind
	}
}

// recurringScheduleParams builds the recurrence query params from an optional
// schedule. A nil schedule yields zero (NULL) params; an unparseable start date
// is rejected with apperr.ErrInvalidStartDate.
func recurringScheduleParams(schedule *models.RecurringSchedule) (interval pgtype.Text, day pgtype.Int4, startDate pgtype.Date, err error) {
	if schedule == nil {
		return
	}
	start, err := utils.ParseDate(schedule.StartDate)
	if err != nil {
		return pgtype.Text{}, pgtype.Int4{}, pgtype.Date{}, apperr.ErrInvalidStartDate
	}
	return pgtype.Text{String: string(schedule.Interval), Valid: true},
		pgtype.Int4{Int32: int32(schedule.DayOfMonth), Valid: true},
		pgtype.Date{Time: start, Valid: true},
		nil
}

// installmentCountParam validates an optional installment count and converts it
// to the pgtype the queries expect. A nil count means "no split" (one-time
// payment). Any non-nil count below 2 is rejected.
func installmentCountParam(count *int) (pgtype.Int4, error) {
	if count == nil {
		return pgtype.Int4{}, nil
	}
	if *count < 2 {
		return pgtype.Int4{}, apperr.ErrInvalidInstallmentCount
	}
	return pgtype.Int4{Int32: int32(*count), Valid: true}, nil
}

func (es *ExpenseService) CreateExpense(ctx context.Context, input dto.CreateExpenseInput) (models.Expense, error) {
	kind, err := resolveKind(input.Kind)
	if err != nil {
		return models.Expense{}, err
	}
	amount, err := utils.FloatToNumeric(input.Amount)
	if err != nil {
		return models.Expense{}, err
	}
	cost, err := utils.FloatToNumeric(input.Cost)
	if err != nil {
		return models.Expense{}, err
	}
	recurrenceInterval, recurrenceDay, recurrenceStartDate, err := recurringScheduleParams(input.RecurringSchedule)
	if err != nil {
		return models.Expense{}, err
	}
	installmentCount, err := installmentCountParam(input.InstallmentCount)
	if err != nil {
		return models.Expense{}, err
	}
	expense, err := es.Queries.CreateExpense(ctx, psql.CreateExpenseParams{
		UserID:              input.UserID,
		Kind:                kind,
		Name:                input.Name,
		Amount:              amount,
		Cost:                cost,
		RecurrenceInterval:  recurrenceInterval,
		RecurrenceDay:       recurrenceDay,
		RecurrenceStartDate: recurrenceStartDate,
		InstallmentCount:    installmentCount,
	})
	if err != nil {
		return models.Expense{}, err
	}
	return es.Mapper.MapToDomain(expense)
}

func (es *ExpenseService) UpdateExpense(ctx context.Context, input dto.UpdateExpenseInput) (models.Expense, error) {
	amount, err := utils.FloatToNumeric(input.Amount)
	if err != nil {
		return models.Expense{}, err
	}
	cost, err := utils.FloatToNumeric(input.Cost)
	if err != nil {
		return models.Expense{}, err
	}
	recurrenceInterval, recurrenceDay, recurrenceStartDate, err := recurringScheduleParams(input.RecurringSchedule)
	if err != nil {
		return models.Expense{}, err
	}
	installmentCount, err := installmentCountParam(input.InstallmentCount)
	if err != nil {
		return models.Expense{}, err
	}
	expense, err := es.Queries.UpdateExpense(ctx, psql.UpdateExpenseParams{
		ID:                  input.ID,
		UserID:              input.UserID,
		Name:                input.Name,
		Amount:              amount,
		Cost:                cost,
		RecurrenceInterval:  recurrenceInterval,
		RecurrenceDay:       recurrenceDay,
		RecurrenceStartDate: recurrenceStartDate,
		InstallmentCount:    installmentCount,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Expense{}, apperr.ErrExpenseNotFound
		}
		return models.Expense{}, err
	}
	return es.Mapper.MapToDomain(expense)
}

func (es *ExpenseService) DeleteExpense(ctx context.Context, userID, id int64) error {
	rows, err := es.Queries.SoftDeleteExpense(ctx, psql.SoftDeleteExpenseParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperr.ErrExpenseNotFound
	}
	return nil
}
