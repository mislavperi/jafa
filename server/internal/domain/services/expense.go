package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/utils"
)

type ExpenseService struct {
	Queries *psql.Queries
	Mapper  *mappers.ExpenseMapper
}

func NewExpenseService(queries *psql.Queries) *ExpenseService {
	return &ExpenseService{
		Queries: queries,
		Mapper:  mappers.NewExpenseMapper(),
	}
}

func (es *ExpenseService) GetAllExpenses(userID int64) ([]models.Expense, error) {
	expenses, err := es.Queries.GetAllExpenses(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return es.Mapper.MapManyToDomain(expenses)
}

func (es *ExpenseService) GetTotalSpendThisMonth(userID int64) (models.MonthlyTotal, error) {
	total, err := es.Queries.GetTotalSpendThisMonth(context.Background(), userID)
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

func (es *ExpenseService) GetDailySpend(userID int64, months int32) ([]models.DailySpend, error) {
	rows, err := es.Queries.GetDailySpend(context.Background(), psql.GetDailySpendParams{
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

func (es *ExpenseService) GetExpensesByMonth(userID int64, year, month int32) ([]models.Expense, error) {
	expenses, err := es.Queries.GetExpensesByMonth(context.Background(), psql.GetExpensesByMonthParams{
		UserID: userID,
		Year:   year,
		Month:  month,
	})
	if err != nil {
		return nil, err
	}
	return es.Mapper.MapManyToDomain(expenses)
}

func (es *ExpenseService) GetFirstExpenseDate(userID int64) (string, error) {
	result, err := es.Queries.GetFirstExpenseDate(context.Background(), userID)
	if err != nil {
		return "", err
	}
	s, ok := result.(string)
	if !ok {
		return "", nil
	}
	return s, nil
}

func (es *ExpenseService) GetDailySpendForMonth(userID int64, year, month int32) ([]models.DailySpend, error) {
	rows, err := es.Queries.GetDailySpendForMonth(context.Background(), psql.GetDailySpendForMonthParams{
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

func (es *ExpenseService) CreateExpense(input CreateExpenseInput) (models.Expense, error) {
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
		recurrenceInterval = pgtype.Text{String: string(input.RecurringSchedule.Interval), Valid: true}
		recurrenceDay = pgtype.Int4{Int32: int32(input.RecurringSchedule.DayOfMonth), Valid: true}
		recurrenceStartDate = pgtype.Date{Time: utils.ParseDate(input.RecurringSchedule.StartDate), Valid: true}
	}
	expense, err := es.Queries.CreateExpense(context.Background(), psql.CreateExpenseParams{
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

type UpdateExpenseInput struct {
	ID                int64
	UserID            int64
	Name              string
	Amount            float32
	Cost              float32
	RecurringSchedule *models.RecurringSchedule
}

func (es *ExpenseService) UpdateExpense(input UpdateExpenseInput) (models.Expense, error) {
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
		recurrenceInterval = pgtype.Text{String: string(input.RecurringSchedule.Interval), Valid: true}
		recurrenceDay = pgtype.Int4{Int32: int32(input.RecurringSchedule.DayOfMonth), Valid: true}
		recurrenceStartDate = pgtype.Date{Time: utils.ParseDate(input.RecurringSchedule.StartDate), Valid: true}
	}
	expense, err := es.Queries.UpdateExpense(context.Background(), psql.UpdateExpenseParams{
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
		return models.Expense{}, err
	}
	return es.Mapper.MapToDomain(expense)
}

func (es *ExpenseService) DeleteExpense(userID, id int64) error {
	return es.Queries.SoftDeleteExpense(context.Background(), psql.SoftDeleteExpenseParams{
		ID:     id,
		UserID: userID,
	})
}

func (es *ExpenseService) GetById(userID, id int64) (models.Expense, error) {
	expense, err := es.Queries.GetExpenseById(context.Background(), psql.GetExpenseByIdParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return models.Expense{}, err
	}
	return es.Mapper.MapToDomain(expense)
}
