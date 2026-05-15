<script setup lang="ts">
import { toRef, computed } from 'vue'
import { useExpenseTags } from '../composables/useTags'
import { useDarkModeStore } from '@/stores/darkMode'

const props = defineProps<{ expenseId: number }>()
const { data: rawTags, isLoading } = useExpenseTags(toRef(props, 'expenseId'))
const tags = computed(() => {
  if (!rawTags.value) return rawTags.value
  const seen = new Set<number>()
  return rawTags.value.filter((t) => { if (seen.has(t.id)) return false; seen.add(t.id); return true })
})
const darkMode = useDarkModeStore()

const LIGHT = [
  { bg: '#fef3c7', color: '#b45309' },
  { bg: '#d1fae5', color: '#065f46' },
  { bg: '#dbeafe', color: '#1e40af' },
  { bg: '#fce7f3', color: '#9d174d' },
  { bg: '#ede9fe', color: '#5b21b6' },
  { bg: '#fee2e2', color: '#991b1b' },
  { bg: '#ccfbf1', color: '#134e4a' },
]

const DARK = [
  { bg: '#78350f', color: '#fcd34d' },
  { bg: '#064e3b', color: '#6ee7b7' },
  { bg: '#1e3a8a', color: '#93c5fd' },
  { bg: '#831843', color: '#f9a8d4' },
  { bg: '#4c1d95', color: '#c4b5fd' },
  { bg: '#7f1d1d', color: '#fca5a5' },
  { bg: '#134e4a', color: '#5eead4' },
]

function tagStyle(tagId: number) {
  const palette = darkMode.isDark ? DARK : LIGHT
  const c = palette[tagId % palette.length]!
  return { backgroundColor: c.bg, color: c.color }
}
</script>

<template>
  <div class="flex flex-wrap gap-1">
    <span v-if="isLoading" class="text-xs text-surface-400">...</span>
    <span v-else-if="!tags?.length" class="text-xs text-surface-400">—</span>
    <span
      v-for="tag in tags"
      :key="tag.id"
      class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium whitespace-nowrap"
      :style="tagStyle(tag.id)"
    >
      {{ tag.name }}
    </span>
  </div>
</template>
