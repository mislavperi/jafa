import type { Expense } from '../models/expense'

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
