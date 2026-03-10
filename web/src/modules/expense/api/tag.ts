import type { Tag } from '../models/expense'

export async function getAllTags(): Promise<Tag[]> {
  const response = await fetch('/api/tags/')
  if (!response.ok) throw new Error('Failed to fetch tags')
  return response.json()
}

export async function createTag(name: string): Promise<Tag> {
  const response = await fetch('/api/tags/', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name }),
  })
  if (!response.ok) throw new Error('Failed to create tag')
  return response.json()
}

export async function getTagsForExpense(expenseId: number): Promise<Tag[]> {
  const response = await fetch(`/api/expense/${expenseId}/tags`)
  if (!response.ok) throw new Error('Failed to fetch expense tags')
  return response.json()
}

export async function addTagToExpense(expenseId: number, tagId: number): Promise<void> {
  const response = await fetch(`/api/expense/${expenseId}/tags`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ tag_id: tagId }),
  })
  if (!response.ok) throw new Error('Failed to add tag to expense')
}

export async function removeTagFromExpense(expenseId: number, tagId: number): Promise<void> {
  const response = await fetch(`/api/expense/${expenseId}/tags/${tagId}`, {
    method: 'DELETE',
  })
  if (!response.ok) throw new Error('Failed to remove tag from expense')
}
