<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Chart from 'primevue/chart'
import InputText from 'primevue/inputtext'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'

import { useExpenses } from '../composables/useExpenses'

const { data: expenses, isLoading } = useExpenses()

const searchQuery = ref('')
const selectedExpenseNames = ref(new Set<string>())

watch(searchQuery, () => {
  selectedExpenseNames.value = new Set()
})

const tableExpenses = computed(() => expenses.value ?? [])

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
  if (!expenses.value) return []
  let source = expenses.value
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
  layout: {
    padding: 5,
  },
  animation: {
    animateRotate: false,
    animateScale: false,
    duration: 300,
  },
  plugins: {
    legend: { display: false },
    tooltip: {
      filter: (tooltipItem: any) => tooltipItem.raw > 0,
    },
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
        :value="tableExpenses"
        :loading="isLoading"
        paginator
        :rows="5"
        scrollable
        scroll-height="flex"
      >
        <Column field="name" header="Name" sortable />
        <Column field="amount" header="Amount" sortable>
          <template #body="{ data }">
            ${{ data.amount.toFixed(2) }}
          </template>
        </Column>
      </DataTable>
    </Panel>

    <Panel header="Expense Breakdown">
      <div class="flex flex-col gap-3">
        <div v-if="aggregatedExpenses.length" class="flex flex-col md:flex-row gap-3 md:h-[350px] overflow-hidden">
          <div class="w-full md:flex-1 min-w-0 h-[200px] md:h-full">
            <Chart
              type="pie"
              :data="chartData"
              :options="chartOptions"
              class="w-full h-full"
            />
          </div>
          <div class="flex flex-col gap-1 min-h-0 md:w-[250px] md:min-w-[150px] max-h-[200px] md:max-h-none md:h-full">
            <IconField class="mb-1 shrink-0">
              <InputIcon class="pi pi-search" />
              <InputText
                v-model="searchQuery"
                placeholder="Search..."
                class="w-full"
                size="small"
              />
            </IconField>
            <div class="overflow-y-auto flex flex-col gap-0.5 min-h-0 flex-1">
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
        <div v-else class="text-center text-surface-500 py-8">
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
