import { useQuery } from '@tanstack/vue-query'
import { getAllExpenses, getExpenseById, getMonthlyTotal, getDailySpend } from '../api/expense'
import type { Ref } from 'vue'

export function useExpenses() {
  return useQuery({
    queryKey: ['expenses'],
    queryFn: getAllExpenses,
  })
}

export function useExpense(id: Ref<number>) {
  return useQuery({
    queryKey: ['expenses', id],
    queryFn: () => getExpenseById(id.value),
  })
}

export function useMonthlyTotal() {
  return useQuery({
    queryKey: ['expenses', 'monthly-total'],
    queryFn: getMonthlyTotal,
  })
}

export function useDailySpend(months: Ref<number>) {
  return useQuery({
    queryKey: ['expenses', 'daily-spend', months],
    queryFn: () => getDailySpend(months.value),
  })
}
