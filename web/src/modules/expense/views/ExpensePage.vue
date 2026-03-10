<script setup lang="ts">
import Root from '@/core/views/Root.vue'
import LatestSpends from '../components/LatestSpends.vue'
import TotalSpendCard from '../components/TotalSpendCard.vue'
import DailySpendChart from '../components/DailySpendChart.vue'
</script>

<template>
  <Root>
    <div class="p-2 sm:p-4">
      <TotalSpendCard class="mb-3 sm:mb-4" />
      <DailySpendChart class="mb-3 sm:mb-4" />
    </div>
    <LatestSpends />
  </Root>
</template>
