package request

type RecurringScheduleRequest struct {
	Interval   string `json:"interval" binding:"required"`
	DayOfMonth int    `json:"dayOfMonth" binding:"required"`
	StartDate  string `json:"startDate" binding:"required"`
}

type CreateExpenseRequest struct {
	Name              string                    `json:"name" binding:"required"`
	Amount            *float32                  `json:"amount" binding:"required"`
	Cost              *float32                  `json:"cost" binding:"required"`
	RecurringSchedule *RecurringScheduleRequest `json:"recurringSchedule,omitempty"`
}

type BulkExpenseItemRequest struct {
	Name   string   `json:"name" binding:"required"`
	Amount *float32 `json:"amount" binding:"required"`
	Cost   *float32 `json:"cost" binding:"required"`
	Tag    string   `json:"tag"`
}

type BulkCreateExpensesRequest struct {
	Expenses []BulkExpenseItemRequest `json:"expenses" binding:"required,min=1,max=200,dive"`
}

type UpdateExpenseRequest struct {
	Name              string                    `json:"name" binding:"required"`
	Amount            *float32                  `json:"amount" binding:"required"`
	Cost              *float32                  `json:"cost" binding:"required"`
	RecurringSchedule *RecurringScheduleRequest `json:"recurringSchedule,omitempty"`
}
