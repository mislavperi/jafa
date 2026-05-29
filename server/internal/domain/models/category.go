package models

type Category struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Icon     string   `json:"icon"`
	Color    string   `json:"color"`
	Budget   float32  `json:"budget"`
	Keywords []string `json:"keywords"`
}
