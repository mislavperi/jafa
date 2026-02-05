package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type TagMapper struct {
}

func NewTagMapper() *TagMapper {
	return &TagMapper{}
}

func (tm *TagMapper) MapToDomain(tag psql.Tag) (models.Tag, error) {
	return models.Tag{
		Id:        tag.ID,
		Name:      tag.Name,
		IsDeleted: tag.IsDeleted,
		CreatedAt: tag.CreatedAt.Time.String(),
		UpdatedAt: tag.UpdatedAt.Time.String(),
	}, nil
}

func (tm *TagMapper) MapManyToDomain(tags []psql.Tag) ([]models.Tag, error) {
	mappedTags := make([]models.Tag, 0, len(tags))
	for _, tag := range tags {
		mappedTag, err := tm.MapToDomain(tag)
		if err != nil {
			return []models.Tag{}, err
		}
		mappedTags = append(mappedTags, mappedTag)
	}
	return mappedTags, nil
}
