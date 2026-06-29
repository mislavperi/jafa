package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type ItemPriceMapper struct {
}

func NewItemPriceMapper() *ItemPriceMapper {
	return &ItemPriceMapper{}
}

func (ipm *ItemPriceMapper) MapToDomain(itemPrice psql.ItemPrice) (models.ItemPrice, error) {
	price, err := itemPrice.Price.Float64Value()
	if err != nil || !price.Valid {
		return models.ItemPrice{}, err
	}
	return models.ItemPrice{
		Id:        itemPrice.ID,
		Price:     float32(price.Float64),
		ItemID:    itemPrice.ItemID,
		IsDeleted: itemPrice.IsDeleted,
		CreatedAt: itemPrice.CreatedAt.Time.String(),
	}, nil
}

func (ipm *ItemPriceMapper) MapManyToDomain(itemPrices []psql.ItemPrice) ([]models.ItemPrice, error) {
	return mapSlice(itemPrices, ipm.MapToDomain)
}
