<script setup lang="ts">
import { ref, computed, watch, watchEffect, onMounted, onUnmounted } from 'vue'
import { useQueries } from '@tanstack/vue-query'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Chart from 'primevue/chart'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import AddExpenseModal from './AddExpenseModal.vue'
import type { Expense, Tag } from '../models/expense'
import { CHART_COLORS } from '../constants'

import { useExpenses, useFirstExpenseDate, useExpensesByMonth, useDeleteExpense } from '../composables/useExpenses'
import { useDarkModeStore } from '@/stores/darkMode'
import { useThemeStore } from '@/stores/theme'
import { currencySymbol, formatCurrency } from '@/core/currency'
const darkMode = useDarkModeStore()
const theme = useThemeStore()
const cs = computed(() => currencySymbol(theme.currency))
const money = (v: number) => formatCurrency(v, theme.currency)
const surfaceColor = computed(() => (darkMode.isDark ? '#131316' : '#ffffff'))
import { useAllTags } from '../composables/useTags'
import { getTagsForExpense } from '../api/tag'
import ExpenseTagsCell from './ExpenseTagsCell.vue'

const { data: allExpenses, isLoading: isLoadingExpenses } = useExpenses()
const { data: rawTags } = useAllTags()
const allTags = computed(() => {
  if (!rawTags.value) return []
  const seen = new Set<number>()
  return rawTags.value.filter((t) => { if (seen.has(t.id)) return false; seen.add(t.id); return true })
})

// Search state
const searchWrapperRef = ref<HTMLElement | null>(null)
const searchInputRef = ref<HTMLInputElement | null>(null)
const showDropdown = ref(false)
const tableSearch = ref('')
const selectedFilterTagIds = ref<Set<number>>(new Set())

// Tag data — lazy: enabled only once user picks a tag, no upfront N+1 cost
const expenseTagQueries = useQueries({
  queries: computed(() =>
    (allExpenses.value ?? []).map((e) => ({
      queryKey: ['expense-tags', e.id],
      queryFn: () => getTagsForExpense(e.id),
      enabled: selectedFilterTagIds.value.size > 0,
      staleTime: 60_000,
    })),
  ),
})

const expenseTagMap = computed(() => {
  const map = new Map<number, Tag[]>()
  if (selectedFilterTagIds.value.size === 0) return map
  ;(allExpenses.value ?? []).forEach((e, i) => {
    const d = expenseTagQueries.value[i]?.data
    if (d) map.set(e.id, d)
  })
  return map
})

const isLoadingTagFilter = computed(() =>
  selectedFilterTagIds.value.size > 0 &&
  expenseTagQueries.value.some((q) => q.isLoading),
)

// Name suggestions derived from allExpenses — zero extra requests
const nameSuggestions = computed(() => {
  if (!tableSearch.value.trim()) return []
  const q = tableSearch.value.toLowerCase()
  const names = new Set<string>()
  for (const e of allExpenses.value ?? []) {
    if (e.name.toLowerCase().includes(q)) {
      names.add(e.name)
      if (names.size >= 6) break
    }
  }
  return [...names]
})

function toggleFilterTag(tagId: number) {
  const next = new Set(selectedFilterTagIds.value)
  if (next.has(tagId)) next.delete(tagId)
  else next.add(tagId)
  selectedFilterTagIds.value = next
}

function tagName(tagId: number) {
  return allTags.value?.find((t) => t.id === tagId)?.name ?? `Tag ${tagId}`
}

function selectSuggestion(name: string) {
  tableSearch.value = name
  showDropdown.value = false
}

function clearSearch() {
  tableSearch.value = ''
  selectedFilterTagIds.value = new Set()
}

function handleClickOutside(e: MouseEvent) {
  if (searchWrapperRef.value && !searchWrapperRef.value.contains(e.target as Node)) {
    showDropdown.value = false
  }
}

onMounted(() => document.addEventListener('mousedown', handleClickOutside))
onUnmounted(() => document.removeEventListener('mousedown', handleClickOutside))

