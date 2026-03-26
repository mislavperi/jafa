export interface Tag {
  id: number
  name: string
}

export type RecurrenceInterval = 'monthly' | 'yearly'

export interface RecurringSchedule {
  interval: RecurrenceInterval
  dayOfMonth: number
  startDate: string
}

export interface Expense {
  id: number
  name: string
  amount: number
  cost?: number
  recurringSchedule?: RecurringSchedule
}

export interface MonthlyTotal {
  total: number
}

export interface DailySpend {
  day: string
  total: number
}

export interface FirstExpenseDate {
  firstDate: string
}
