// Report shapes returned by the backend /reports endpoints. The backend owns
// the category taxonomy and all aggregation, so the frontend only renders.

export interface CategoryBreakdown {
  name: string
  icon: string
  color: string
  budget: number
  spent: number
  remaining: number
  pct: number
}

export interface MonthlySpend {
  month: string // YYYY-MM
  total: number
}
