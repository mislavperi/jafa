<script setup lang="ts">
import { ref, computed } from 'vue'
import Panel from 'primevue/panel'
import Chart from 'primevue/chart'
import Select from 'primevue/select'
import Skeleton from 'primevue/skeleton'
import { useDailySpend } from '../composables/useExpenses'

const TIMEFRAMES = [
  { label: 'Last month', months: 1 },
  { label: 'Last 3 months', months: 3 },
  { label: 'Last 6 months', months: 6 },
  { label: 'Last year', months: 12 },
  { label: 'Last 2 years', months: 24 },
]

type Timeframe = { label: string; months: number }
const selectedTimeframe = ref<Timeframe>(TIMEFRAMES[1]!)
const months = computed(() => selectedTimeframe.value.months)
const { data: dailySpend, isLoading } = useDailySpend(months)

const chartData = computed(() => {
  const spendByDay = new Map<string, number>()
  for (const entry of dailySpend.value ?? []) {
    spendByDay.set(entry.day, entry.total)
  }

  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const start = new Date(today)
  start.setMonth(start.getMonth() - months.value)
  start.setDate(1)

  const labels: string[] = []
  const data: number[] = []
  let cumulative = 0

  const showYear = months.value > 12
  for (const d = new Date(start); d <= today; d.setDate(d.getDate() + 1)) {
    const key = d.toISOString().slice(0, 10)
    cumulative += spendByDay.get(key) ?? 0
    labels.push(d.toLocaleDateString('default', { month: 'short', day: 'numeric', year: showYear ? 'numeric' : undefined }))
    data.push(cumulative)
  }

  if (!data.length) return { labels: [], datasets: [] }

  return {
    labels,
    datasets: [
      {
        label: selectedTimeframe.value.label,
        data,
        borderColor: '#6366f1',
        backgroundColor: '#6366f133',
        tension: 0,
        fill: true,
        pointRadius: 0,
        pointHoverRadius: 4,
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
        label: (ctx: { raw: number }) => `$${ctx.raw.toFixed(2)}`,
      },
    },
  },
  scales: {
    x: {
      ticks: { maxTicksLimit: 12, maxRotation: 0 },
    },
    y: {
      title: { display: true, text: 'Cumulative Spend ($)' },
      beginAtZero: true,
    },
  },
}
</script>

<template>
  <Panel header="Cumulative Spend">
    <div class="flex items-center gap-2 mb-3">
      <Select
        v-model="selectedTimeframe"
        :options="TIMEFRAMES"
        option-label="label"
        size="small"
      />
    </div>
    <Skeleton v-if="isLoading" height="350px" />
    <div v-else-if="chartData.datasets.length" class="h-[350px]">
      <Chart type="line" :data="chartData" :options="chartOptions" class="w-full h-full" />
    </div>
    <div v-else class="text-center text-surface-500 py-8">
      No spend data available
    </div>
  </Panel>
</template>
