import { describe, it, expect, vi, beforeEach } from 'vitest'
import type { Mock } from 'vitest'

vi.mock('@/core/api', () => ({ apiFetch: vi.fn() }))

import { apiFetch } from '@/core/api'
import { getAllTags, createTag, getTagsForExpense, addTagToExpense, removeTagFromExpense } from '@/modules/expense/api/tag'

const mockFetch = apiFetch as Mock

function ok(data: unknown): Response {
  return { ok: true, json: () => Promise.resolve(data) } as unknown as Response
}
const fail: Response = { ok: false } as unknown as Response

beforeEach(() => vi.clearAllMocks())

describe('getAllTags', () => {
  it('returns tags on success', async () => {
    const tags = [{ id: 1, name: 'food' }]
    mockFetch.mockResolvedValue(ok(tags))
    expect(await getAllTags()).toEqual(tags)
  })

  it('throws on failure', async () => {
    mockFetch.mockResolvedValue(fail)
    await expect(getAllTags()).rejects.toThrow('Failed to fetch tags')
  })
})

describe('createTag', () => {
  it('posts and returns new tag', async () => {
    const tag = { id: 2, name: 'transport' }
    mockFetch.mockResolvedValue(ok(tag))
    const result = await createTag('transport')
    expect(result).toEqual(tag)
    expect(mockFetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({ method: 'POST' }),
    )
  })

  it('throws on failure', async () => {
    mockFetch.mockResolvedValue(fail)
    await expect(createTag('x')).rejects.toThrow('Failed to create tag')
  })
})

describe('getTagsForExpense', () => {
  it('returns tags for a given expense', async () => {
    const tags = [{ id: 3, name: 'dining' }]
    mockFetch.mockResolvedValue(ok(tags))
    expect(await getTagsForExpense(10)).toEqual(tags)
    expect(mockFetch).toHaveBeenCalledWith(expect.stringContaining('10'))
  })

  it('throws on failure', async () => {
    mockFetch.mockResolvedValue(fail)
    await expect(getTagsForExpense(10)).rejects.toThrow('Failed to fetch expense tags')
  })
})

describe('addTagToExpense', () => {
  it('resolves without value on success', async () => {
    mockFetch.mockResolvedValue({ ok: true } as Response)
    await expect(addTagToExpense(1, 2)).resolves.toBeUndefined()
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('1'),
      expect.objectContaining({ method: 'POST' }),
    )
  })

  it('throws on failure', async () => {
    mockFetch.mockResolvedValue(fail)
    await expect(addTagToExpense(1, 2)).rejects.toThrow('Failed to add tag to expense')
  })
})

describe('removeTagFromExpense', () => {
  it('resolves without value on success', async () => {
    mockFetch.mockResolvedValue({ ok: true } as Response)
    await expect(removeTagFromExpense(1, 2)).resolves.toBeUndefined()
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/1/tags/2'),
      expect.objectContaining({ method: 'DELETE' }),
    )
  })

  it('throws on failure', async () => {
    mockFetch.mockResolvedValue(fail)
    await expect(removeTagFromExpense(1, 2)).rejects.toThrow('Failed to remove tag from expense')
  })
})
