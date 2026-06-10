import type { Expense, MonthlyTotal, DailySpend, FirstExpenseDate, RecurringSchedule } from '../models/expense'
import { apiFetch } from '@/core/api'

const BASE_URL = '/api/expense/'

export async function createExpense(payload: {
  name: string
  amount: number
  cost: number
  recurringSchedule?: RecurringSchedule
}): Promise<Expense> {
  const response = await apiFetch(BASE_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    throw new Error('Failed to create expense')
  }
  return response.json()
}

export interface BulkExpenseItem {
  name: string
  amount: number
  cost: number
  tag?: string
}

export async function bulkCreateExpenses(items: BulkExpenseItem[]): Promise<Expense[]> {
  const response = await apiFetch(`${BASE_URL}bulk`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ expenses: items }),
  })
  if (!response.ok) {
    throw new Error('Failed to import expenses')
  }
  return response.json()
}

export async function updateExpense(
  id: number,
  payload: {
    name: string
    amount: number
    cost: number
    recurringSchedule?: RecurringSchedule
  },
): Promise<Expense> {
  const response = await apiFetch(`${BASE_URL}${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    throw new Error('Failed to update expense')
  }
  return response.json()
}

export async function deleteExpense(id: number): Promise<void> {
  const response = await apiFetch(`${BASE_URL}${id}`, { method: 'DELETE' })
  if (!response.ok) {
    throw new Error('Failed to delete expense')
  }
}

export async function getAllExpenses(): Promise<Expense[]> {
  const response = await apiFetch(BASE_URL)
  if (!response.ok) {
    throw new Error('Failed to fetch expenses')
  }
  return response.json()
}

export async function getExpenseById(id: number): Promise<Expense> {
  const response = await apiFetch(`${BASE_URL}${id}`)
  if (!response.ok) {
    throw new Error(`Failed to fetch expense ${id}`)
  }
  return response.json()
}

export async function getMonthlyTotal(): Promise<MonthlyTotal> {
  const response = await apiFetch('/api/expense-stats/monthly-total')
  if (!response.ok) {
    throw new Error('Failed to fetch monthly total')
  }
  return response.json()
}

export async function getDailySpend(months: number): Promise<DailySpend[]> {
  const response = await apiFetch(`/api/expense-stats/daily-spend?months=${months}`)
  if (!response.ok) {
    throw new Error('Failed to fetch daily spend')
  }
  return response.json()
}

export async function getExpensesByMonth(year: number, month: number): Promise<Expense[]> {
  const response = await apiFetch(`/api/expense-stats/expenses-by-month?year=${year}&month=${month}`)
  if (!response.ok) {
    throw new Error('Failed to fetch expenses by month')
  }
  return response.json()
}

export async function getFirstExpenseDate(): Promise<FirstExpenseDate> {
  const response = await apiFetch('/api/expense-stats/first-expense-date')
  if (!response.ok) {
    throw new Error('Failed to fetch first expense date')
  }
  return response.json()
}

export async function getDailySpendForMonth(year: number, month: number): Promise<DailySpend[]> {
  const response = await apiFetch(`/api/expense-stats/daily-spend-for-month?year=${year}&month=${month}`)
  if (!response.ok) {
    throw new Error('Failed to fetch daily spend for month')
  }
  return response.json()
}
