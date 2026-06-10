import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { getAllExpenses, getExpenseById, getMonthlyTotal, getDailySpend, getFirstExpenseDate, getDailySpendForMonth, getExpensesByMonth, createExpense, bulkCreateExpenses, updateExpense, deleteExpense } from '../api/expense'
import type { RecurringSchedule } from '../models/expense'
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

export function useFirstExpenseDate() {
  return useQuery({
    queryKey: ['expenses', 'first-expense-date'],
    queryFn: getFirstExpenseDate,
  })
}

export function useDailySpendForMonth(year: Ref<number>, month: Ref<number>) {
  return useQuery({
    queryKey: ['expenses', 'daily-spend-for-month', year, month],
    queryFn: () => getDailySpendForMonth(year.value, month.value),
  })
}

export function useExpensesByMonth(year: Ref<number>, month: Ref<number>) {
  return useQuery({
    queryKey: ['expenses', 'by-month', year, month],
    queryFn: () => getExpensesByMonth(year.value, month.value),
  })
}

export function useCreateExpense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: createExpense,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
    },
  })
}

export function useBulkCreateExpenses() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: bulkCreateExpenses,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      queryClient.invalidateQueries({ queryKey: ['tags'] })
    },
  })
}

export function useUpdateExpense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: { id: number; name: string; amount: number; cost: number; recurringSchedule?: RecurringSchedule }) =>
      updateExpense(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
    },
  })
}

export function useDeleteExpense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: deleteExpense,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
    },
  })
}
