import type { Sample } from '../models/receiptScan'

// Demo receipts used by the scanner's "try a sample" affordance until a real
// OCR backend is wired up.
export const SAMPLES: Sample[] = [
  {
    id: 'whole-foods',
    merchant: 'WHOLE FOODS MARKET',
    address: '1234 Market St · San Francisco, CA',
    date: '2026-05-26',
    total: 87.42,
    items: [
      { name: 'Organic Bananas', amount: 3.49, category: 'groceries', confidence: 0.95 },
      { name: 'Almond Milk · 64oz', amount: 5.99, category: 'groceries', confidence: 0.92 },
      { name: 'Sourdough Loaf', amount: 6.5, category: 'groceries', confidence: 0.9 },
      { name: 'Avocados · 4ct', amount: 7.96, category: 'groceries', confidence: 0.94 },
      { name: 'Chicken Breast · 1.2lb', amount: 14.38, category: 'groceries', confidence: 0.88 },
      { name: 'Misc HBC', amount: 12.99, category: null, confidence: 0.32 },
      { name: 'Olive Oil EVOO', amount: 18.49, category: 'groceries', confidence: 0.86 },
      { name: 'Spinach · Organic', amount: 4.29, category: 'groceries', confidence: 0.93 },
      { name: 'Greek Yogurt', amount: 6.99, category: 'groceries', confidence: 0.91 },
      { name: 'Sales Tax', amount: 6.34, category: null, confidence: 0.25 },
    ],
  },
  {
    id: 'pharmacy',
    merchant: 'CVS PHARMACY',
    address: '500 Main St · Brooklyn, NY',
    date: '2026-05-25',
    total: 42.18,
    items: [
      { name: 'Ibuprofen 200mg · 100ct', amount: 8.99, category: 'health', confidence: 0.93 },
      { name: 'Multivitamin · 90ct', amount: 14.49, category: 'health', confidence: 0.91 },
      { name: 'Toothpaste', amount: 4.79, category: 'shopping', confidence: 0.65 },
      { name: 'Hand Sanitizer', amount: 3.49, category: 'health', confidence: 0.78 },
      { name: 'Snack Bar', amount: 1.99, category: null, confidence: 0.42 },
      { name: 'Allergy Tablets', amount: 10.43, category: 'health', confidence: 0.89 },
    ],
  },
  {
    id: 'restaurant',
    merchant: 'TAQUERIA EL FAROLITO',
    address: '2779 Mission St',
    date: '2026-05-24',
    total: 38.5,
    items: [
      { name: 'Super Burrito · Carnitas', amount: 13.5, category: 'dining', confidence: 0.94 },
      { name: 'Chips & Guacamole', amount: 7.25, category: 'dining', confidence: 0.9 },
      { name: 'Horchata · Large', amount: 5.0, category: 'dining', confidence: 0.86 },
      { name: 'Tacos al Pastor · 3pc', amount: 9.75, category: 'dining', confidence: 0.92 },
      { name: 'Tip · 18%', amount: 3.0, category: null, confidence: 0.3 },
    ],
  },
]

// Keyword hints used to guess a category for an item the scanner couldn't
// confidently classify.
export const CATEGORY_HINTS: Record<string, string[]> = {
  groceries: ['milk', 'bread', 'banana', 'avocado', 'fruit', 'vegetable', 'spinach', 'yogurt', 'oil', 'pasta', 'rice', 'chicken', 'beef', 'cheese', 'coffee', 'tea'],
  dining: ['burrito', 'taco', 'pizza', 'sushi', 'tip', 'meal', 'lunch', 'dinner', 'sandwich', 'horchata', 'drink', 'soda', 'chips'],
  health: ['vitamin', 'tablet', 'pill', 'ibuprofen', 'medicine', 'pharmacy', 'allergy', 'sanitizer', 'bandage'],
  shopping: ['toothpaste', 'soap', 'shampoo', 'detergent', 'paper', 'hbc'],
  transport: ['gas', 'petrol', 'uber', 'lyft', 'fare', 'parking'],
  bills: ['internet', 'electric', 'water', 'phone', 'subscription'],
  entertainment: ['movie', 'cinema', 'concert', 'spotify', 'netflix'],
}

export function guessCategory(name: string): string {
  const lower = name.toLowerCase()
  for (const [cat, hints] of Object.entries(CATEGORY_HINTS)) {
    if (hints.some((h) => lower.includes(h))) return cat
  }
  return 'other'
}
