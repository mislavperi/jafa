import { useQuery } from '@tanstack/vue-query'
import { getAllExpenses, getExpenseById } from '../api/expense'
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
