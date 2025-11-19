package models

import "github.com/mislavperi/jafa/server/utils/enums"

type Expense struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Currency enums.Currency `json:"currency"`
	Type     ExpenseType    `json:"type"`
	Amount   float32        `json:"amount"`
}
