<script setup lang="ts">
import { computed } from 'vue'
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import { useExpenses } from '@/modules/expense/composables/useExpenses'
import type { Expense } from '@/modules/expense/models/expense'

const { data: expenses, isLoading } = useExpenses()

interface CategoryConfig {
  name: string
  icon: string
  color: string
  budget: number
  keywords: string[]
}

const CATEGORIES: CategoryConfig[] = [
  { name: 'Groceries', icon: 'pi pi-shopping-cart', color: '#f5c518', budget: 400, keywords: ['grocer', 'supermarket', 'food', 'market'] },
  { name: 'Dining', icon: 'pi pi-utensils', color: '#f97316', budget: 300, keywords: ['restaurant', 'cafe', 'coffee', 'dining', 'eat', 'lunch', 'dinner', 'breakfast'] },
  { name: 'Transport', icon: 'pi pi-car', color: '#3b82f6', budget: 200, keywords: ['uber', 'lyft', 'bus', 'train', 'transit', 'fuel', 'gas', 'transport', 'taxi'] },
  { name: 'Bills', icon: 'pi pi-file-edit', color: '#a855f7', budget: 500, keywords: ['bill', 'electric', 'water', 'internet', 'phone', 'rent', 'insurance'] },
  { name: 'Shopping', icon: 'pi pi-tag', color: '#ec4899', budget: 250, keywords: ['shop', 'amazon', 'clothing', 'clothes', 'shoes', 'mall'] },
  { name: 'Entertainment', icon: 'pi pi-star', color: '#14b8a6', budget: 150, keywords: ['netflix', 'spotify', 'movie', 'game', 'entertain', 'stream'] },
  { name: 'Health', icon: 'pi pi-heart', color: '#ef4444', budget: 100, keywords: ['gym', 'health', 'pharmacy', 'doctor', 'medical', 'fitness'] },
  { name: 'Other', icon: 'pi pi-ellipsis-h', color: '#71717a', budget: 200, keywords: [] },
]

function categorize(expense: Expense): string {
  const name = expense.name.toLowerCase()
  for (const cat of CATEGORIES.slice(0, -1)) {
    if (cat.keywords.some((k) => name.includes(k))) return cat.name
  }
  return 'Other'
}

const categoryStats = computed(() => {
  const totals: Record<string, number> = {}
  for (const e of expenses.value ?? []) {
    const cat = categorize(e)
    totals[cat] = (totals[cat] ?? 0) + (e.amount)
  }
  return CATEGORIES.map((c) => {
    const spent = totals[c.name] ?? 0
    const pct = Math.min(100, Math.round((spent / c.budget) * 100))
    const remaining = c.budget - spent
    return { ...c, spent, pct, remaining }
  })
})

function fmt(val: number) {
  return '€' + Math.abs(val).toFixed(2)
}
</script>

<template>
  <Root>
    <div class="flex flex-col gap-5 h-full min-w-0 p-8 overflow-auto">
      <AppPageHeader title="Categories" subtitle="Budget vs. spending by category" />

      <div v-if="isLoading" class="grid grid-cols-2 gap-4">
        <div v-for="i in 8" :key="i" class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 h-32 animate-pulse" />
      </div>

      <div v-else class="grid grid-cols-2 gap-4">
        <div
          v-for="cat in categoryStats"
          :key="cat.name"
          class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 flex flex-col gap-3"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-10 h-10 rounded-[10px] flex items-center justify-center shrink-0"
              :style="{ backgroundColor: cat.color + '22', color: cat.color }"
            >
              <i :class="cat.icon" class="text-[calc(16px*var(--jafa-text-scale,1))]" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-[var(--jafa-text)] font-semibold text-sm leading-tight">{{ cat.name }}</p>
              <p class="text-[var(--jafa-text-muted)] text-xs mt-0.5">Budget: {{ fmt(cat.budget) }}</p>
            </div>
            <div class="text-right">
              <p class="text-[var(--jafa-text)] font-bold text-base tabular-nums">{{ fmt(cat.spent) }}</p>
              <p class="text-[calc(11px*var(--jafa-text-scale,1))] mt-0.5" :class="cat.remaining >= 0 ? 'text-emerald-400' : 'text-red-400'">
                {{ cat.remaining >= 0 ? fmt(cat.remaining) + ' left' : fmt(cat.remaining) + ' over' }}
              </p>
            </div>
          </div>

          <div class="flex flex-col gap-1">
            <div class="h-1.5 rounded-full bg-[var(--jafa-border)] overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :style="{
                  width: cat.pct + '%',
                  backgroundColor: cat.pct >= 100 ? '#ef4444' : cat.color,
                }"
              />
            </div>
            <p class="text-[calc(11px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)]">{{ cat.pct }}% of budget used</p>
          </div>
        </div>
      </div>
    </div>
  </Root>
</template>
