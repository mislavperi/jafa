import type { Expense, MonthlyTotal, DailySpend, FirstExpenseDate } from '../models/expense'

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

export async function getExpensesByMonth(year: number, month: number): Promise<Expense[]> {
  const response = await fetch(`/api/expense-stats/expenses-by-month?year=${year}&month=${month}`)
  if (!response.ok) {
    throw new Error('Failed to fetch expenses by month')
  }
  return response.json()
}

export async function getFirstExpenseDate(): Promise<FirstExpenseDate> {
  const response = await fetch('/api/expense-stats/first-expense-date')
  if (!response.ok) {
    throw new Error('Failed to fetch first expense date')
  }
  return response.json()
}

export async function getDailySpendForMonth(year: number, month: number): Promise<DailySpend[]> {
  const response = await fetch(`/api/expense-stats/daily-spend-for-month?year=${year}&month=${month}`)
  if (!response.ok) {
    throw new Error('Failed to fetch daily spend for month')
  }
  return response.json()
}
