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

func (ts *TagService) GetAllTags() ([]models.Tag, error) {
	tags, err := ts.Queries.GetAllTags(context.Background())
	if err != nil {
		return nil, err
	}
	return ts.Mapper.MapManyToDomain(tags)
}

func (ts *TagService) CreateTag(name string) (models.Tag, error) {
	tag, err := ts.Queries.CreateTag(context.Background(), name)
	if err != nil {
		return models.Tag{}, err
	}
	return ts.Mapper.MapToDomain(tag)
}

func (ts *TagService) GetTagsForExpense(expenseID int64) ([]models.Tag, error) {
	tags, err := ts.Queries.GetTagsForExpense(context.Background(), expenseID)
	if err != nil {
		return nil, err
	}
	return ts.Mapper.MapManyToDomain(tags)
}

func (ts *TagService) AddTagToExpense(expenseID, tagID int64) error {
	return ts.Queries.AddTagToExpense(context.Background(), psql.AddTagToExpenseParams{
		ExpenseID: expenseID,
		TagID:     tagID,
	})
}

func (ts *TagService) RemoveTagFromExpense(expenseID, tagID int64) error {
	return ts.Queries.RemoveTagFromExpense(context.Background(), psql.RemoveTagFromExpenseParams{
		ExpenseID: expenseID,
		TagID:     tagID,
	})
}
