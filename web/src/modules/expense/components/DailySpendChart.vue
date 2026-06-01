<script setup lang="ts">
import { ref, computed, watchEffect } from 'vue'
import Panel from 'primevue/panel'
import Chart from 'primevue/chart'
import Select from 'primevue/select'
import Skeleton from 'primevue/skeleton'
import { useFirstExpenseDate, useDailySpendForMonth } from '../composables/useExpenses'
import { useThemeStore } from '@/stores/theme'

const theme = useThemeStore()

function hexWithAlpha(hex: string, alphaHex: string): string {
  if (/^#[0-9a-fA-F]{6}$/.test(hex)) return hex + alphaHex
  return hex
}

type MonthOption = { label: string; year: number; month: number }

const { data: firstExpenseDateData, isLoading: isLoadingFirst } = useFirstExpenseDate()

const availableMonths = computed<MonthOption[]>(() => {
  const firstDateStr = firstExpenseDateData.value?.firstDate
  if (!firstDateStr) return []

  const first = new Date(firstDateStr + 'T00:00:00')
  const today = new Date()
  const months: MonthOption[] = []

  const d = new Date(first.getFullYear(), first.getMonth(), 1)
  while (d.getFullYear() < today.getFullYear() || (d.getFullYear() === today.getFullYear() && d.getMonth() <= today.getMonth())) {
    months.push({
      label: d.toLocaleDateString('default', { month: 'long', year: 'numeric' }),
      year: d.getFullYear(),
      month: d.getMonth() + 1, // 1-indexed for the API
    })
    d.setMonth(d.getMonth() + 1)
  }

  return months.reverse()
})

const selectedMonth = ref<MonthOption | null>(null)

watchEffect(() => {
  if (selectedMonth.value === null && availableMonths.value.length > 0) {
    const today = new Date()
    const current = availableMonths.value.find(
      m => m.year === today.getFullYear() && m.month === today.getMonth() + 1,
    )
    selectedMonth.value = current ?? availableMonths.value[0] ?? null
  }
})

const selectedYear = computed(() => selectedMonth.value?.year ?? new Date().getFullYear())
const selectedMonthNum = computed(() => selectedMonth.value?.month ?? new Date().getMonth() + 1)

const { data: dailySpend, isLoading: isLoadingSpend } = useDailySpendForMonth(selectedYear, selectedMonthNum)

const isLoading = computed(() => isLoadingFirst.value || isLoadingSpend.value)

const chartData = computed(() => {
  const m = selectedMonth.value
  if (!m) return { labels: [], datasets: [] }

  const spendByDay = new Map<string, number>()
  for (const entry of dailySpend.value ?? []) {
    spendByDay.set(entry.day, entry.total)
  }

  const start = new Date(m.year, m.month - 1, 1)
  const end = new Date(m.year, m.month, 0)
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const labels: string[] = []
  const data: number[] = []
  const colors: string[] = []
  let cumulative = 0

  for (const d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    cumulative += spendByDay.get(key) ?? 0
    data.push(cumulative)
    const accent = theme.currentAccent().color
    if (d.getTime() === today.getTime()) colors.push(accent)
    else if (d < today) colors.push(hexWithAlpha(accent, '99'))
    else colors.push('#a1a1aa55')
    labels.push(d.toLocaleDateString('default', { month: 'short', day: 'numeric' }))
  }

  if (!data.length) return { labels: [], datasets: [] }

  return {
    labels,
    datasets: [
      {
        label: m.label,
        data,
        backgroundColor: colors,
        borderColor: colors,
        borderWidth: 0,
        borderRadius: 3,
        barPercentage: 0.85,
        categoryPercentage: 0.9,
      },
    ],
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    mode: 'index' as const,
    intersect: false,
  },
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        label: (ctx: { raw: number | null }) => ctx.raw !== null ? `€${ctx.raw.toFixed(2)}` : '',
      },
    },
  },
  scales: {
    x: {
      ticks: { maxTicksLimit: 12, maxRotation: 0, color: '#71717a' },
      grid: { display: false },
    },
    y: {
      title: { display: true, text: 'Cumulative Spend (€)', color: '#71717a' },
      beginAtZero: true,
      ticks: { color: '#71717a' },
      grid: { color: 'rgba(128,128,128,0.15)' },
    },
  },
}
</script>

<template>
  <Panel
    header="Cumulative Spend"
    class="h-full flex flex-col"
    :pt="{ toggleableContent: { class: 'flex-1 flex flex-col min-h-0' }, content: { class: 'flex-1 flex flex-col min-h-0' } }"
  >
    <div class="flex items-center gap-2 mb-3 shrink-0">
      <Select
        v-model="selectedMonth"
        :options="availableMonths"
        option-label="label"
        size="small"
        :loading="isLoadingFirst"
        placeholder="Select month"
      />
    </div>
    <Skeleton v-if="isLoading" class="flex-1" />
    <div v-else-if="chartData.datasets.length" class="flex-1 min-h-0">
      <Chart type="bar" :data="chartData" :options="chartOptions" class="w-full h-full" />
    </div>
    <div v-else class="text-center text-surface-500 py-8">
      No spend data available
    </div>
  </Panel>
</template>
