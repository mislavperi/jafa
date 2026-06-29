package services

import (
	"context"

	"github.com/mislavperi/jafa/server/internal/domain/dto"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/utils"
)

// BulkCreateExpenses creates all items (and their tag links) in one
// transaction, so a receipt import either fully succeeds or leaves no trace.
func (es *ExpenseService) BulkCreateExpenses(ctx context.Context, userID int64, items []dto.BulkExpenseItem) ([]models.Expense, error) {
	return psql.RunInTx(ctx, es.Pool, func(q *psql.Queries) ([]models.Expense, error) {
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
			row, err := q.CreateExpense(ctx, psql.CreateExpenseParams{
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
				tag, err := q.UpsertTag(ctx, psql.UpsertTagParams{
					Name:   item.Tag,
					UserID: userID,
				})
				if err != nil {
					return nil, err
				}
				if err := q.AddTagToExpense(ctx, psql.AddTagToExpenseParams{
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
		return created, nil
	})
}
