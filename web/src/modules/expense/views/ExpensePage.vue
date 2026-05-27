<script setup lang="ts">
import { ref } from 'vue'
import Button from 'primevue/button'
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import AppStatCard from '@/core/components/AppStatCard.vue'
import LatestSpends from '../components/LatestSpends.vue'
import TotalSpendCard from '../components/TotalSpendCard.vue'
import AddExpenseModal from '../components/AddExpenseModal.vue'

const showModal = ref(false)
</script>

<template>
  <Root>
    <div
      class="flex flex-col gap-5 h-full min-w-0 p-8 overflow-auto"
    >
      <AppPageHeader title="Dashboard" subtitle="Your spending at a glance">
        <Button
          label="Add Expense"
          icon="pi pi-plus"
          size="small"
          @click="showModal = true"
        />
      </AppPageHeader>

      <div class="grid grid-cols-3 gap-4">
        <TotalSpendCard />
        <AppStatCard label="Budget" icon="pi pi-chart-pie" tone="brand" />
        <AppStatCard label="Savings" icon="pi pi-piggy-bank" tone="positive" />
      </div>

      <LatestSpends class="flex-1 min-h-0" />
    </div>
    <AddExpenseModal v-model:visible="showModal" />
  </Root>
</template>
