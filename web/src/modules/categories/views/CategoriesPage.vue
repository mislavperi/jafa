<script setup lang="ts">
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import { useCategoryBreakdown } from '@/modules/reports/composables/useReports'
import { useThemeStore } from '@/stores/theme'
import { formatCurrency } from '@/core/currency'

const { data: categories, isLoading } = useCategoryBreakdown()
const theme = useThemeStore()

function fmt(val: number) {
  return formatCurrency(Math.abs(val), theme.currency)
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
          v-for="cat in categories"
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
