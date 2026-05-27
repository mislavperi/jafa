<script setup lang="ts">
import { ref, computed } from 'vue'
import Button from 'primevue/button'
import Chart from 'primevue/chart'
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import AppStatCard from '@/core/components/AppStatCard.vue'
import TotalSpendCard from '../components/TotalSpendCard.vue'
import DailySpendChart from '../components/DailySpendChart.vue'
import AddExpenseModal from '../components/AddExpenseModal.vue'
import { useExpenses, useExpensesByMonth, useMonthlyTotal } from '../composables/useExpenses'

const COLORS = ['#f5c518','#f97316','#3b82f6','#a855f7','#ec4899','#14b8a6','#ef4444','#71717a']

const showModal = ref(false)

const { data: allExpenses } = useExpenses()
const { data: monthlyTotal } = useMonthlyTotal()

// Current + last month for insights
const now = new Date()
const currentYear = ref(now.getFullYear())
const currentMonth = ref(now.getMonth() + 1)
const lastYear = ref(now.getMonth() === 0 ? now.getFullYear() - 1 : now.getFullYear())
const lastMonth = ref(now.getMonth() === 0 ? 12 : now.getMonth())

const { data: currentMonthExpenses } = useExpensesByMonth(currentYear, currentMonth)
const { data: lastMonthExpenses } = useExpensesByMonth(lastYear, lastMonth)

const currentTotal = computed(() => monthlyTotal.value?.total ?? 0)

const lastMonthTotal = computed(() =>
  (lastMonthExpenses.value ?? []).reduce((s, e) => s + e.amount, 0)
)

const monthDeltaPct = computed(() => {
  if (!lastMonthTotal.value) return 0
  return ((currentTotal.value - lastMonthTotal.value) / lastMonthTotal.value) * 100
})

// Biggest category by expense name grouping
const topCategory = computed(() => {
  const expenses = currentMonthExpenses.value ?? []
  if (!expenses.length) return null
  const map = new Map<string, number>()
  for (const e of expenses) map.set(e.name, (map.get(e.name) ?? 0) + e.amount)
  const [name, amount] = [...map.entries()].sort((a, b) => b[1] - a[1])[0] ?? []
  if (!name) return null
  const pct = currentTotal.value ? Math.round((amount / currentTotal.value) * 100) : 0
  return { name, amount, pct }
})

const dayOfMonth = now.getDate()
const daysInMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate()
const expectedSoFar = computed(() => {
  // budget unknown from API — use last month as proxy
  return lastMonthTotal.value ? (lastMonthTotal.value / daysInMonth) * dayOfMonth : 0
})
const onTrack = computed(() => currentTotal.value <= expectedSoFar.value)

const insights = computed(() => [
  {
    icon: onTrack.value ? 'pi pi-check-circle' : 'pi pi-exclamation-triangle',
    tone: onTrack.value ? 'green' : 'red',
    title: 'Budget pace',
    text: onTrack.value
      ? `On track — $${Math.abs(expectedSoFar.value - currentTotal.value).toFixed(0)} under pace for day ${dayOfMonth}`
      : `Ahead of pace by $${Math.abs(currentTotal.value - expectedSoFar.value).toFixed(0)} for day ${dayOfMonth}`,
  },
  {
    icon: monthDeltaPct.value < 0 ? 'pi pi-arrow-down' : 'pi pi-arrow-up',
    tone: monthDeltaPct.value < 0 ? 'green' : 'red',
    title: 'vs. last month',
    text: `Spent ${Math.abs(monthDeltaPct.value).toFixed(0)}% ${monthDeltaPct.value < 0 ? 'less' : 'more'} so far vs. same point last month`,
  },
  {
    icon: 'pi pi-star',
    tone: 'blue',
    title: 'Top expense',
    text: topCategory.value
      ? `"${topCategory.value.name}" is top at $${topCategory.value.amount.toFixed(0)} (${topCategory.value.pct}% of total)`
      : 'Add expenses to see trends',
  },
])

// Breakdown items for pie chart
const breakdownItems = computed(() => {
  const map = new Map<string, number>()
  for (const e of currentMonthExpenses.value ?? []) {
    map.set(e.name, (map.get(e.name) ?? 0) + e.amount)
  }
  return [...map.entries()]
    .map(([name, total]) => ({ name, total }))
    .sort((a, b) => b.total - a.total)
    .slice(0, 8)
})

// Recent 5 expenses for dashboard table
const recentExpenses = computed(() =>
  [...(allExpenses.value ?? [])]
    .sort((a, b) => new Date(b.created_at ?? '').getTime() - new Date(a.created_at ?? '').getTime())
    .slice(0, 5)
)

function formatDate(d?: string) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-GB', { day: 'numeric', month: 'short' })
}
</script>

