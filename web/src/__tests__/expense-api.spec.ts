import { describe, it, expect, vi, beforeEach } from 'vitest'
import type { Mock } from 'vitest'

vi.mock('@/core/api', () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from '@/core/api'
import {
  getAllExpenses,
  getExpenseById,
  createExpense,
  updateExpense,
  deleteExpense,
  getMonthlyTotal,
  getDailySpend,
  bulkCreateExpenses,
} from '@/modules/expense/api/expense'

const mockFetch = apiFetch as Mock

function okResponse(data: unknown): Response {
  return {
    ok: true,
    json: () => Promise.resolve(data),
  } as unknown as Response
}

function failResponse(): Response {
  return { ok: false } as unknown as Response
}

describe('expense API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAllExpenses', () => {
    it('returns parsed expenses on success', async () => {
      const expenses = [{ id: 1, name: 'Coffee', amount: 3.5 }]
      mockFetch.mockResolvedValue(okResponse(expenses))

      const result = await getAllExpenses()
      expect(result).toEqual(expenses)
    })

    it('throws on non-ok response', async () => {
      mockFetch.mockResolvedValue(failResponse())
      await expect(getAllExpenses()).rejects.toThrow('Failed to fetch expenses')
    })
  })

  describe('getExpenseById', () => {
    it('returns expense on success', async () => {
      const expense = { id: 7, name: 'Lunch', amount: 12.5 }
      mockFetch.mockResolvedValue(okResponse(expense))

      const result = await getExpenseById(7)
      expect(result).toEqual(expense)
      expect(mockFetch).toHaveBeenCalledWith(expect.stringContaining('7'))
    })

    it('throws on non-ok response', async () => {
      mockFetch.mockResolvedValue(failResponse())
      await expect(getExpenseById(7)).rejects.toThrow('Failed to fetch expense 7')
    })
  })

  describe('createExpense', () => {
    it('posts payload and returns new expense', async () => {
      const created = { id: 10, name: 'Dinner', amount: 25 }
      mockFetch.mockResolvedValue(okResponse(created))

      const result = await createExpense({ name: 'Dinner', amount: 25, cost: 25 })
      expect(result).toEqual(created)
      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({ method: 'POST' }),
      )
    })

    it('throws on non-ok response', async () => {
      mockFetch.mockResolvedValue(failResponse())
      await expect(createExpense({ name: 'X', amount: 1, cost: 1 })).rejects.toThrow(
        'Failed to create expense',
      )
    })
  })

  describe('updateExpense', () => {
    it('patches and returns updated expense', async () => {
      const updated = { id: 5, name: 'Updated', amount: 10 }
      mockFetch.mockResolvedValue(okResponse(updated))

      const result = await updateExpense(5, { name: 'Updated', amount: 10, cost: 10 })
      expect(result).toEqual(updated)
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('5'),
        expect.objectContaining({ method: 'PATCH' }),
      )
    })

    it('throws on non-ok response', async () => {
      mockFetch.mockResolvedValue(failResponse())
      await expect(updateExpense(5, { name: 'X', amount: 1, cost: 1 })).rejects.toThrow(
        'Failed to update expense',
      )
    })
  })

  describe('deleteExpense', () => {
    it('calls DELETE and resolves on success', async () => {
      mockFetch.mockResolvedValue({ ok: true } as Response)

      await expect(deleteExpense(3)).resolves.toBeUndefined()
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('3'),
        expect.objectContaining({ method: 'DELETE' }),
      )
    })

    it('throws on non-ok response', async () => {
      mockFetch.mockResolvedValue(failResponse())
      await expect(deleteExpense(3)).rejects.toThrow('Failed to delete expense')
    })
  })

  describe('getMonthlyTotal', () => {
    it('returns total on success', async () => {
      const total = { total: 450.75 }
      mockFetch.mockResolvedValue(okResponse(total))

      const result = await getMonthlyTotal()
      expect(result).toEqual(total)
    })

    it('throws on non-ok response', async () => {
      mockFetch.mockResolvedValue(failResponse())
      await expect(getMonthlyTotal()).rejects.toThrow('Failed to fetch monthly total')
    })
  })

  describe('getDailySpend', () => {
    it('passes months param and returns data', async () => {
      const rows = [{ day: '2024-06-01', total: 30 }]
      mockFetch.mockResolvedValue(okResponse(rows))

      const result = await getDailySpend(3)
      expect(result).toEqual(rows)
      expect(mockFetch).toHaveBeenCalledWith(expect.stringContaining('months=3'))
    })

    it('throws on non-ok response', async () => {
      mockFetch.mockResolvedValue(failResponse())
      await expect(getDailySpend(1)).rejects.toThrow('Failed to fetch daily spend')
    })
  })

  describe('bulkCreateExpenses', () => {
    it('posts items and returns created expenses', async () => {
      const created = [
        { id: 1, name: 'Milk', amount: 3 },
        { id: 2, name: 'Bread', amount: 2 },
      ]
      mockFetch.mockResolvedValue(okResponse(created))

      const result = await bulkCreateExpenses([
        { name: 'Milk', amount: 3, cost: 3 },
        { name: 'Bread', amount: 2, cost: 2 },
      ])
      expect(result).toEqual(created)
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('bulk'),
        expect.objectContaining({ method: 'POST' }),
      )
    })

    it('throws on non-ok response', async () => {
      mockFetch.mockResolvedValue(failResponse())
      await expect(bulkCreateExpenses([])).rejects.toThrow('Failed to import expenses')
    })
  })
})
