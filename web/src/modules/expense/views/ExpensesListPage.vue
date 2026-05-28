<script setup lang="ts">
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import LatestSpends from '../components/LatestSpends.vue'
import Button from 'primevue/button'
import AddExpenseModal from '../components/AddExpenseModal.vue'
import ReceiptScannerModal from '../components/ReceiptScannerModal.vue'
import { ref } from 'vue'
const showModal = ref(false)
const showScanner = ref(false)
</script>
<template>
  <Root>
    <div class="flex flex-col gap-5 h-full min-w-0 p-8 overflow-auto">
      <AppPageHeader title="Expenses" subtitle="All your transactions">
        <Button label="Scan Receipt" icon="pi pi-receipt" size="small" severity="secondary" @click="showScanner = true" />
        <Button label="Add Expense" icon="pi pi-plus" size="small" @click="showModal = true" />
      </AppPageHeader>
      <LatestSpends class="flex-1 min-h-0" />
      <AddExpenseModal v-model:visible="showModal" />
      <ReceiptScannerModal v-model:visible="showScanner" />
    </div>
  </Root>
</template>
