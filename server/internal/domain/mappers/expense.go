package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type ExpenseMapper struct {
}

func NewExpenseMapper() *ExpenseMapper {
	return &ExpenseMapper{}
}

func (em *ExpenseMapper) MapToDomain(expense psql.Expense) (models.Expense, error) {
	amount, err := expense.Amount.Float64Value()
	if err != nil || !amount.Valid {
		return models.Expense{}, err
	}
	cost, err := expense.Cost.Float64Value()
	if err != nil || !cost.Valid {
		return models.Expense{}, err
	}
	var recurringSchedule *models.RecurringSchedule
	if expense.RecurrenceInterval.Valid && expense.RecurrenceDay.Valid {
		recurringSchedule = &models.RecurringSchedule{
			Interval:   models.RecurrenceInterval(expense.RecurrenceInterval.String),
			DayOfMonth: int(expense.RecurrenceDay.Int32),
			StartDate:  expense.RecurrenceStartDate.Time.Format("2006-01-02"),
		}
	}
	return models.Expense{
		Id:                expense.ID,
		Name:              expense.Name,
		Amount:            float32(amount.Float64),
		Cost:              float32(cost.Float64),
		ItemID:            expense.ItemID.Int64,
		IsDeleted:         expense.IsDeleted,
		CreatedAt:         expense.CreatedAt.Time.String(),
		UpdatedAt:         expense.UpdatedAt.Time.String(),
		RecurringSchedule: recurringSchedule,
	}, nil
}

func (em *ExpenseMapper) MapManyToDomain(expenses []psql.Expense) ([]models.Expense, error) {
	mappedExpenses := make([]models.Expense, 0, len(expenses))
	for _, expense := range expenses {
		mappedExpense, err := em.MapToDomain(expense)
		if err != nil {
			return []models.Expense{}, err
		}
		mappedExpenses = append(mappedExpenses, mappedExpense)
	}
	return mappedExpenses, nil
}

func (em *ExpenseMapper) MapRowToDomain(expense psql.GetExpensesByMonthRow) (models.Expense, error) {
	amount, err := expense.Amount.Float64Value()
	if err != nil || !amount.Valid {
		return models.Expense{}, err
	}
	cost, err := expense.Cost.Float64Value()
	if err != nil || !cost.Valid {
		return models.Expense{}, err
	}
	return models.Expense{
		Id:        expense.ID,
		Name:      expense.Name,
		Amount:    float32(amount.Float64),
		Cost:      float32(cost.Float64),
		ItemID:    expense.ItemID.Int64,
		IsDeleted: expense.IsDeleted,
		CreatedAt: expense.CreatedAt.Time.String(),
		UpdatedAt: expense.UpdatedAt.Time.String(),
	}, nil
}

func (em *ExpenseMapper) MapManyRowsToDomain(expenses []psql.GetExpensesByMonthRow) ([]models.Expense, error) {
	mappedExpenses := make([]models.Expense, 0, len(expenses))
	for _, expense := range expenses {
		mappedExpense, err := em.MapRowToDomain(expense)
		if err != nil {
			return []models.Expense{}, err
		}
		mappedExpenses = append(mappedExpenses, mappedExpense)
	}
	return mappedExpenses, nil
}
