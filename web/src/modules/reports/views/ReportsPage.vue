<script setup lang="ts">
import { computed } from 'vue'
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import AppStatCard from '@/core/components/AppStatCard.vue'
import Chart from 'primevue/chart'
import { useMonthlySpend, useCategoryBreakdown } from '@/modules/reports/composables/useReports'
import { useThemeStore } from '@/stores/theme'
import { currencySymbol, formatCurrency } from '@/core/currency'

const { data: monthly, isLoading } = useMonthlySpend()
const { data: breakdown } = useCategoryBreakdown()
const theme = useThemeStore()

function fmt(val: number) {
  return formatCurrency(val, theme.currency)
}

// Backend returns months oldest-first as { month: 'YYYY-MM', total }.
const sortedMonths = computed(() => monthly.value ?? [])

const last12 = computed(() => sortedMonths.value.slice(-12))

const thisMonthKey = computed(() => {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
})

const thisMonthTotal = computed(
  () => sortedMonths.value.find((m) => m.month === thisMonthKey.value)?.total ?? 0,
)

const avg12Month = computed(() => {
  if (!last12.value.length) return 0
  const total = last12.value.reduce((s, m) => s + m.total, 0)
  return total / last12.value.length
})

const highestMonth = computed(() =>
  sortedMonths.value.reduce((max, m) => (m.total > max ? m.total : max), 0),
)

const chartData = computed(() => {
  const labels = last12.value.map((m) => {
    const [y, mo] = m.month.split('-')
    return new Date(Number(y), Number(mo) - 1).toLocaleString('default', { month: 'short', year: '2-digit' })
  })
  const data = last12.value.map((m) => m.total)
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

const chartOptions = computed(() => {
  const symbol = currencySymbol(theme.currency)
  return {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          label: (ctx: { parsed: { y: number } }) => symbol + ctx.parsed.y.toFixed(2),
        },
      },
    },
    scales: {
      x: {
        ticks: { color: '#71717a' },
        grid: { color: 'var(--jafa-border)' },
      },
      y: {
        ticks: {
          color: '#71717a',
          callback: (v: number) => symbol + v,
        },
        grid: { color: 'var(--jafa-border)' },
      },
    },
  }
})

// Bar widths are scaled relative to the largest category spend — purely a
// presentation concern; the totals themselves come from the backend.
const categoryAverages = computed(() => {
  const cats = (breakdown.value ?? []).filter((c) => c.spent > 0)
  const maxVal = Math.max(1, ...cats.map((c) => c.spent))
  return cats.map((c) => ({
    name: c.name,
    color: c.color,
    total: c.spent,
    pct: Math.round((c.spent / maxVal) * 100),
  }))
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

      <div class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 flex flex-col gap-3">
        <p class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.08em] text-[var(--jafa-text-muted)]">Monthly Spend (last 12 months)</p>
        <div class="h-64">
          <Chart v-if="!isLoading && sortedMonths.length" type="bar" :data="chartData" :options="chartOptions" class="h-full" />
          <div v-else-if="isLoading" class="h-full bg-[var(--jafa-surface-3)] rounded animate-pulse" />
          <div v-else class="h-full flex items-center justify-center text-[var(--jafa-text-muted)] text-sm">No data yet</div>
        </div>
      </div>

      <div class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 flex flex-col gap-4">
        <p class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.08em] text-[var(--jafa-text-muted)]">Spending by Category</p>
        <div v-if="isLoading" class="flex flex-col gap-3">
          <div v-for="i in 4" :key="i" class="h-8 bg-[var(--jafa-surface-3)] rounded animate-pulse" />
        </div>
        <div v-else-if="categoryAverages.length" class="flex flex-col gap-3">
          <div v-for="cat in categoryAverages" :key="cat.name" class="flex flex-col gap-1">
            <div class="flex items-center justify-between">
              <span class="text-sm text-[var(--jafa-text)]">{{ cat.name }}</span>
              <span class="text-sm text-[var(--jafa-text-muted)]">{{ fmt(cat.total) }}</span>
            </div>
            <div class="h-2 rounded-full bg-[var(--jafa-border)] overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :style="{ width: cat.pct + '%', backgroundColor: cat.color }"
              />
            </div>
          </div>
        </div>
        <div v-else class="text-[var(--jafa-text-muted)] text-sm">No category data yet</div>
      </div>
    </div>
  </Root>
</template>
