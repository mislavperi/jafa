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

export interface InstallmentPlan {
  count: number
  paymentAmount: number
}

export type ExpenseKind = 'expense' | 'income'

export interface Expense {
  id: number
  kind?: ExpenseKind
  name: string
  amount: number
  cost?: number
  recurringSchedule?: RecurringSchedule
  installmentPlan?: InstallmentPlan
  created_at?: string
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
