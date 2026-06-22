<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Button from 'primevue/button'
import Chart from 'primevue/chart'
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import AppStatCard from '@/core/components/AppStatCard.vue'
import TotalSpendCard from '../components/TotalSpendCard.vue'
import DailySpendChart from '../components/DailySpendChart.vue'
import AddExpenseModal from '../components/AddExpenseModal.vue'
import ReceiptScannerModal from '../components/ReceiptScannerModal.vue'
import ExpenseTagsCell from '../components/ExpenseTagsCell.vue'
import { useExpenses, useAllEntries, useExpensesByMonth, useMonthlyTotal, useMonthlyIncome, useDeleteExpense } from '../composables/useExpenses'
import { useDarkModeStore } from '@/stores/darkMode'
import { useThemeStore } from '@/stores/theme'
import { useOnboardingStore } from '@/stores/onboarding'
import { currencySymbol, formatCurrency } from '@/core/currency'
import type { Expense } from '../models/expense'

const COLORS = ['#f5c518','#f97316','#3b82f6','#a855f7','#ec4899','#14b8a6','#ef4444','#71717a']
const darkMode = useDarkModeStore()
const theme = useThemeStore()
const cs = computed(() => currencySymbol(theme.currency))
const money = (v: number) => formatCurrency(v, theme.currency)
const surfaceColor = computed(() => (darkMode.isDark ? '#131316' : '#ffffff'))

// First visit to the dashboard kicks off the onboarding tour (no-op once
// the user has finished or skipped it).
const onboarding = useOnboardingStore()
onMounted(() => onboarding.maybeStart())

const showModal = ref(false)
const showScanner = ref(false)
const editingExpense = ref<Expense | null>(null)
const { mutateAsync: deleteExpense } = useDeleteExpense()
function openEdit(e: Expense) {
  editingExpense.value = e
  showModal.value = true
}
async function removeExpense(e: Expense) {
  if (!confirm(`Delete "${e.name}"?`)) return
  await deleteExpense(e.id)
}

const { data: allExpenses } = useExpenses()
// Both kinds, for the recent transactions table (expenses-only `allExpenses`
// still drives upcoming bills so income isn't treated as a bill).
const { data: allEntries } = useAllEntries()
const isIncome = (e: Expense) => e.kind === 'income'
const { data: monthlyTotal } = useMonthlyTotal()
const { data: monthlyIncome } = useMonthlyIncome()

// Current + last month for insights
const now = new Date()
const currentYear = ref(now.getFullYear())
const currentMonth = ref(now.getMonth() + 1)
const lastYear = ref(now.getMonth() === 0 ? now.getFullYear() - 1 : now.getFullYear())
const lastMonth = ref(now.getMonth() === 0 ? 12 : now.getMonth())

const { data: currentMonthExpenses } = useExpensesByMonth(currentYear, currentMonth)
const { data: lastMonthExpenses } = useExpensesByMonth(lastYear, lastMonth)

const currentTotal = computed(() => monthlyTotal.value?.total ?? 0)
const incomeTotal = computed(() => monthlyIncome.value?.total ?? 0)

const lastMonthTotal = computed(() =>
  (lastMonthExpenses.value ?? []).reduce((s, e) => s + e.amount, 0)
)

const monthDeltaPct = computed(() => {
  if (!lastMonthTotal.value) return 0
  return ((currentTotal.value - lastMonthTotal.value) / lastMonthTotal.value) * 100
})

// Current-month spend grouped by expense name, largest first. Shared by the
// "top expense" insight and the breakdown pie so we only aggregate once.
const spendByName = computed(() => {
  const map = new Map<string, number>()
  for (const e of currentMonthExpenses.value ?? []) {
    map.set(e.name, (map.get(e.name) ?? 0) + e.amount)
  }
  return [...map.entries()]
    .map(([name, total]) => ({ name, total }))
    .sort((a, b) => b.total - a.total)
})

const topCategory = computed(() => {
  const top = spendByName.value[0]
  if (!top) return null
  const pct = currentTotal.value ? Math.round((top.total / currentTotal.value) * 100) : 0
  return { name: top.name, amount: top.total, pct }
})

