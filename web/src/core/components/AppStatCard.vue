<script setup lang="ts">
defineProps<{
  label: string
  value?: string | number
  subtitle?: string
  loading?: boolean
  icon?: string
  tone?: 'brand' | 'positive' | 'muted'
}>()
</script>

<template>
  <div
    class="bg-[#131316] border border-[#26262c] rounded-[14px] px-5 py-4 flex items-center gap-4"
  >
    <div
      v-if="icon"
      class="w-11 h-11 rounded-[10px] flex items-center justify-center shrink-0"
      :class="{
        'bg-[#3d2f0a] text-[#f5c518]': !tone || tone === 'brand',
        'bg-emerald-500/10 text-emerald-400': tone === 'positive',
        'bg-[#1f1f24] text-zinc-400': tone === 'muted',
      }"
    >
      <i :class="icon" class="text-[18px]" />
    </div>
    <div class="min-w-0 flex-1">
      <p class="text-[11px] font-semibold uppercase tracking-[0.08em] text-zinc-400 mb-1">
        {{ label }}
      </p>
      <div v-if="loading" class="h-7 w-24 bg-[#1f1f24] rounded animate-pulse" />
      <p v-else class="text-[22px] font-bold tracking-tight leading-none tabular text-white">
        {{ value ?? '—' }}
      </p>
      <p v-if="subtitle" class="text-[12px] text-zinc-500 mt-1">{{ subtitle }}</p>
    </div>
  </div>
</template>
