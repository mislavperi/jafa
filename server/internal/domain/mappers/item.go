package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type ItemMapper struct {
}

func NewItemMapper() *ItemMapper {
	return &ItemMapper{}
}

func (im *ItemMapper) MapToDomain(item psql.Item) (models.Item, error) {
	return models.Item{
		Id:        item.ID,
		Name:      item.Name,
		IsDeleted: item.IsDeleted,
		CreatedAt: item.CreatedAt.Time.String(),
		UpdatedAt: item.UpdatedAt.Time.String(),
	}, nil
}

func (im *ItemMapper) MapManyToDomain(items []psql.Item) ([]models.Item, error) {
	mappedItems := make([]models.Item, 0, len(items))
	for _, item := range items {
		mappedItem, err := im.MapToDomain(item)
		if err != nil {
			return []models.Item{}, err
		}
		mappedItems = append(mappedItems, mappedItem)
	}
	return mappedItems, nil
}
