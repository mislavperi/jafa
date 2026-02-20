import type { Expense, MonthlyTotal, DailySpend } from '../models/expense'

const BASE_URL = '/api/expense/'

export async function getAllExpenses(): Promise<Expense[]> {
  const response = await fetch(BASE_URL)
  if (!response.ok) {
    throw new Error('Failed to fetch expenses')
  }
  return response.json()
}

export async function getExpenseById(id: number): Promise<Expense> {
  const response = await fetch(`${BASE_URL}/${id}`)
  if (!response.ok) {
    throw new Error(`Failed to fetch expense ${id}`)
  }
  return response.json()
}

export async function getMonthlyTotal(): Promise<MonthlyTotal> {
  const response = await fetch('/api/expense-stats/monthly-total')
  if (!response.ok) {
    throw new Error('Failed to fetch monthly total')
  }
  return response.json()
}

export async function getDailySpend(months: number): Promise<DailySpend[]> {
  const response = await fetch(`/api/expense-stats/daily-spend?months=${months}`)
  if (!response.ok) {
    throw new Error('Failed to fetch daily spend')
  }
  return response.json()
}
