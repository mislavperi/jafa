<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<template>
  <div class="flex items-center justify-between pb-3 border-b border-surface shrink-0">
    <div>
      <h1 class="text-base font-bold tracking-tight">{{ title }}</h1>
      <p v-if="subtitle" class="text-xs text-surface-400 mt-0.5">{{ subtitle }}</p>
    </div>
    <div class="flex items-center gap-2">
      <slot />
    </div>
  </div>
</template>
