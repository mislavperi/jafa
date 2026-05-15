<script setup lang="ts">
import AppStatCard from '@/core/components/AppStatCard.vue'
import { useMonthlyTotal } from '../composables/useExpenses'

const { data: monthlyTotal, isLoading } = useMonthlyTotal()
const currentMonth = new Date().toLocaleString('default', { month: 'long', year: 'numeric' })
</script>

<template>
  <AppStatCard
    label="Spent"
    :value="monthlyTotal ? `$${monthlyTotal.total.toFixed(2)}` : '$0.00'"
    :subtitle="currentMonth"
    :loading="isLoading"
    icon="pi pi-wallet"
  />
</template>