// Monthly budget comes from user preferences (set on the Settings page); 0
// means no budget configured.
const monthlyBudget = computed(() => theme.monthlyBudget)
// Income tops up the configured budget: this month's income is added to the
// monthly budget to form the spendable pool.
const effectiveBudget = computed(() => monthlyBudget.value + incomeTotal.value)
const hasBudget = computed(() => effectiveBudget.value > 0)
const budgetPctUsed = computed(() =>
  effectiveBudget.value > 0 ? Math.round((currentTotal.value / effectiveBudget.value) * 100) : 0,
)
const budgetRemaining = computed(() => effectiveBudget.value - currentTotal.value)

const dayOfMonth = now.getDate()
const daysInMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate()
const expectedSoFar = computed(() => {
  // Pace against the configured budget; fall back to last month as a proxy.
  const reference = effectiveBudget.value > 0 ? effectiveBudget.value : lastMonthTotal.value
  return reference ? (reference / daysInMonth) * dayOfMonth : 0
})
const onTrack = computed(() => currentTotal.value <= expectedSoFar.value)

const insights = computed(() => [
  {
    icon: onTrack.value ? 'pi pi-check-circle' : 'pi pi-exclamation-triangle',
    tone: onTrack.value ? 'green' : 'red',
    title: 'Budget pace',
    text: onTrack.value
      ? `On track — ${cs.value}${Math.abs(expectedSoFar.value - currentTotal.value).toFixed(0)} under pace for day ${dayOfMonth}`
      : `Ahead of pace by ${cs.value}${Math.abs(currentTotal.value - expectedSoFar.value).toFixed(0)} for day ${dayOfMonth}`,
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
      ? `"${topCategory.value.name}" is top at ${cs.value}${topCategory.value.amount.toFixed(0)} (${topCategory.value.pct}% of total)`
      : 'Add expenses to see trends',
  },
])

// Breakdown items for pie chart (top 8 of the shared aggregation)
const breakdownItems = computed(() => spendByName.value.slice(0, 8))

// Upcoming bills: recurring expenses, sorted by next due date
interface UpcomingBill {
  id: number
  name: string
  cost: number
  interval: string
  nextDue: Date
  daysUntil: number
}

function nextOccurrence(interval: string, day: number, startDate: string): Date {
  const t = new Date()
  t.setHours(0, 0, 0, 0)
  const start = new Date(startDate + 'T00:00:00')
  if (interval === 'yearly') {
    const cand = new Date(t.getFullYear(), start.getMonth(), Math.min(day, 28))
    if (cand < t) cand.setFullYear(cand.getFullYear() + 1)
    return cand
  }
  // monthly
  const cand = new Date(t.getFullYear(), t.getMonth(), Math.min(day, 28))
  if (cand < t) cand.setMonth(cand.getMonth() + 1)
  return cand
}

const upcomingBills = computed<UpcomingBill[]>(() => {
  const t = new Date()
  t.setHours(0, 0, 0, 0)
  return (allExpenses.value ?? [])
    .filter((e) => e.recurringSchedule)
    .map((e) => {
      const s = e.recurringSchedule!
      const due = nextOccurrence(s.interval, s.dayOfMonth, s.startDate)
      const daysUntil = Math.round((due.getTime() - t.getTime()) / 86400000)
      return {
        id: e.id,
        name: e.name,
        cost: e.cost ?? e.amount,
        interval: s.interval,
        nextDue: due,
        daysUntil,
      }
    })
    .sort((a, b) => a.daysUntil - b.daysUntil)
    .slice(0, 6)
})

const upcomingTotal7d = computed(() =>
  upcomingBills.value.filter((b) => b.daysUntil <= 7).reduce((s, b) => s + b.cost, 0),
)

