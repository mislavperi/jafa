package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type ExpensesTagMapper struct {
}

func NewExpensesTagMapper() *ExpensesTagMapper {
	return &ExpensesTagMapper{}
}

func (etm *ExpensesTagMapper) MapToDomain(expensesTag psql.ExpensesTag) (models.ExpensesTag, error) {
	return models.ExpensesTag{
		ExpenseID: expensesTag.ExpenseID,
		TagID:     expensesTag.TagID,
	}, nil
}

func (etm *ExpensesTagMapper) MapManyToDomain(expensesTags []psql.ExpensesTag) ([]models.ExpensesTag, error) {
	mappedExpensesTags := make([]models.ExpensesTag, 0, len(expensesTags))
	for _, expensesTag := range expensesTags {
		mappedExpensesTag, err := etm.MapToDomain(expensesTag)
		if err != nil {
			return []models.ExpensesTag{}, err
		}
		mappedExpensesTags = append(mappedExpensesTags, mappedExpensesTag)
	}
	return mappedExpensesTags, nil
}
