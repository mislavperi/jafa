<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<template>
  <div class="flex items-start justify-between gap-6 shrink-0">
    <div>
      <h1 class="text-[18px] font-semibold tracking-tight text-white whitespace-nowrap">
        {{ title }}
      </h1>
      <p v-if="subtitle" class="text-[13px] text-zinc-400 mt-0.5">{{ subtitle }}</p>
    </div>
    <div class="flex items-center gap-2">
      <slot />
    </div>
  </div>
</template>
