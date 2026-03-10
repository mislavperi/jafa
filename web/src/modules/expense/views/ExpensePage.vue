<script setup lang="ts">
import { ref } from 'vue'
import Button from 'primevue/button'
import LatestSpends from '../components/LatestSpends.vue'
import TotalSpendCard from '../components/TotalSpendCard.vue'
import DailySpendChart from '../components/DailySpendChart.vue'
import AddExpenseModal from '../components/AddExpenseModal.vue'

const showModal = ref(false)
</script>

<template>
  <div class="p-2 sm:p-4">
    <div class="flex justify-end mb-3 sm:mb-4">
      <Button label="Add Expense" icon="pi pi-plus" @click="showModal = true" />
    </div>
    <TotalSpendCard class="mb-3 sm:mb-4" />
    <DailySpendChart class="mb-3 sm:mb-4" />
  </div>
  <LatestSpends />
  <AddExpenseModal v-model:visible="showModal" />
</template>
