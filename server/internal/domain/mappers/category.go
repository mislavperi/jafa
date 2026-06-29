package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type CategoryMapper struct {
}

func NewCategoryMapper() *CategoryMapper {
	return &CategoryMapper{}
}

func (cm *CategoryMapper) MapToDomain(category psql.Category) (models.Category, error) {
	budget, err := category.Budget.Float64Value()
	if err != nil {
		return models.Category{}, err
	}
	keywords := category.Keywords
	if keywords == nil {
		keywords = []string{}
	}
	return models.Category{
		ID:       category.ID,
		Name:     category.Name,
		Icon:     category.Icon,
		Color:    category.Color,
		Budget:   float32(budget.Float64),
		Keywords: keywords,
	}, nil
}

func (cm *CategoryMapper) MapManyToDomain(categories []psql.Category) ([]models.Category, error) {
	return mapSlice(categories, cm.MapToDomain)
}
