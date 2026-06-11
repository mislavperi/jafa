package services

import (
	"context"

	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type TagService struct {
	Queries *psql.Queries
	Mapper  *mappers.TagMapper
}

func NewTagService(queries *psql.Queries) *TagService {
	return &TagService{
		Queries: queries,
		Mapper:  mappers.NewTagMapper(),
	}
}

func (ts *TagService) GetAllTags(ctx context.Context, userID int64) ([]models.Tag, error) {
	tags, err := ts.Queries.GetAllTags(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ts.Mapper.MapManyToDomain(tags)
}

func (ts *TagService) CreateTag(ctx context.Context, userID int64, name string) (models.Tag, error) {
	tag, err := ts.Queries.CreateTag(ctx, psql.CreateTagParams{
		Name:   name,
		UserID: userID,
	})
	if err != nil {
		return models.Tag{}, err
	}
	return ts.Mapper.MapToDomain(tag)
}

func (ts *TagService) GetTagsForExpense(ctx context.Context, userID, expenseID int64) ([]models.Tag, error) {
	tags, err := ts.Queries.GetTagsForExpense(ctx, psql.GetTagsForExpenseParams{
		ExpenseID: expenseID,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}
	return ts.Mapper.MapManyToDomain(tags)
}

func (ts *TagService) AddTagToExpense(ctx context.Context, userID, expenseID, tagID int64) error {
	return ts.Queries.AddTagToExpense(ctx, psql.AddTagToExpenseParams{
		ExpenseID: expenseID,
		TagID:     tagID,
		UserID:    userID,
	})
}

func (ts *TagService) RemoveTagFromExpense(ctx context.Context, userID, expenseID, tagID int64) error {
	return ts.Queries.RemoveTagFromExpense(ctx, psql.RemoveTagFromExpenseParams{
		ExpenseID: expenseID,
		TagID:     tagID,
		UserID:    userID,
	})
}
