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
    class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] px-5 py-4 flex items-center gap-4"
  >
    <div
      v-if="icon"
      class="w-11 h-11 rounded-[10px] flex items-center justify-center shrink-0"
      :class="{
        'bg-[var(--jafa-accent)]/15 text-[var(--jafa-accent)]': !tone || tone === 'brand',
        'bg-emerald-500/10 text-emerald-400': tone === 'positive',
        'bg-[var(--jafa-surface-3)] text-[var(--jafa-text-muted)]': tone === 'muted',
      }"
    >
      <i :class="icon" class="text-[calc(18px*var(--jafa-text-scale,1))]" />
    </div>
    <div class="min-w-0 flex-1">
      <p class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.08em] text-[var(--jafa-text-muted)] mb-1">
        {{ label }}
      </p>
      <div v-if="loading" class="h-7 w-24 bg-[var(--jafa-surface-3)] rounded animate-pulse" />
      <p v-else class="text-[calc(22px*var(--jafa-text-scale,1))] font-bold tracking-tight leading-none tabular text-[var(--jafa-text)]">
        {{ value ?? '—' }}
      </p>
      <p v-if="subtitle" class="text-[calc(12px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] mt-1">{{ subtitle }}</p>
    </div>
  </div>
</template>
