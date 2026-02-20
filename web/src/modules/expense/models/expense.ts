export interface Expense {
  id: number
  name: string
  amount: number
}

export interface MonthlyTotal {
  total: number
}

export interface DailySpend {
  day: string
  total: number
}
