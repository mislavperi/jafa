// Types for the receipt scanner flow (upload → scan → review).

export type ScanStep = 'upload' | 'scanning' | 'review'

export interface SampleItem {
  name: string
  amount: number
  category: string | null
  confidence: number
}

export interface Sample {
  id: string
  merchant: string
  address: string
  date: string
  total: number
  items: SampleItem[]
}

export interface ReviewItem extends SampleItem {
  id: number
  included: boolean
  suggestedTag: string
  finalTag: string
  needsReview: boolean
}
