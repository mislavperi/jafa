// Package dto holds the command/input objects passed across the
// controller→service boundary. These are application-layer DTOs, distinct from
// the domain entities in models (which represent business state) and from the
// wire/transport DTOs in models/request (which carry JSON binding tags). A dto
// here has no identity and no tags; the controller builds one from a request
// and hands it to a service.
package dto

import "github.com/mislavperi/jafa/server/internal/domain/models"

// CreateExpenseInput is the command ExpenseService.CreateExpense needs to record
// a new entry. RecurringSchedule and InstallmentCount are optional.
type CreateExpenseInput struct {
	UserID            int64
	Kind              string
	Name              string
	Amount            float32
	Cost              float32
	RecurringSchedule *models.RecurringSchedule
	InstallmentCount  *int
}

// UpdateExpenseInput is the command ExpenseService.UpdateExpense needs. Kind is
// not editable after creation, so it is absent here.
type UpdateExpenseInput struct {
	ID                int64
	UserID            int64
	Name              string
	Amount            float32
	Cost              float32
	RecurringSchedule *models.RecurringSchedule
	InstallmentCount  *int
}

// BulkExpenseItem is one expense in a bulk import (e.g. a receipt line item).
// Tag is optional; when set, the tag is created if it does not exist yet and
// linked to the new expense.
type BulkExpenseItem struct {
	Name   string
	Amount float32
	Cost   float32
	Tag    string
}
