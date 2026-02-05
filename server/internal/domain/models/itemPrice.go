package models

type ItemPrice struct {
	Id        int64   `json:"id"`
	Price     float32 `json:"price" faker:"oneof: 10.0, 20.0, 30.0, 40.0, 50.0"`
	ItemID    int64   `json:"item_id" faker:"oneof: 1, 2, 3, 4, 5"`
	IsDeleted bool    `json:"is_deleted"`
	CreatedAt string  `json:"created_at"`
}