<template>
  <Root>
    <div class="flex flex-col gap-5 h-full min-w-0 p-8 overflow-auto">
      <AppPageHeader title="Dashboard" subtitle="Your spending at a glance">
        <Button label="Add Expense" icon="pi pi-plus" size="small" @click="showModal = true" />
      </AppPageHeader>

      <!-- Stat cards -->
      <div class="grid grid-cols-3 gap-4">
        <TotalSpendCard />
        <AppStatCard label="Budget" icon="pi pi-chart-pie" tone="brand" subtitle="No budget set" />
        <AppStatCard label="Savings" icon="pi pi-piggy-bank" tone="positive" subtitle="Track savings" />
      </div>

      <!-- Insights strip -->
      <div class="grid grid-cols-3 gap-3">
        <div
          v-for="insight in insights"
          :key="insight.title"
          class="flex items-start gap-3 p-4 rounded-[14px] border border-[#26262c] bg-[#131316]"
        >
          <div
            class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0"
            :class="{
              'bg-emerald-500/10 text-emerald-400': insight.tone === 'green',
              'bg-red-500/10 text-red-400': insight.tone === 'red',
              'bg-blue-500/10 text-blue-400': insight.tone === 'blue',
            }"
          >
            <i :class="insight.icon" class="text-[13px]" />
          </div>
          <div class="min-w-0">
            <p class="text-[11px] font-medium text-zinc-400 mb-0.5">{{ insight.title }}</p>
            <p class="text-[13px] text-white leading-snug">{{ insight.text }}</p>
          </div>
        </div>
      </div>

      <!-- Main grid: recent expenses + breakdown -->
      <div class="grid grid-cols-[1.4fr_1fr] gap-5">
        <!-- Recent expenses -->
        <div class="bg-[#131316] border border-[#26262c] rounded-[14px] p-5">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-[11px] font-semibold uppercase tracking-[0.06em] text-zinc-400">Recent Expenses</h2>
            <RouterLink to="/expenses" class="text-[12px] text-zinc-400 hover:text-white flex items-center gap-1">
              View all <i class="pi pi-chevron-right text-[10px]" />
            </RouterLink>
          </div>
          <table class="w-full text-[13.5px]">
            <thead>
              <tr class="border-b border-[#26262c]">
                <th class="text-left text-[11px] font-semibold uppercase tracking-[0.06em] text-zinc-400 pb-2 px-1">Name</th>
                <th class="text-right text-[11px] font-semibold uppercase tracking-[0.06em] text-zinc-400 pb-2 px-1">Amount</th>
                <th class="text-right text-[11px] font-semibold uppercase tracking-[0.06em] text-zinc-400 pb-2 px-1">Date</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="exp in recentExpenses"
                :key="exp.id"
                class="border-b border-[#26262c] last:border-0 hover:bg-white/[0.015]"
              >
                <td class="py-3 px-1 font-medium text-white">{{ exp.name }}</td>
                <td class="py-3 px-1 text-right tabular-nums font-medium text-white">${{ exp.amount.toFixed(2) }}</td>
                <td class="py-3 px-1 text-right text-zinc-400 tabular-nums">{{ formatDate(exp.created_at) }}</td>
              </tr>
              <tr v-if="!recentExpenses.length">
                <td colspan="3" class="py-8 text-center text-zinc-500">No expenses yet</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Expense breakdown pie chart -->
        <div class="bg-[#131316] border border-[#26262c] rounded-[14px] p-5 flex flex-col gap-4">
          <h2 class="text-[11px] font-semibold uppercase tracking-[0.06em] text-zinc-400">Expense Breakdown</h2>
          <template v-if="currentMonthExpenses?.length">
            <div class="relative h-[200px]">
              <Chart
                type="doughnut"
                :data="{
                  labels: breakdownItems.map(i => i.name),
                  datasets: [{
                    data: breakdownItems.map(i => i.total),
                    backgroundColor: breakdownItems.map((_, idx) => COLORS[idx % COLORS.length]),
                    hoverBackgroundColor: breakdownItems.map((_, idx) => COLORS[idx % COLORS.length]),
                    borderWidth: 2,
                    borderColor: '#131316',
                  }]
                }"
                :options="{
                  responsive: true,
                  maintainAspectRatio: false,
                  cutout: '68%',
                  plugins: {
                    legend: { display: false },
                    tooltip: { callbacks: { label: (ctx) => ` $${Number(ctx.raw).toFixed(2)}` } }
                  }
                }"
                class="w-full h-full"
              />
            </div>
            <div class="flex flex-col gap-2">
              <div
                v-for="(item, i) in breakdownItems"
                :key="item.name"
                class="flex items-center gap-2"
              >
                <span class="w-2.5 h-2.5 rounded-sm shrink-0" :style="{ background: COLORS[i % COLORS.length] }" />
                <span class="flex-1 text-zinc-400 text-[12px] truncate">{{ item.name }}</span>
                <span class="text-white tabular-nums text-[12px] font-medium">${{ item.total.toFixed(0) }}</span>
                <span class="text-zinc-500 text-[11px] w-8 text-right tabular-nums">{{ currentTotal ? Math.round(item.total / currentTotal * 100) : 0 }}%</span>
              </div>
            </div>
          </template>
          <p v-else class="text-zinc-500 text-[13px] py-4 text-center">No data for this month</p>
        </div>
      </div>

      <!-- Bottom grid: daily chart + upcoming bills -->
      <div class="grid grid-cols-[1.6fr_1fr] gap-5">
        <DailySpendChart />

        <!-- Upcoming bills placeholder -->
        <div class="bg-[#131316] border border-[#26262c] rounded-[14px] p-5">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-[11px] font-semibold uppercase tracking-[0.06em] text-zinc-400">Upcoming Bills</h2>
            <button class="text-[12px] text-zinc-400 hover:text-white">Manage</button>
          </div>
          <div class="flex flex-col items-center justify-center py-8 text-zinc-500 gap-2">
            <i class="pi pi-calendar text-[24px] text-zinc-600" />
            <p class="text-[13px]">No upcoming bills</p>
            <p class="text-[12px] text-zinc-600">Recurring expense tracking coming soon</p>
          </div>
        </div>
      </div>
    </div>

    <AddExpenseModal v-model:visible="showModal" />
  </Root>
</template>