const filteredExpenses = computed<Expense[]>(() => {
  const expenses = allExpenses.value ?? []
  let result = expenses

  if (tableSearch.value.trim()) {
    const q = tableSearch.value.toLowerCase()
    result = result.filter((e) => e.name.toLowerCase().includes(q))
  }

  if (selectedFilterTagIds.value.size > 0 && !isLoadingTagFilter.value) {
    result = result.filter((e) => {
      const tags = expenseTagMap.value.get(e.id)
      if (!tags) return false
      return [...selectedFilterTagIds.value].some((id) => tags.some((t) => t.id === id))
    })
  }

  return result
})

const selectedRows = ref<Expense[]>([])

let syncing = false
watch(selectedRows, (rows) => {
  if (syncing) return
  syncing = true
  selectedExpenseNames.value = new Set(rows.map((r) => r.name))
  syncing = false
})

function formatDate(dateStr?: string) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' })
}

const editingExpense = ref<Expense | null>(null)
const showEditModal = ref(false)
const { mutateAsync: deleteExpense } = useDeleteExpense()

function openEdit(expense: Expense) {
  editingExpense.value = expense
  showEditModal.value = true
}

async function removeExpense(expense: Expense) {
  if (!confirm(`Delete "${expense.name}"?`)) return
  await deleteExpense(expense.id)
}

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
  syncRowsFromNames(new Set())
})

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
    color: CHART_COLORS[i % CHART_COLORS.length]!,
  }))
})

const colorByName = computed(() => {
  const m = new Map<string, string>()
  aggregatedExpenses.value.forEach((e) => m.set(e.name, e.color))
  return m
})

const chartData = computed(() => {
  if (!selectedExpenseNames.value.size) return { labels: [], datasets: [] }

  const totals = new Map<string, number>()
  for (const row of selectedRows.value) {
    if (!selectedExpenseNames.value.has(row.name)) continue
    totals.set(row.name, (totals.get(row.name) ?? 0) + row.amount)
  }

  for (const name of selectedExpenseNames.value) {
    if (totals.has(name)) continue
    const agg = aggregatedExpenses.value.find((a) => a.name === name)
    if (agg) totals.set(name, agg.total)
  }

  if (!totals.size) return { labels: [], datasets: [] }

  const labels = [...totals.keys()]
  const data = labels.map((n) => totals.get(n) ?? 0)
  const colors = labels.map((n, i) => colorByName.value.get(n) ?? CHART_COLORS[i % CHART_COLORS.length]!)

  return {
    labels,
    datasets: [
      {
        data,
        backgroundColor: colors,
        hoverBackgroundColor: colors,
        borderWidth: 2,
        borderColor: surfaceColor.value,
      },
    ],
  }
})

const chartTotal = computed(() =>
  chartData.value.datasets[0]?.data.reduce((s: number, n: number) => s + n, 0) ?? 0,
)

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  cutout: '68%',
  layout: { padding: 5 },
  animation: { animateRotate: false, animateScale: false, duration: 300 },
  plugins: {
    legend: { display: false },
    tooltip: { callbacks: { label: (ctx: any) => ` ${money(Number(ctx.raw))}` } },
  },
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
  syncRowsFromNames(next)
}

function clearSelection() {
  selectedExpenseNames.value = new Set()
  syncRowsFromNames(new Set())
}

