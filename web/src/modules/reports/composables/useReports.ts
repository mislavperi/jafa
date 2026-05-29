import { useQuery } from '@tanstack/vue-query'
import { getCategoryBreakdown, getMonthlySpend } from '../api/report'

export function useCategoryBreakdown() {
  return useQuery({
    queryKey: ['reports', 'category-breakdown'],
    queryFn: getCategoryBreakdown,
  })
}

export function useMonthlySpend() {
  return useQuery({
    queryKey: ['reports', 'monthly-spend'],
    queryFn: getMonthlySpend,
  })
}
