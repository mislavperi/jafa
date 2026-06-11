package services

import (
	"context"
	"strings"

	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type CategoryService struct {
	Queries *psql.Queries
	Mapper  *mappers.CategoryMapper
}

func NewCategoryService(queries *psql.Queries) *CategoryService {
	return &CategoryService{
		Queries: queries,
		Mapper:  mappers.NewCategoryMapper(),
	}
}

func (cs *CategoryService) List(ctx context.Context) ([]models.Category, error) {
	categories, err := cs.Queries.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	return cs.Mapper.MapManyToDomain(categories)
}

// Categorize resolves the category name for an expense by matching its name
// against each category's keywords. The fallback ("Other") category has no
// keywords and is used when nothing else matches.
func Categorize(name string, categories []models.Category) string {
	lower := strings.ToLower(name)
	for _, category := range categories {
		for _, keyword := range category.Keywords {
			if keyword != "" && strings.Contains(lower, keyword) {
				return category.Name
			}
		}
	}
	return "Other"
}