function syncRowsFromNames(names: Set<string>) {
  if (syncing) return
  syncing = true
  selectedRows.value = (filteredExpenses.value ?? []).filter((e) => names.has(e.name))
  syncing = false
}
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 h-full">
    <Panel
      header="RECENT EXPENSES"
      class="h-full flex flex-col bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] overflow-hidden"
      :pt="{
        root: { class: '!bg-[var(--jafa-surface)] !border-[var(--jafa-border)]' },
        header: { class: '!bg-[var(--jafa-surface)] !border-b !border-[var(--jafa-border)] !px-5 !py-4 [&_.p-panel-title]:!text-[calc(11px*var(--jafa-text-scale,1))] [&_.p-panel-title]:!font-semibold [&_.p-panel-title]:!uppercase [&_.p-panel-title]:!tracking-[0.06em] [&_.p-panel-title]:!text-[var(--jafa-text-muted)]' },
        toggleableContent: { class: 'flex-1 flex flex-col min-h-0' },
        content: { class: 'flex-1 flex flex-col min-h-0 !p-0 !bg-[var(--jafa-surface)] overflow-hidden' }
      }"
    >
      <!-- Smart search bar -->
      <div ref="searchWrapperRef" class="relative px-3 py-2 border-b border-surface shrink-0">
        <div
          class="flex flex-wrap items-center gap-1 px-2 py-1.5 rounded-lg border border-surface bg-white dark:bg-surface-800 focus-within:border-primary/60 transition-colors cursor-text"
          @click="searchInputRef?.focus(); showDropdown = true"
        >
          <i class="pi pi-search text-surface-400 text-xs mr-0.5 shrink-0" />

          <!-- Active tag chips -->
          <span
            v-for="tagId in selectedFilterTagIds"
            :key="tagId"
            class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded bg-primary/15 text-primary text-xs font-medium shrink-0"
          >
            {{ tagName(tagId) }}
            <button class="hover:text-primary/60 leading-none text-sm" @click.stop="toggleFilterTag(tagId)">×</button>
          </span>

          <input
            ref="searchInputRef"
            v-model="tableSearch"
            placeholder="Search or filter by tag..."
            class="flex-1 min-w-[6rem] bg-transparent outline-none text-sm text-surface-700 dark:text-surface-200 placeholder:text-surface-400"
            @focus="showDropdown = true"
            @keydown.escape="showDropdown = false; tableSearch = ''"
          />

          <i
            v-if="isLoadingTagFilter"
            class="pi pi-spin pi-spinner text-surface-400 text-xs ml-1 shrink-0"
          />
          <button
            v-else-if="tableSearch || selectedFilterTagIds.size > 0"
            class="text-surface-400 hover:text-surface-600 ml-1 shrink-0"
            @click.stop="clearSearch"
          >
            <i class="pi pi-times text-xs" />
          </button>
        </div>

        <!-- Dropdown -->
        <div
          v-if="showDropdown && (allTags?.length || nameSuggestions.length)"
          class="absolute left-3 right-3 top-full z-50 mt-1 bg-white dark:bg-surface-800 border border-surface rounded-xl shadow-lg overflow-hidden"
        >
          <!-- Tags section -->
          <div v-if="allTags?.length" class="px-3 pt-3 pb-2">
            <p class="text-[0.6rem] font-semibold uppercase tracking-widest text-surface-400 mb-2">Filter by tag</p>
            <div class="flex flex-wrap gap-1.5">
              <button
                v-for="tag in allTags"
                :key="tag.id"
                class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium transition-all border"
                :class="selectedFilterTagIds.has(tag.id)
                  ? 'bg-primary text-[var(--jafa-text)] border-primary'
                  : 'bg-transparent border-surface text-surface-500 hover:border-primary/50 hover:text-primary'"
                @mousedown.prevent="toggleFilterTag(tag.id)"
              >
                <i v-if="selectedFilterTagIds.has(tag.id)" class="pi pi-check text-[0.55rem]" />
                {{ tag.name }}
              </button>
            </div>
          </div>

          <!-- Name suggestions -->
          <div v-if="nameSuggestions.length" class="border-t border-surface py-1">
            <p class="text-[0.6rem] font-semibold uppercase tracking-widest text-surface-400 px-3 py-1.5">Suggestions</p>
            <div
              v-for="name in nameSuggestions"
              :key="name"
              class="flex items-center gap-2 px-3 py-2 cursor-pointer hover:bg-surface-50 dark:hover:bg-surface-700 text-sm"
              @mousedown.prevent="selectSuggestion(name)"
            >
              <i class="pi pi-arrow-up-right text-surface-400 text-xs shrink-0" />
              {{ name }}
            </div>
          </div>
        </div>
      </div>

      <!-- Bulk action bar -->
        <div
          v-if="selectedRows.length"
          class="flex items-center gap-3 px-4 py-2 bg-surface-50 dark:bg-surface-800 border-b border-surface text-sm shrink-0"
        >
          <span class="font-medium text-surface-600 dark:text-surface-300">{{ selectedRows.length }} selected</span>
          <span class="text-surface-300 dark:text-surface-600">|</span>
          <button class="flex items-center gap-1.5 text-red-500 hover:text-red-600 transition-colors font-medium">
            <i class="pi pi-trash text-xs" /> Delete
          </button>
          <button class="ml-auto text-surface-400 hover:text-surface-600 transition-colors" @click="selectedRows = []">
            <i class="pi pi-times text-xs" />
          </button>
        </div>

      <DataTable
        :value="filteredExpenses"
        :loading="isLoadingExpenses"
        v-model:selection="selectedRows"
        dataKey="id"
        paginator
        :rows="7"
        scrollable
        scroll-height="flex"
        class="jafa-table"
        :pt="{
          table: { class: '!bg-[var(--jafa-surface)]' },
          thead: { class: '!bg-[var(--jafa-surface)]' },
          headerRow: { class: '!border-b !border-[var(--jafa-border)]' },
          headerCell: { class: '!bg-[var(--jafa-surface)] !border-0 !text-[calc(11px*var(--jafa-text-scale,1))] !font-semibold !uppercase !tracking-[0.06em] !text-[var(--jafa-text-muted)] !px-3 !py-2.5' },
          tbody: { class: '!bg-[var(--jafa-surface)]' },
          bodyRow: { class: '!bg-[var(--jafa-surface)] hover:!bg-[var(--jafa-hover)] !border-0' },
          bodyCell: { class: '!bg-transparent !border-b !border-[var(--jafa-border)] !text-[calc(13.5px*var(--jafa-text-scale,1))] !px-3 !py-3.5 !text-[var(--jafa-text)]' },
        }"
      >
        <Column selectionMode="multiple" style="width: 2.75rem; padding-left: 1rem; padding-right: 0" />

        <Column field="name" header="Expense Name" sortable>
          <template #body="{ data }">
            <span class="inline-flex items-center gap-2 font-medium text-[calc(13.5px*var(--jafa-text-scale,1))] text-[var(--jafa-text)]">
              {{ data.name }}
              <span
                v-if="data.recurringSchedule"
                class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-[var(--jafa-accent)]/15 text-[var(--jafa-accent)] text-[calc(10px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.06em]"
              >
                <i class="pi pi-replay text-[calc(9px*var(--jafa-text-scale,1))]" />
                {{ data.recurringSchedule.interval }}
              </span>
            </span>
          </template>
        </Column>

        <Column header="Tags">
          <template #body="{ data }">
            <ExpenseTagsCell :expense-id="data.id" />
          </template>
        </Column>

        <Column field="amount" header="Amount" sortable style="width: 8rem">
          <template #body="{ data }">
            <span class="font-semibold tabular-nums text-[var(--jafa-text)]">−{{ money(data.cost ?? data.amount) }}</span>
          </template>
        </Column>

        <Column field="created_at" header="Date" sortable style="width: 9rem">
          <template #body="{ data }">
            <span class="text-[var(--jafa-text-muted)] text-[calc(13px*var(--jafa-text-scale,1))] tabular-nums">{{ formatDate(data.created_at) }}</span>
          </template>
        </Column>

        <Column header="Actions" style="width: 6rem">
          <template #header>
            <span class="w-full text-right block">Actions</span>
          </template>
          <template #body="{ data }">
            <div class="flex items-center justify-end gap-1 opacity-60 hover:opacity-100 transition">
              <button
                class="w-7 h-7 inline-flex items-center justify-center rounded-md text-[var(--jafa-text-muted)] hover:bg-[var(--jafa-border)] hover:text-[var(--jafa-text)] transition"
                @click="openEdit(data)"
              >
                <i class="pi pi-pencil text-[calc(12px*var(--jafa-text-scale,1))]" />
              </button>
              <button
                class="w-7 h-7 inline-flex items-center justify-center rounded-md text-[var(--jafa-text-muted)] hover:bg-[var(--jafa-border)] hover:text-red-400 transition"
                @click="removeExpense(data)"
              >
                <i class="pi pi-trash text-[calc(12px*var(--jafa-text-scale,1))]" />
              </button>
            </div>
          </template>
        </Column>

        <template #empty>
          <div class="flex flex-col items-center gap-2 py-12 text-[var(--jafa-text-muted)]">
            <div class="w-10 h-10 rounded-full bg-[var(--jafa-surface-3)] flex items-center justify-center">
              <i class="pi pi-wallet text-[calc(16px*var(--jafa-text-scale,1))]" />
            </div>
            <div class="text-[calc(13px*var(--jafa-text-scale,1))]">No expenses to show</div>
          </div>
        </template>
      </DataTable>

      <AddExpenseModal
        v-model:visible="showEditModal"
        :expense="editingExpense"
        @update:visible="(v: boolean) => { if (!v) editingExpense = null }"
      />
    </Panel>

    <Panel
      header="Expense Breakdown"
      class="h-full flex flex-col"
      :pt="{ toggleableContent: { class: 'flex-1 flex flex-col min-h-0' }, content: { class: 'flex-1 flex flex-col min-h-0' } }"
    >
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
      <div class="flex flex-col gap-4 flex-1 min-h-0">
        <div v-if="aggregatedExpenses.length" class="flex flex-col gap-4 flex-1 min-h-0">
          <div class="relative h-[220px] shrink-0 flex items-center justify-center">
            <template v-if="selectedExpenseNames.size > 0">
              <Chart
                type="doughnut"
                :data="chartData"
                :options="chartOptions"
                class="w-full h-full"
              />
              <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
                <span class="text-[calc(10px*var(--jafa-text-scale,1))] uppercase tracking-[0.14em] text-[var(--jafa-text-muted)]">Selected</span>
                <span class="text-[calc(22px*var(--jafa-text-scale,1))] font-semibold tabular-nums text-[var(--jafa-text)] leading-tight mt-0.5">{{ money(chartTotal) }}</span>
                <span class="text-[calc(10px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] mt-1">{{ selectedExpenseNames.size }} item{{ selectedExpenseNames.size === 1 ? '' : 's' }}</span>
              </div>
            </template>
            <div v-else class="text-center text-surface-400 text-sm px-4">
              <i class="pi pi-chart-pie text-2xl mb-2 block opacity-30" />
              Select an expense to see breakdown
            </div>
          </div>

          <IconField class="shrink-0">
            <InputIcon class="pi pi-search" />
            <InputText
              v-model="searchQuery"
              placeholder="Search..."
              class="w-full"
              size="small"
            />
          </IconField>

          <div class="flex flex-col gap-2 overflow-y-auto min-h-0 flex-1">
            <div
              v-for="item in aggregatedExpenses"
              :key="item.name"
              class="flex items-center gap-2 px-2 py-1 rounded cursor-pointer transition-colors"
              :class="[
                selectedExpenseNames.has(item.name)
                  ? 'bg-primary/10 ring-1 ring-primary'
                  : 'hover:bg-white/[0.025]',
              ]"
              @click="toggleExpense(item.name)"
            >
              <span
                class="w-2.5 h-2.5 rounded-sm shrink-0"
                :style="{ background: item.color }"
              />
              <span class="flex-1 text-[var(--jafa-text-muted)] text-[calc(12px*var(--jafa-text-scale,1))] truncate">{{ item.name }}</span>
              <span class="text-[var(--jafa-text)] tabular-nums text-[calc(12px*var(--jafa-text-scale,1))] font-medium">{{ cs }}{{ item.total.toFixed(0) }}</span>
              <span class="text-[var(--jafa-text-muted)] text-[calc(11px*var(--jafa-text-scale,1))] w-8 text-right tabular-nums">
                {{ chartTotal && selectedExpenseNames.has(item.name) ? Math.round(item.total / chartTotal * 100) : 0 }}%
              </span>
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
