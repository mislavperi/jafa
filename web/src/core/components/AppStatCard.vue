<script setup lang="ts">
defineProps<{
  label: string
  value?: string | number
  subtitle?: string
  loading?: boolean
  icon?: string
}>()
</script>

<template>
  <div class="border border-surface rounded bg-white dark:bg-surface-800 px-4 py-3 flex items-center gap-4">
    <div v-if="icon" class="w-9 h-9 rounded bg-brand/10 flex items-center justify-center shrink-0">
      <i :class="icon" class="text-brand text-base" />
    </div>
    <div class="min-w-0">
      <p class="text-[0.65rem] font-semibold uppercase tracking-widest text-surface-400 mb-1">{{ label }}</p>
      <div v-if="loading" class="h-7 w-24 bg-surface-200 dark:bg-surface-700 rounded animate-pulse" />
      <p v-else class="text-2xl font-bold tracking-tight leading-none">{{ value ?? '—' }}</p>
      <p v-if="subtitle" class="text-xs text-surface-400 mt-1">{{ subtitle }}</p>
    </div>
  </div>
</template>
