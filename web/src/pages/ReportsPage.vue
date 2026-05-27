<script setup lang="ts">
import { computed } from 'vue'
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import AppStatCard from '@/core/components/AppStatCard.vue'
import Chart from 'primevue/chart'
import { useExpenses } from '@/modules/expense/composables/useExpenses'
import type { Expense } from '@/modules/expense/models/expense'

const { data: expenses, isLoading } = useExpenses()

function fmt(val: number) {
  return '$' + val.toFixed(2)
}

const monthlyTotals = computed(() => {
  const map: Record<string, number> = {}
  for (const e of expenses.value ?? []) {
    const date = e.created_at ? new Date(e.created_at) : null
    if (!date) continue
    const key = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
    map[key] = (map[key] ?? 0) + (e.cost ?? 0)
  }
  return map
})

const sortedMonths = computed(() => Object.keys(monthlyTotals.value).sort())

const thisMonthKey = computed(() => {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
})

const thisMonthTotal = computed(() => monthlyTotals.value[thisMonthKey.value] ?? 0)

const avg12Month = computed(() => {
  const months = sortedMonths.value.slice(-12)
  if (!months.length) return 0
  const total = months.reduce((s, m) => s + (monthlyTotals.value[m] ?? 0), 0)
  return total / months.length
})

const highestMonth = computed(() => {
  let max = 0
  for (const v of Object.values(monthlyTotals.value)) {
    if (v > max) max = v
  }
  return max
})

const chartData = computed(() => {
  const labels = sortedMonths.value.slice(-12).map((k) => {
    const [y, m] = k.split('-')
    return new Date(Number(y), Number(m) - 1).toLocaleString('default', { month: 'short', year: '2-digit' })
  })
  const data = sortedMonths.value.slice(-12).map((m) => monthlyTotals.value[m] ?? 0)
  return {
    labels,
    datasets: [
      {
        label: 'Monthly Spend',
        data,
        backgroundColor: '#f5c518cc',
        borderColor: '#f5c518',
        borderWidth: 1,
        borderRadius: 6,
      },
    ],
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        label: (ctx: { parsed: { y: number } }) => '$' + ctx.parsed.y.toFixed(2),
      },
    },
  },
  scales: {
    x: {
      ticks: { color: '#71717a' },
      grid: { color: '#26262c' },
    },
    y: {
      ticks: {
        color: '#71717a',
        callback: (v: number) => '$' + v,
      },
      grid: { color: '#26262c' },
    },
  },
}

const CATEGORIES = [
  { name: 'Groceries', color: '#f5c518', keywords: ['grocer', 'supermarket', 'food', 'market'] },
  { name: 'Dining', color: '#f97316', keywords: ['restaurant', 'cafe', 'coffee', 'dining', 'eat', 'lunch', 'dinner', 'breakfast'] },
  { name: 'Transport', color: '#3b82f6', keywords: ['uber', 'lyft', 'bus', 'train', 'transit', 'fuel', 'gas', 'transport', 'taxi'] },
  { name: 'Bills', color: '#a855f7', keywords: ['bill', 'electric', 'water', 'internet', 'phone', 'rent', 'insurance'] },
  { name: 'Shopping', color: '#ec4899', keywords: ['shop', 'amazon', 'clothing', 'clothes', 'shoes', 'mall'] },
  { name: 'Entertainment', color: '#14b8a6', keywords: ['netflix', 'spotify', 'movie', 'game', 'entertain', 'stream'] },
  { name: 'Health', color: '#ef4444', keywords: ['gym', 'health', 'pharmacy', 'doctor', 'medical', 'fitness'] },
  { name: 'Other', color: '#71717a', keywords: [] },
]

function categorize(expense: Expense): string {
  const name = expense.name.toLowerCase()
  for (const cat of CATEGORIES.slice(0, -1)) {
    if (cat.keywords.some((k) => name.includes(k))) return cat.name
  }
  return 'Other'
}

const categoryAverages = computed(() => {
  const totals: Record<string, number> = {}
  for (const e of expenses.value ?? []) {
    const cat = categorize(e)
    totals[cat] = (totals[cat] ?? 0) + (e.cost ?? 0)
  }
  const maxVal = Math.max(1, ...Object.values(totals))
  return CATEGORIES.map((c) => ({
    name: c.name,
    color: c.color,
    total: totals[c.name] ?? 0,
    pct: Math.round(((totals[c.name] ?? 0) / maxVal) * 100),
  })).filter((c) => c.total > 0)
})
</script>

<template>
  <Root>
    <div class="flex flex-col gap-5 h-full min-w-0 p-8 overflow-auto">
      <AppPageHeader title="Reports" subtitle="Spending insights and trends" />

      <div class="grid grid-cols-3 gap-4">
        <AppStatCard
          label="This Month"
          :value="fmt(thisMonthTotal)"
          icon="pi pi-calendar"
          tone="brand"
          :loading="isLoading"
        />
        <AppStatCard
          label="12-Month Avg"
          :value="fmt(avg12Month)"
          icon="pi pi-chart-line"
          tone="muted"
          :loading="isLoading"
        />
        <AppStatCard
          label="Highest Month"
          :value="fmt(highestMonth)"
          icon="pi pi-arrow-up"
          tone="positive"
          :loading="isLoading"
        />
      </div>

      <div class="bg-[#131316] border border-[#26262c] rounded-[14px] p-5 flex flex-col gap-3">
        <p class="text-[11px] font-semibold uppercase tracking-[0.08em] text-zinc-400">Monthly Spend (last 12 months)</p>
        <div class="h-64">
          <Chart v-if="!isLoading && sortedMonths.length" type="bar" :data="chartData" :options="chartOptions" class="h-full" />
          <div v-else-if="isLoading" class="h-full bg-[#1f1f24] rounded animate-pulse" />
          <div v-else class="h-full flex items-center justify-center text-zinc-500 text-sm">No data yet</div>
        </div>
      </div>

      <div class="bg-[#131316] border border-[#26262c] rounded-[14px] p-5 flex flex-col gap-4">
        <p class="text-[11px] font-semibold uppercase tracking-[0.08em] text-zinc-400">Spending by Category</p>
        <div v-if="isLoading" class="flex flex-col gap-3">
          <div v-for="i in 4" :key="i" class="h-8 bg-[#1f1f24] rounded animate-pulse" />
        </div>
        <div v-else-if="categoryAverages.length" class="flex flex-col gap-3">
          <div v-for="cat in categoryAverages" :key="cat.name" class="flex flex-col gap-1">
            <div class="flex items-center justify-between">
              <span class="text-sm text-white">{{ cat.name }}</span>
              <span class="text-sm text-zinc-400">{{ fmt(cat.total) }}</span>
            </div>
            <div class="h-2 rounded-full bg-[#26262c] overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :style="{ width: cat.pct + '%', backgroundColor: cat.color }"
              />
            </div>
          </div>
        </div>
        <div v-else class="text-zinc-500 text-sm">No category data yet</div>
      </div>
    </div>
  </Root>
</template>
