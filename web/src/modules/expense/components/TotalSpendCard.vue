<script setup lang="ts">
import Card from 'primevue/card'
import Skeleton from 'primevue/skeleton'
import { useMonthlyTotal } from '../composables/useExpenses'

const { data: monthlyTotal, isLoading } = useMonthlyTotal()

const currentMonth = new Date().toLocaleString('default', { month: 'long', year: 'numeric' })
</script>

<template>
  <Card class="w-full">
    <template #title>
      Total Spend This Month
    </template>
    <template #subtitle>
      {{ currentMonth }}
    </template>
    <template #content>
      <Skeleton v-if="isLoading" width="8rem" height="2.5rem" />
      <div v-else class="text-3xl font-bold text-primary">
        ${{ monthlyTotal?.total.toFixed(2) ?? '0.00' }}
      </div>
    </template>
  </Card>
</template>
