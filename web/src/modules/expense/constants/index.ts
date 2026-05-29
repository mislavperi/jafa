// Shared palette used to colour tags and expense breakdown charts across the
// expense module (tag chips, the add/edit modal, the breakdown doughnut, …).
export const TAG_COLORS = [
  '#f5c518',
  '#f97316',
  '#22c55e',
  '#3b82f6',
  '#a855f7',
  '#ec4899',
  '#14b8a6',
  '#ef4444',
] as const

// Palette for aggregated/grouped chart series. Ordered for good adjacent
// contrast in pie/doughnut slices.
export const CHART_COLORS = [
  '#f5c518',
  '#f97316',
  '#3b82f6',
  '#a855f7',
  '#ec4899',
  '#14b8a6',
  '#ef4444',
  '#71717a',
] as const

// Resolve a stable colour for a tag from its id.
export function tagColor(tagId: number): string {
  return TAG_COLORS[tagId % TAG_COLORS.length]!
}
