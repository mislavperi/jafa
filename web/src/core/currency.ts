// Currency formatting shared across the app. The active currency lives in the
// user's preferences (see the theme store); these helpers turn a numeric amount
// into a display string using the matching symbol.

export const CURRENCY_SYMBOLS: Record<string, string> = {
  EUR: '€',
  USD: '$',
}

export const CURRENCY_OPTIONS = Object.keys(CURRENCY_SYMBOLS)

export function currencySymbol(currency: string): string {
  return CURRENCY_SYMBOLS[currency] ?? currency
}

export function formatCurrency(value: number, currency: string): string {
  return currencySymbol(currency) + value.toFixed(2)
}
