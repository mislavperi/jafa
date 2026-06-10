import { describe, it, expect } from 'vitest'
import { currencySymbol, formatCurrency, CURRENCY_SYMBOLS } from '@/core/currency'

describe('currencySymbol', () => {
  it('returns € for EUR', () => {
    expect(currencySymbol('EUR')).toBe('€')
  })

  it('returns $ for USD', () => {
    expect(currencySymbol('USD')).toBe('$')
  })

  it('returns the raw code for unknown currencies', () => {
    expect(currencySymbol('GBP')).toBe('GBP')
    expect(currencySymbol('JPY')).toBe('JPY')
  })
})

describe('formatCurrency', () => {
  it('prefixes with symbol and formats to 2 decimal places', () => {
    expect(formatCurrency(42, 'EUR')).toBe('€42.00')
    expect(formatCurrency(9.5, 'USD')).toBe('$9.50')
  })

  it('handles zero', () => {
    expect(formatCurrency(0, 'EUR')).toBe('€0.00')
  })

  it('falls back to raw code for unknown currency', () => {
    expect(formatCurrency(10, 'GBP')).toBe('GBP10.00')
  })
})

describe('CURRENCY_SYMBOLS', () => {
  it('includes EUR and USD', () => {
    expect(CURRENCY_SYMBOLS).toHaveProperty('EUR')
    expect(CURRENCY_SYMBOLS).toHaveProperty('USD')
  })
})
