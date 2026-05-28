<script setup lang="ts">
import { toRef, computed } from 'vue'
import { useExpenseTags } from '../composables/useTags'

const props = defineProps<{ expenseId: number }>()
const { data: rawTags, isLoading } = useExpenseTags(toRef(props, 'expenseId'))
const tags = computed(() => {
  if (!rawTags.value) return rawTags.value
  const seen = new Set<number>()
  return rawTags.value.filter((t) => { if (seen.has(t.id)) return false; seen.add(t.id); return true })
})

const TAG_COLORS = ['#f5c518', '#f97316', '#22c55e', '#3b82f6', '#a855f7', '#ec4899', '#14b8a6', '#ef4444']

function tagColor(tagId: number) {
  return TAG_COLORS[tagId % TAG_COLORS.length]!
}
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
