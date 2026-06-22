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

// ErrInvalidInstallmentCount is returned when an expense is split into fewer
// than 2 payments. Controllers map this to HTTP 400.
var ErrInvalidInstallmentCount = errors.New("installment count must be at least 2")

// ErrInvalidKind is returned when an entry's kind is neither "expense" nor
// "income". Controllers map this to HTTP 400.
var ErrInvalidKind = errors.New("kind must be 'expense' or 'income'")

// resolveKind defaults an empty kind to expense and validates the value.
func resolveKind(kind string) (string, error) {
	if kind == "" {
		return string(models.ExpenseKindExpense), nil
	}
	switch models.ExpenseKind(kind) {
	case models.ExpenseKindExpense, models.ExpenseKindIncome:
		return kind, nil
	default:
		return "", ErrInvalidKind
	}
}

// ExpenseQuerier is the subset of psql.Queries used by ExpenseService.
// Extracted as an interface to allow test doubles.
type ExpenseQuerier interface {
	GetAllExpenses(ctx context.Context, userID int64) ([]psql.Expense, error)
	GetAllEntries(ctx context.Context, userID int64) ([]psql.Expense, error)
	GetExpenseById(ctx context.Context, arg psql.GetExpenseByIdParams) (psql.Expense, error)
	CreateExpense(ctx context.Context, arg psql.CreateExpenseParams) (psql.Expense, error)
	UpdateExpense(ctx context.Context, arg psql.UpdateExpenseParams) (psql.Expense, error)
	SoftDeleteExpense(ctx context.Context, arg psql.SoftDeleteExpenseParams) (int64, error)
	GetTotalSpendThisMonth(ctx context.Context, userID int64) (pgtype.Numeric, error)
	GetTotalIncomeThisMonth(ctx context.Context, userID int64) (pgtype.Numeric, error)
	GetDailySpend(ctx context.Context, arg psql.GetDailySpendParams) ([]psql.GetDailySpendRow, error)
	GetExpensesByMonth(ctx context.Context, arg psql.GetExpensesByMonthParams) ([]psql.Expense, error)
	GetFirstExpenseDate(ctx context.Context, userID int64) (interface{}, error)
	GetDailySpendForMonth(ctx context.Context, arg psql.GetDailySpendForMonthParams) ([]psql.GetDailySpendForMonthRow, error)
	UpsertTag(ctx context.Context, arg psql.UpsertTagParams) (psql.Tag, error)
	AddTagToExpense(ctx context.Context, arg psql.AddTagToExpenseParams) error
	WithTx(tx pgx.Tx) ExpenseQuerier
}

// psqlQueriesWrapper wraps *psql.Queries so its WithTx satisfies ExpenseQuerier.
type psqlQueriesWrapper struct {
	*psql.Queries
}

func (w *psqlQueriesWrapper) WithTx(tx pgx.Tx) ExpenseQuerier {
	return &psqlQueriesWrapper{w.Queries.WithTx(tx)}
}

func wrapExpenseQueries(q *psql.Queries) ExpenseQuerier {
	return &psqlQueriesWrapper{q}
}

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

// GetAllEntries returns both expenses and income for the user, newest first.
// Used by the transactions table where the two kinds are shown together.
func (es *ExpenseService) GetAllEntries(ctx context.Context, userID int64) ([]models.Expense, error) {
	entries, err := es.Queries.GetAllEntries(ctx, userID)
	if err != nil {
		return nil, err
	}
	return es.Mapper.MapManyToDomain(entries)
}

// numericToMonthlyTotal converts a pgtype.Numeric sum into a MonthlyTotal,
// treating a NULL/invalid value as a zero total.
func numericToMonthlyTotal(n pgtype.Numeric) (models.MonthlyTotal, error) {
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return models.MonthlyTotal{}, err
	}
	return models.MonthlyTotal{Total: float32(f.Float64)}, nil
}

func (es *ExpenseService) GetTotalSpendThisMonth(ctx context.Context, userID int64) (models.MonthlyTotal, error) {
	total, err := es.Queries.GetTotalSpendThisMonth(ctx, userID)
	if err != nil {
		return models.MonthlyTotal{}, err
	}
	return numericToMonthlyTotal(total)
}

// dailySpendFromRow converts a (day, total) pair from a daily-spend query row
// into a DailySpend. A NULL/invalid total yields a zero-value DailySpend.
func dailySpendFromRow(day pgtype.Date, total pgtype.Numeric) (models.DailySpend, error) {
	f, err := total.Float64Value()
	if err != nil || !f.Valid {
		return models.DailySpend{}, err
	}
	return models.DailySpend{
		Day:   day.Time.Format("2006-01-02"),
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
		spend, err := dailySpendFromRow(row.Day, row.Total)
		if err != nil {
			return nil, err
		}
		result = append(result, spend)
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
		spend, err := dailySpendFromRow(row.Day, row.Total)
		if err != nil {
			return nil, err
		}
		result = append(result, spend)
	}
	return result, nil
}

type CreateExpenseInput struct {
	UserID            int64
	Kind              string
	Name              string
	Amount            float32
	Cost              float32
	RecurringSchedule *models.RecurringSchedule
	InstallmentCount  *int
}

// recurringScheduleParams builds the recurrence query params from an optional
// schedule. A nil schedule yields zero (NULL) params; an unparseable start date
// is rejected with ErrInvalidStartDate.
func recurringScheduleParams(schedule *models.RecurringSchedule) (interval pgtype.Text, day pgtype.Int4, startDate pgtype.Date, err error) {
	if schedule == nil {
		return
	}
	start, err := utils.ParseDate(schedule.StartDate)
	if err != nil {
		return pgtype.Text{}, pgtype.Int4{}, pgtype.Date{}, ErrInvalidStartDate
	}
	return pgtype.Text{String: string(schedule.Interval), Valid: true},
		pgtype.Int4{Int32: int32(schedule.DayOfMonth), Valid: true},
		pgtype.Date{Time: start, Valid: true},
		nil
}

func (es *ExpenseService) CreateExpense(ctx context.Context, input CreateExpenseInput) (models.Expense, error) {
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

func (es *ExpenseService) GetTotalIncomeThisMonth(ctx context.Context, userID int64) (models.MonthlyTotal, error) {
	total, err := es.Queries.GetTotalIncomeThisMonth(ctx, userID)
	if err != nil {
		return models.MonthlyTotal{}, err
	}
	return numericToMonthlyTotal(total)
}

// installmentCountParam validates an optional installment count and converts it
// to the pgtype the queries expect. A nil count means "no split" (one-time
// payment). Any non-nil count below 2 is rejected.
func installmentCountParam(count *int) (pgtype.Int4, error) {
	if count == nil {
		return pgtype.Int4{}, nil
	}
	if *count < 2 {
		return pgtype.Int4{}, ErrInvalidInstallmentCount
	}
	return pgtype.Int4{Int32: int32(*count), Valid: true}, nil
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
			Kind:   string(models.ExpenseKindExpense),
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
	InstallmentCount  *int
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
