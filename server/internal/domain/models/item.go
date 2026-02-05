package models

type Item struct {
	Id        int64  `json:"id"`
	Name      string `json:"name" faker:"word"`
	IsDeleted bool   `json:"is_deleted"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
