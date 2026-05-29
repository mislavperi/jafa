<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<template>
  <div class="flex items-start justify-between gap-6 shrink-0">
    <div>
      <h1 class="text-[calc(18px*var(--jafa-text-scale,1))] font-semibold tracking-tight text-[var(--jafa-text)] whitespace-nowrap">
        {{ title }}
      </h1>
      <p v-if="subtitle" class="text-[calc(13px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] mt-0.5">{{ subtitle }}</p>
    </div>
    <div class="flex items-center gap-2">
      <slot />
    </div>
  </div>
</template>
