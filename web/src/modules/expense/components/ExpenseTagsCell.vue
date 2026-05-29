<script setup lang="ts">
import { toRef, computed } from 'vue'
import { useExpenseTags } from '../composables/useTags'
import { tagColor } from '../constants'

const props = defineProps<{ expenseId: number }>()
const { data: rawTags, isLoading } = useExpenseTags(toRef(props, 'expenseId'))
const tags = computed(() => {
  if (!rawTags.value) return []
  const seen = new Set<number>()
  return rawTags.value.filter((t) => { if (seen.has(t.id)) return false; seen.add(t.id); return true })
})
</script>

<template>
  <div class="flex flex-wrap gap-1.5">
    <span v-if="isLoading" class="text-[calc(11px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)]">...</span>
    <span v-else-if="!tags?.length" class="text-[calc(12px*var(--jafa-text-scale,1))] text-[var(--jafa-text-dim)]">—</span>
    <span
      v-for="tag in tags"
      :key="tag.id"
      class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-[var(--jafa-surface-3)] text-[calc(12px*var(--jafa-text-scale,1))] font-medium whitespace-nowrap text-zinc-300"
    >
      <span class="w-1.5 h-1.5 rounded-full" :style="{ background: tagColor(tag.id) }" />
      {{ tag.name }}
    </span>
  </div>
</template>
