package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/utils"
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
		CreatedAt: utils.FormatRFC3339(tag.CreatedAt),
		UpdatedAt: utils.FormatRFC3339(tag.UpdatedAt),
	}, nil
}

func (tm *TagMapper) MapManyToDomain(tags []psql.Tag) ([]models.Tag, error) {
	return mapSlice(tags, tm.MapToDomain)
}
