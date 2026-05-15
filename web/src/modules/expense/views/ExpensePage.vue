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
      class="grid flex-1 min-w-0 p-4 gap-4 overflow-hidden"
      style="grid-template-rows: auto auto 1fr"
    >
      <AppPageHeader title="Dashboard" subtitle="Your expenses at a glance">
        <Button label="Add Expense" icon="pi pi-plus" size="small" @click="showModal = true" />
      </AppPageHeader>
      <div class="grid grid-cols-3 gap-4">
        <TotalSpendCard />
        <AppStatCard label="Budget" icon="pi pi-chart-pie" />
        <AppStatCard label="Savings" icon="pi pi-piggy-bank" />
      </div>
      <LatestSpends class="h-full" />
    </div>
    <AddExpenseModal v-model:visible="showModal" />
  </Root>
</template>
