package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/utils"
)

// ErrExpenseNotFound is returned when an expense does not exist or is not owned
// by the requesting user. Controllers map this to HTTP 404.
var ErrExpenseNotFound = errors.New("expense not found")

// ErrInvalidStartDate is returned when a recurring schedule's start date is not
// a valid YYYY-MM-DD date. Controllers map this to HTTP 400.
var ErrInvalidStartDate = errors.New("invalid recurring schedule start date")

type ExpenseService struct {
	Queries ExpenseQuerier
	Pool    *pgxpool.Pool
	Mapper  *mappers.ExpenseMapper
}

func NewExpenseService(queries *psql.Queries, pool *pgxpool.Pool) *ExpenseService {
	return &ExpenseService{
		Queries: wrapExpenseQueries(queries),
		Pool:    pool,
		Mapper:  mappers.NewExpenseMapper(),
	}
}

func (es *ExpenseService) GetAllExpenses(ctx context.Context, userID int64) ([]models.Expense, error) {
	expenses, err := es.Queries.GetAllExpenses(ctx, userID)
	if err != nil {
		return nil, err
	}
	return es.Mapper.MapManyToDomain(expenses)
}

func (es *ExpenseService) GetTotalSpendThisMonth(ctx context.Context, userID int64) (models.MonthlyTotal, error) {
	total, err := es.Queries.GetTotalSpendThisMonth(ctx, userID)
	if err != nil {
		return models.MonthlyTotal{}, err
	}
	f, err := total.Float64Value()
	if err != nil || !f.Valid {
		return models.MonthlyTotal{}, err
	}
	return models.MonthlyTotal{
		Total: float32(f.Float64),
	}, nil
}

