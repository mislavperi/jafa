package models

type RecurrenceInterval string

const (
	RecurrenceIntervalMonthly RecurrenceInterval = "monthly"
	RecurrenceIntervalYearly  RecurrenceInterval = "yearly"
)

type RecurringSchedule struct {
	Interval   RecurrenceInterval `json:"interval"`
	DayOfMonth int                `json:"dayOfMonth"`
	StartDate  string             `json:"startDate"`
}

type Expense struct {
	Id                int64              `json:"id"`
	Name              string             `json:"name" faker:"name"`
	Amount            float32            `json:"amount" faker:"oneof: 1.0, 2.0, 3.0, 4.0, 5.0"`
	Cost              float32            `json:"cost" faker:"oneof: 10.0, 20.0, 30.0, 40.0, 50.0"`
	ItemID            int64              `json:"item_id" faker:"oneof: 1, 2, 3, 4, 5"`
	IsDeleted         bool               `json:"is_deleted"`
	CreatedAt         string             `json:"created_at"`
	UpdatedAt         string             `json:"updated_at"`
	RecurringSchedule *RecurringSchedule `json:"recurringSchedule,omitempty"`
}
