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
