<script setup lang="ts">
import AppStatCard from '@/core/components/AppStatCard.vue'
import { useMonthlyTotal } from '../composables/useExpenses'
import { useThemeStore } from '@/stores/theme'
import { formatCurrency } from '@/core/currency'

const { data: monthlyTotal, isLoading } = useMonthlyTotal()
const theme = useThemeStore()
const currentMonth = new Date().toLocaleString('default', { month: 'long', year: 'numeric' })
</script>

<template>
  <AppStatCard
    label="Spent"
    :value="formatCurrency(monthlyTotal?.total ?? 0, theme.currency)"
    :subtitle="currentMonth"
    :loading="isLoading"
    icon="pi pi-wallet"
  />
</template>
