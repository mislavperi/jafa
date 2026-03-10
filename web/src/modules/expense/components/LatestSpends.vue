<script setup lang="ts">
import { ref, computed, watch, watchEffect } from 'vue'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Chart from 'primevue/chart'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'

import { useExpenses, useFirstExpenseDate, useExpensesByMonth } from '../composables/useExpenses'
import ExpenseTagManager from './ExpenseTagManager.vue'

// Recent expenses (all, not month-filtered)
const { data: allExpenses, isLoading: isLoadingExpenses } = useExpenses()

// Row expansion state
const expandedRows = ref<Record<number, boolean>>({})

// Month selector for the pie chart
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
      month: d.getMonth() + 1,
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

const { data: monthExpenses, isLoading: isLoadingMonth } = useExpensesByMonth(selectedYear, selectedMonthNum)

// Search / filter for pie chart
const searchQuery = ref('')
const selectedExpenseNames = ref(new Set<string>())

watch(searchQuery, () => {
  selectedExpenseNames.value = new Set()
})

const COLORS = [
  '#6366f1', '#f59e0b', '#10b981', '#ef4444', '#8b5cf6',
  '#ec4899', '#14b8a6', '#f97316', '#06b6d4', '#84cc16',
]

interface AggregatedExpense {
  name: string
  total: number
  color: string
}

const aggregatedExpenses = computed<AggregatedExpense[]>(() => {
  if (!monthExpenses.value) return []
  let source = monthExpenses.value
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    source = source.filter((e) => e.name.toLowerCase().includes(query))
  }
  const map = new Map<string, number>()
  for (const e of source) {
    map.set(e.name, (map.get(e.name) ?? 0) + e.amount)
  }
  return [...map.entries()].map(([name, total], i) => ({
    name,
    total,
    color: COLORS[i % COLORS.length]!,
  }))
})

const chartData = computed(() => {
  const items = aggregatedExpenses.value
  if (!items.length) return { labels: [], datasets: [] }

  const hasSelection = selectedExpenseNames.value.size > 0

  return {
    labels: items.map((e) => e.name),
    datasets: [
      {
        data: items.map((e) =>
          hasSelection && !selectedExpenseNames.value.has(e.name) ? 0 : e.total,
        ),
        backgroundColor: items.map((e) =>
          hasSelection && !selectedExpenseNames.value.has(e.name)
            ? e.color + '20'
            : e.color,
        ),
        hoverBackgroundColor: items.map((e) => e.color),
      },
    ],
  }
})

const chartOptions = {
  layout: { padding: 5 },
  animation: { animateRotate: false, animateScale: false, duration: 300 },
  plugins: {
    legend: { display: false },
    tooltip: { filter: (tooltipItem: any) => tooltipItem.raw > 0 },
  },
  responsive: true,
  maintainAspectRatio: false,
  onClick: (_event: any, elements: any[]) => {
    if (!elements.length) return
    const index = elements[0].index
    const clickedName = chartData.value.labels[index]
    if (clickedName) toggleExpense(clickedName)
  },
}

function toggleExpense(name: string) {
  const next = new Set(selectedExpenseNames.value)
  if (next.has(name)) {
    next.delete(name)
  } else {
    next.add(name)
  }
  selectedExpenseNames.value = next
}

function clearSelection() {
  selectedExpenseNames.value = new Set()
}
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-3 sm:gap-4 p-2 sm:p-4">
    <Panel header="Recent Expenses">
      <DataTable
        :value="allExpenses ?? []"
        :loading="isLoadingExpenses"
        v-model:expandedRows="expandedRows"
        dataKey="id"
        paginator
        :rows="5"
        scrollable
        scroll-height="flex"
      >
        <Column expander style="width: 3rem" />
        <Column field="name" header="Name" sortable />
        <Column field="amount" header="Amount" sortable>
          <template #body="{ data }">
            ${{ data.amount.toFixed(2) }}
          </template>
        </Column>
        <template #expansion="{ data }">
          <div class="px-2">
            <p class="text-xs text-surface-400 font-medium mb-1 uppercase tracking-wide">Tags</p>
            <ExpenseTagManager :expense-id="data.id" />
          </div>
        </template>
      </DataTable>
    </Panel>

    <Panel header="Expense Breakdown">
      <template #icons>
        <Select
          v-model="selectedMonth"
          :options="availableMonths"
          option-label="label"
          size="small"
          :loading="isLoadingFirst"
          placeholder="Select month"
        />
      </template>
      <div class="flex flex-col gap-3">
        <div v-if="aggregatedExpenses.length" class="flex flex-col md:flex-row gap-3 md:h-[350px]">
          <div class="relative min-w-0 h-[200px] md:h-full md:flex-1">
            <Chart
              type="pie"
              :data="chartData"
              :options="chartOptions"
              class="absolute inset-0 w-full h-full"
            />
          </div>
          <div class="flex flex-col gap-1 min-h-0 md:w-2/5 md:max-w-[250px] md:h-full md:overflow-hidden">
            <IconField class="mb-1 shrink-0">
              <InputIcon class="pi pi-search" />
              <InputText
                v-model="searchQuery"
                placeholder="Search..."
                class="w-full"
                size="small"
              />
            </IconField>
            <div class="overflow-y-auto flex flex-col gap-0.5 min-h-0 flex-1 max-h-[150px] md:max-h-none">
              <div
                v-for="item in aggregatedExpenses"
                :key="item.name"
                class="flex items-center gap-1.5 px-2 py-1 rounded cursor-pointer transition-colors"
                :class="[
                  selectedExpenseNames.has(item.name)
                    ? 'bg-primary/10 ring-1 ring-primary'
                    : 'hover:bg-surface-100 dark:hover:bg-surface-800',
                ]"
                @click="toggleExpense(item.name)"
              >
                <span
                  class="inline-block w-2.5 h-2.5 rounded-full shrink-0"
                  :style="{ backgroundColor: item.color }"
                />
                <div class="flex flex-col min-w-0">
                  <span class="text-xs truncate">{{ item.name }}</span>
                  <span class="text-xs text-surface-500">${{ item.total.toFixed(2) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-else-if="!isLoadingMonth" class="text-center text-surface-500 py-8">
          No expenses found
        </div>

        <div v-if="selectedExpenseNames.size > 0" class="text-sm text-surface-500 text-center">
          {{ selectedExpenseNames.size }} selected
          <span class="ml-2 cursor-pointer underline" @click="clearSelection">Clear all</span>
        </div>
      </div>
    </Panel>
  </div>
</template>
