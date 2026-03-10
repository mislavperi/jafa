import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { getAllTags, createTag, getTagsForExpense, addTagToExpense, removeTagFromExpense } from '../api/tag'
import type { Ref } from 'vue'

export function useAllTags() {
  return useQuery({
    queryKey: ['tags'],
    queryFn: getAllTags,
  })
}

export function useCreateTag() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: createTag,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tags'] })
    },
  })
}

export function useExpenseTags(expenseId: Ref<number>) {
  return useQuery({
    queryKey: ['expense-tags', expenseId],
    queryFn: () => getTagsForExpense(expenseId.value),
  })
}

export function useAddTagToExpense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ expenseId, tagId }: { expenseId: number; tagId: number }) =>
      addTagToExpense(expenseId, tagId),
    onSuccess: (_data, { expenseId }) => {
      queryClient.invalidateQueries({ queryKey: ['expense-tags', expenseId] })
    },
  })
}

export function useRemoveTagFromExpense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ expenseId, tagId }: { expenseId: number; tagId: number }) =>
      removeTagFromExpense(expenseId, tagId),
    onSuccess: (_data, { expenseId }) => {
      queryClient.invalidateQueries({ queryKey: ['expense-tags', expenseId] })
    },
  })
}
