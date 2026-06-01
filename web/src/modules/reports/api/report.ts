import type { CategoryBreakdown, MonthlySpend } from '../models/report'
import { apiFetch } from '@/core/api'

export async function getCategoryBreakdown(): Promise<CategoryBreakdown[]> {
  const response = await apiFetch('/api/reports/category-breakdown')
  if (!response.ok) {
    throw new Error('Failed to fetch category breakdown')
  }
  return response.json()
}

export async function getMonthlySpend(): Promise<MonthlySpend[]> {
  const response = await apiFetch('/api/reports/monthly-spend')
  if (!response.ok) {
    throw new Error('Failed to fetch monthly spend')
  }
  return response.json()
}