// Recent 5 transactions (expenses + income) for dashboard table
const recentExpenses = computed(() =>
  [...(allEntries.value ?? [])]
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
        <Button label="Scan Receipt" icon="pi pi-receipt" size="small" severity="secondary" data-tour="scan-receipt" @click="showScanner = true" />
        <Button label="Add Expense" icon="pi pi-plus" size="small" data-tour="add-expense" @click="showModal = true" />
      </AppPageHeader>

      <!-- Stat cards -->
      <div class="grid grid-cols-3 gap-4" data-tour="stat-cards">
        <TotalSpendCard />
        <AppStatCard
          label="Budget"
          icon="pi pi-chart-pie"
          tone="brand"
          :value="hasBudget ? money(effectiveBudget) : undefined"
          :subtitle="hasBudget
            ? (incomeTotal > 0
              ? `${budgetPctUsed}% used · +${money(incomeTotal)} income`
              : `${budgetPctUsed}% used this month`)
            : 'No budget set — add one in Settings'"
        />
        <AppStatCard
          label="Left to spend"
          icon="pi pi-piggy-bank"
          :tone="hasBudget && budgetRemaining < 0 ? 'brand' : 'positive'"
          :value="hasBudget ? money(Math.max(budgetRemaining, 0)) : undefined"
          :subtitle="hasBudget
            ? (budgetRemaining >= 0 ? 'Remaining this month' : `${money(-budgetRemaining)} over budget`)
            : 'Set a budget to track this'"
        />
      </div>

      <!-- Insights strip -->
      <div class="grid grid-cols-3 gap-3">
        <div
          v-for="insight in insights"
          :key="insight.title"
          class="flex items-start gap-3 p-4 rounded-[14px] border border-[var(--jafa-border)] bg-[var(--jafa-surface)]"
        >
          <div
            class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0"
            :class="{
              'bg-emerald-500/10 text-emerald-400': insight.tone === 'green',
              'bg-red-500/10 text-red-400': insight.tone === 'red',
              'bg-blue-500/10 text-blue-400': insight.tone === 'blue',
            }"
          >
            <i :class="insight.icon" class="text-[calc(13px*var(--jafa-text-scale,1))]" />
          </div>
          <div class="min-w-0">
            <p class="text-[calc(11px*var(--jafa-text-scale,1))] font-medium text-[var(--jafa-text-muted)] mb-0.5">{{ insight.title }}</p>
            <p class="text-[calc(13px*var(--jafa-text-scale,1))] text-[var(--jafa-text)] leading-snug">{{ insight.text }}</p>
          </div>
        </div>
      </div>

      <!-- Main grid: recent expenses + breakdown -->
      <div class="grid grid-cols-[1.4fr_1fr] gap-5">
        <!-- Recent expenses -->
        <div class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5" data-tour="recent-expenses">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em] text-[var(--jafa-text-muted)]">Recent Transactions</h2>
            <RouterLink to="/expenses" class="text-[calc(12px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] hover:text-[var(--jafa-text)] flex items-center gap-1">
              View all <i class="pi pi-chevron-right text-[calc(10px*var(--jafa-text-scale,1))]" />
            </RouterLink>
          </div>
          <table class="w-full text-[calc(13.5px*var(--jafa-text-scale,1))]">
            <thead>
              <tr class="border-b border-[var(--jafa-border)]">
                <th class="text-left text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em] text-[var(--jafa-text-muted)] pb-2.5 px-2">Expense Name</th>
                <th class="text-left text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em] text-[var(--jafa-text-muted)] pb-2.5 px-2">Tags</th>
                <th class="text-right text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em] text-[var(--jafa-text-muted)] pb-2.5 px-2">Amount</th>
                <th class="text-left text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em] text-[var(--jafa-text-muted)] pb-2.5 px-2">Date</th>
                <th class="text-right text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em] text-[var(--jafa-text-muted)] pb-2.5 px-2">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="exp in recentExpenses"
                :key="exp.id"
                class="border-b border-[var(--jafa-border)] last:border-0 hover:bg-[var(--jafa-hover)] group"
              >
                <td class="py-3.5 px-2 font-medium text-[var(--jafa-text)]">
                  <span class="inline-flex items-center gap-2">
                    {{ exp.name }}
                    <span
                      v-if="isIncome(exp)"
                      class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-emerald-500/15 text-emerald-500 text-[calc(10px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em]"
                    >
                      <i class="pi pi-arrow-down-left text-[calc(9px*var(--jafa-text-scale,1))]" />
                      income
                    </span>
                    <span
                      v-if="exp.recurringSchedule"
                      class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-[var(--jafa-accent)]/15 text-[var(--jafa-accent)] text-[calc(10px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em]"
                    >
                      <i class="pi pi-replay text-[calc(9px*var(--jafa-text-scale,1))]" />
                      {{ exp.recurringSchedule.interval }}
                    </span>
                  </span>
                </td>
                <td class="py-3.5 px-2"><ExpenseTagsCell :expense-id="exp.id" /></td>
                <td
                  class="py-3.5 px-2 text-right tabular-nums font-semibold"
                  :class="isIncome(exp) ? 'text-emerald-500' : 'text-red-500'"
                >{{ money(exp.cost ?? exp.amount) }}</td>
                <td class="py-3.5 px-2 text-[var(--jafa-text-muted)] tabular-nums">{{ formatDate(exp.created_at) }}</td>
                <td class="py-3.5 px-2 text-right">
                  <div class="flex items-center justify-end gap-1 opacity-60 group-hover:opacity-100 transition">
                    <button
                      class="w-7 h-7 inline-flex items-center justify-center rounded-md text-[var(--jafa-text-muted)] hover:bg-[var(--jafa-border)] hover:text-[var(--jafa-text)] transition"
                      @click="openEdit(exp)"
                    >
                      <i class="pi pi-pencil text-[calc(12px*var(--jafa-text-scale,1))]" />
                    </button>
                    <button
                      class="w-7 h-7 inline-flex items-center justify-center rounded-md text-[var(--jafa-text-muted)] hover:bg-[var(--jafa-border)] hover:text-red-400 transition"
                      @click="removeExpense(exp)"
                    >
                      <i class="pi pi-trash text-[calc(12px*var(--jafa-text-scale,1))]" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="!recentExpenses.length">
                <td colspan="5" class="py-8 text-center text-[var(--jafa-text-muted)]">No transactions yet</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Expense breakdown pie chart -->
        <div class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 flex flex-col gap-4" data-tour="breakdown">
          <h2 class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em] text-[var(--jafa-text-muted)]">Expense Breakdown</h2>
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
                    borderColor: surfaceColor,
                  }]
                }"
                :options="{
                  responsive: true,
                  maintainAspectRatio: false,
                  cutout: '68%',
                  plugins: {
                    legend: { display: false },
                    tooltip: { callbacks: { label: (ctx: { raw: unknown }) => ` ${money(Number(ctx.raw))}` } }
                  }
                }"
                class="w-full h-full"
              />
              <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
                <span class="text-[calc(10px*var(--jafa-text-scale,1))] uppercase tracking-[0.14em] text-[var(--jafa-text-muted)]">This Month</span>
                <span class="text-[calc(22px*var(--jafa-text-scale,1))] font-semibold tabular-nums text-[var(--jafa-text)] leading-tight mt-0.5">{{ money(currentTotal) }}</span>
              </div>
            </div>
            <div class="flex flex-col gap-2">
              <div
                v-for="(item, i) in breakdownItems"
                :key="item.name"
                class="flex items-center gap-2"
              >
                <span class="w-2.5 h-2.5 rounded-sm shrink-0" :style="{ background: COLORS[i % COLORS.length] }" />
                <span class="flex-1 text-[var(--jafa-text-muted)] text-[calc(12px*var(--jafa-text-scale,1))] truncate">{{ item.name }}</span>
                <span class="text-[var(--jafa-text)] tabular-nums text-[calc(12px*var(--jafa-text-scale,1))] font-medium">{{ cs }}{{ item.total.toFixed(0) }}</span>
                <span class="text-[var(--jafa-text-muted)] text-[calc(11px*var(--jafa-text-scale,1))] w-8 text-right tabular-nums">{{ currentTotal ? Math.round(item.total / currentTotal * 100) : 0 }}%</span>
              </div>
            </div>
          </template>
          <p v-else class="text-[var(--jafa-text-muted)] text-[calc(13px*var(--jafa-text-scale,1))] py-4 text-center">No data for this month</p>
        </div>
      </div>

      <!-- Bottom grid: daily chart + upcoming bills -->
      <div class="grid grid-cols-[1.6fr_1fr] gap-5">
        <DailySpendChart />

        <!-- Upcoming bills -->
        <div class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 flex flex-col" data-tour="upcoming-bills">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em] text-[var(--jafa-text-muted)]">Upcoming Bills</h2>
            <RouterLink to="/expenses" class="text-[calc(12px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] hover:text-[var(--jafa-text)]">Manage</RouterLink>
          </div>
          <template v-if="upcomingBills.length">
            <div class="flex flex-col">
              <div
                v-for="b in upcomingBills"
                :key="b.id"
                class="flex items-center gap-3 py-2.5 border-b border-[var(--jafa-border)] last:border-0"
              >
                <div
                  class="w-9 h-9 rounded-lg flex flex-col items-center justify-center shrink-0 bg-[var(--jafa-surface-3)]"
                >
                  <span class="text-[calc(9px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-wider text-[var(--jafa-text-muted)] leading-none">
                    {{ b.nextDue.toLocaleDateString('default', { month: 'short' }) }}
                  </span>
                  <span class="text-[calc(13px*var(--jafa-text-scale,1))] font-semibold text-[var(--jafa-text)] tabular-nums leading-tight">
                    {{ b.nextDue.getDate() }}
                  </span>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-1.5">
                    <span class="text-[calc(13px*var(--jafa-text-scale,1))] font-medium text-[var(--jafa-text)] truncate">{{ b.name }}</span>
                    <span
                      v-if="b.daysUntil <= 3"
                      class="px-1.5 py-0.5 rounded bg-[var(--jafa-accent)]/15 text-[var(--jafa-accent)] text-[calc(9px*var(--jafa-text-scale,1))] font-bold uppercase tracking-wider"
                    >
                      Soon
                    </span>
                  </div>
                  <div class="text-[calc(11px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] mt-0.5">
                    <span class="uppercase tracking-wider">{{ b.interval }}</span>
                    <span class="mx-1.5">·</span>
                    <span>{{ b.daysUntil === 0 ? 'Today' : b.daysUntil === 1 ? 'Tomorrow' : `in ${b.daysUntil}d` }}</span>
                  </div>
                </div>
                <div class="text-[calc(13px*var(--jafa-text-scale,1))] font-semibold text-[var(--jafa-text)] tabular-nums">{{ money(b.cost) }}</div>
              </div>
            </div>
            <div class="mt-auto pt-3 flex items-center justify-between border-t border-[var(--jafa-border)]">
              <span class="text-[calc(11px*var(--jafa-text-scale,1))] uppercase tracking-[0.08em] text-[var(--jafa-text-muted)]">Next 7 days</span>
              <span class="text-[calc(13px*var(--jafa-text-scale,1))] font-semibold tabular-nums text-[var(--jafa-text)]">{{ money(upcomingTotal7d) }}</span>
            </div>
          </template>
          <div v-else class="flex flex-col items-center justify-center py-8 text-[var(--jafa-text-muted)] gap-2">
            <i class="pi pi-calendar text-[calc(24px*var(--jafa-text-scale,1))] text-[var(--jafa-text-dim)]" />
            <p class="text-[calc(13px*var(--jafa-text-scale,1))]">No upcoming bills</p>
            <p class="text-[calc(12px*var(--jafa-text-scale,1))] text-[var(--jafa-text-dim)]">Mark expenses as recurring to track them here</p>
          </div>
        </div>
      </div>
    </div>

    <AddExpenseModal v-model:visible="showModal" :expense="editingExpense" @update:visible="(v) => { if (!v) editingExpense = null }" />
    <ReceiptScannerModal v-model:visible="showScanner" />
  </Root>
</template>