func (es *ExpenseService) GetDailySpend(ctx context.Context, userID int64, months int32) ([]models.DailySpend, error) {
	rows, err := es.Queries.GetDailySpend(ctx, psql.GetDailySpendParams{
		UserID: userID,
		Months: months,
	})
	if err != nil {
		return nil, err
	}
	result := make([]models.DailySpend, 0, len(rows))
	for _, row := range rows {
		f, err := row.Total.Float64Value()
		if err != nil || !f.Valid {
			return nil, err
		}
		result = append(result, models.DailySpend{
			Day:   row.Day.Time.Format("2006-01-02"),
			Total: float32(f.Float64),
		})
	}
	return result, nil
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

func (es *ExpenseService) GetDailySpendForMonth(ctx context.Context, userID int64, year, month int32) ([]models.DailySpend, error) {
	rows, err := es.Queries.GetDailySpendForMonth(ctx, psql.GetDailySpendForMonthParams{
		UserID: userID,
		Year:   year,
		Month:  month,
	})
	if err != nil {
		return nil, err
	}
	result := make([]models.DailySpend, 0, len(rows))
	for _, row := range rows {
		f, err := row.Total.Float64Value()
		if err != nil || !f.Valid {
			return nil, err
		}
		result = append(result, models.DailySpend{
			Day:   row.Day.Time.Format("2006-01-02"),
			Total: float32(f.Float64),
		})
	}
	return result, nil
}

type CreateExpenseInput struct {
	UserID            int64
	Name              string
	Amount            float32
	Cost              float32
	RecurringSchedule *models.RecurringSchedule
}

func (es *ExpenseService) CreateExpense(ctx context.Context, input CreateExpenseInput) (models.Expense, error) {
	amount, err := utils.FloatToNumeric(input.Amount)
	if err != nil {
		return models.Expense{}, err
	}
	cost, err := utils.FloatToNumeric(input.Cost)
	if err != nil {
		return models.Expense{}, err
	}
	var recurrenceInterval pgtype.Text
	var recurrenceDay pgtype.Int4
	var recurrenceStartDate pgtype.Date
	if input.RecurringSchedule != nil {
		startDate, err := utils.ParseDate(input.RecurringSchedule.StartDate)
		if err != nil {
			return models.Expense{}, ErrInvalidStartDate
		}
		recurrenceInterval = pgtype.Text{String: string(input.RecurringSchedule.Interval), Valid: true}
		recurrenceDay = pgtype.Int4{Int32: int32(input.RecurringSchedule.DayOfMonth), Valid: true}
		recurrenceStartDate = pgtype.Date{Time: startDate, Valid: true}
	}
	expense, err := es.Queries.CreateExpense(ctx, psql.CreateExpenseParams{
		UserID:              input.UserID,
		Name:                input.Name,
		Amount:              amount,
		Cost:                cost,
		RecurrenceInterval:  recurrenceInterval,
		RecurrenceDay:       recurrenceDay,
		RecurrenceStartDate: recurrenceStartDate,
	})
	if err != nil {
		return models.Expense{}, err
	}
	return es.Mapper.MapToDomain(expense)
}

// BulkExpenseItem is one expense in a bulk import (e.g. a receipt line item).
// Tag is optional; when set, the tag is created if it does not exist yet and
// linked to the new expense.
type BulkExpenseItem struct {
	Name   string
	Amount float32
	Cost   float32
	Tag    string
}

// BulkCreateExpenses creates all items (and their tag links) in one
// transaction, so a receipt import either fully succeeds or leaves no trace.
func (es *ExpenseService) BulkCreateExpenses(ctx context.Context, userID int64, items []BulkExpenseItem) ([]models.Expense, error) {
	tx, err := es.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := es.Queries.WithTx(tx)

	created := make([]models.Expense, 0, len(items))
	for _, item := range items {
		amount, err := utils.FloatToNumeric(item.Amount)
		if err != nil {
			return nil, err
		}
		cost, err := utils.FloatToNumeric(item.Cost)
		if err != nil {
			return nil, err
		}
		row, err := qtx.CreateExpense(ctx, psql.CreateExpenseParams{
			UserID: userID,
			Name:   item.Name,
			Amount: amount,
			Cost:   cost,
		})
		if err != nil {
			return nil, err
		}
		if item.Tag != "" {
			tag, err := qtx.UpsertTag(ctx, psql.UpsertTagParams{
				Name:   item.Tag,
				UserID: userID,
			})
			if err != nil {
				return nil, err
			}
			if err := qtx.AddTagToExpense(ctx, psql.AddTagToExpenseParams{
				ExpenseID: row.ID,
				TagID:     tag.ID,
				UserID:    userID,
			}); err != nil {
				return nil, err
			}
		}
		expense, err := es.Mapper.MapToDomain(row)
		if err != nil {
			return nil, err
		}
		created = append(created, expense)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}

type UpdateExpenseInput struct {
	ID                int64
	UserID            int64
	Name              string
	Amount            float32
	Cost              float32
	RecurringSchedule *models.RecurringSchedule
}

func (es *ExpenseService) UpdateExpense(ctx context.Context, input UpdateExpenseInput) (models.Expense, error) {
	amount, err := utils.FloatToNumeric(input.Amount)
	if err != nil {
		return models.Expense{}, err
	}
	cost, err := utils.FloatToNumeric(input.Cost)
	if err != nil {
		return models.Expense{}, err
	}
	var recurrenceInterval pgtype.Text
	var recurrenceDay pgtype.Int4
	var recurrenceStartDate pgtype.Date
	if input.RecurringSchedule != nil {
		startDate, err := utils.ParseDate(input.RecurringSchedule.StartDate)
		if err != nil {
			return models.Expense{}, ErrInvalidStartDate
		}
		recurrenceInterval = pgtype.Text{String: string(input.RecurringSchedule.Interval), Valid: true}
		recurrenceDay = pgtype.Int4{Int32: int32(input.RecurringSchedule.DayOfMonth), Valid: true}
		recurrenceStartDate = pgtype.Date{Time: startDate, Valid: true}
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
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Expense{}, ErrExpenseNotFound
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
		return ErrExpenseNotFound
	}
	return nil
}

func (es *ExpenseService) GetById(ctx context.Context, userID, id int64) (models.Expense, error) {
	expense, err := es.Queries.GetExpenseById(ctx, psql.GetExpenseByIdParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Expense{}, ErrExpenseNotFound
		}
		return models.Expense{}, err
	}
	return es.Mapper.MapToDomain(expense)
}
