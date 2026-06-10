<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch, nextTick } from 'vue'
import type { CSSProperties } from 'vue'
import Button from 'primevue/button'
import { TOUR_STEPS, useOnboardingStore } from '@/stores/onboarding'

const tour = useOnboardingStore()
// The store clamps stepIndex to valid bounds, so the lookup never misses.
const step = computed(() => TOUR_STEPS[tour.stepIndex]!)
const isFirst = computed(() => tour.stepIndex === 0)
const isLast = computed(() => tour.stepIndex === TOUR_STEPS.length - 1)

const SPOT_PAD = 8
const POP_W = 340
const POP_H = 230 // estimate used only for viewport clamping
const GAP = 14
const MARGIN = 12

const rect = ref<{ top: number; left: number; width: number; height: number } | null>(null)
let targetEl: HTMLElement | null = null
let rafId = 0

// The dashboard renders async data (queries, charts), so the layout shifts
// after the tour starts. Track the target every frame instead of trying to
// enumerate every event that can move it.
function syncRect() {
  if (!targetEl || !targetEl.isConnected) {
    targetEl = step.value?.target
      ? document.querySelector<HTMLElement>(`[data-tour="${step.value.target}"]`)
      : null
  }
  const r = targetEl?.getBoundingClientRect() ?? null
  const cur = rect.value
  if (!r) {
    if (cur) rect.value = null
  } else if (
    !cur ||
    cur.top !== r.top ||
    cur.left !== r.left ||
    cur.width !== r.width ||
    cur.height !== r.height
  ) {
    rect.value = { top: r.top, left: r.left, width: r.width, height: r.height }
  }
  rafId = requestAnimationFrame(syncRect)
}

watch(
  () => [tour.active, tour.stepIndex],
  async () => {
    cancelAnimationFrame(rafId)
    targetEl = null
    rect.value = null
    if (!tour.active) return
    await nextTick()
    targetEl = step.value?.target
      ? document.querySelector<HTMLElement>(`[data-tour="${step.value.target}"]`)
      : null
    targetEl?.scrollIntoView({ block: 'nearest' })
    syncRect()
  },
  { immediate: true },
)

function clamp(v: number, lo: number, hi: number) {
  return Math.min(Math.max(v, lo), hi)
}

const spotlightStyle = computed<CSSProperties>(() => {
  const r = rect.value
  if (!r) return { display: 'none' }
  return {
    top: `${r.top - SPOT_PAD}px`,
    left: `${r.left - SPOT_PAD}px`,
    width: `${r.width + SPOT_PAD * 2}px`,
    height: `${r.height + SPOT_PAD * 2}px`,
  }
})

const popoverStyle = computed<CSSProperties>(() => {
  const r = rect.value
  if (!r || step.value.placement === 'center') {
    return { top: '50%', left: '50%', transform: 'translate(-50%, -50%)' }
  }
  const vw = window.innerWidth
  const vh = window.innerHeight
  const gap = GAP + SPOT_PAD

  let p = step.value.placement
  if (p === 'right' && r.left + r.width + gap + POP_W > vw - MARGIN) p = 'bottom'
  if (p === 'left' && r.left - gap - POP_W < MARGIN) p = 'bottom'
  if (p === 'bottom' && r.top + r.height + gap + POP_H > vh - MARGIN) p = 'top'
  if (p === 'top' && r.top - gap - POP_H < MARGIN) p = 'bottom'

  const sideTop = `${clamp(r.top + r.height / 2 - POP_H / 2, MARGIN, vh - POP_H - MARGIN)}px`
  const centeredLeft = `${clamp(r.left + r.width / 2 - POP_W / 2, MARGIN, vw - POP_W - MARGIN)}px`
  switch (p) {
    case 'right':
      return { top: sideTop, left: `${r.left + r.width + gap}px` }
    case 'left':
      return { top: sideTop, left: `${r.left - gap - POP_W}px` }
    case 'top':
      return { bottom: `${vh - r.top + gap}px`, left: centeredLeft }
    default:
      return { top: `${r.top + r.height + gap}px`, left: centeredLeft }
  }
})

function onKeydown(e: KeyboardEvent) {
  if (!tour.active) return
  if (e.key === 'Escape') tour.finish()
  else if (e.key === 'ArrowRight') tour.next()
  else if (e.key === 'ArrowLeft') tour.back()
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
  cancelAnimationFrame(rafId)
})
</script>

<template>
  <Teleport to="body">
    <div v-if="tour.active" class="fixed inset-0 z-[1200]">
      <!-- Full dim for centered steps; spotlight casts the dim otherwise -->
      <div v-if="!rect" class="absolute inset-0 bg-black/55" />
      <div
        v-else
        class="fixed rounded-xl pointer-events-none transition-all duration-200"
        style="box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.55); outline: 2px solid var(--jafa-accent)"
        :style="spotlightStyle"
      />

      <div
        class="fixed w-[340px] bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 shadow-2xl flex flex-col gap-3 transition-all duration-200"
        :style="popoverStyle"
        role="dialog"
        aria-modal="true"
        :aria-label="step.title"
      >
        <div class="flex items-start justify-between gap-3">
          <h3 class="text-[calc(15px*var(--jafa-text-scale,1))] font-semibold text-[var(--jafa-text)]">
            {{ step.title }}
          </h3>
          <button
            class="w-7 h-7 -mt-1 -mr-1 inline-flex items-center justify-center rounded-md shrink-0 text-[var(--jafa-text-muted)] hover:bg-[var(--jafa-surface-2)] hover:text-[var(--jafa-text)] transition"
            aria-label="Skip tour"
            @click="tour.finish"
          >
            <i class="pi pi-times text-[calc(12px*var(--jafa-text-scale,1))]" />
          </button>
        </div>

        <p class="text-[calc(13px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] leading-relaxed">
          {{ step.text }}
        </p>

        <div class="flex items-center justify-between mt-1">
          <span class="text-[calc(11px*var(--jafa-text-scale,1))] tabular-nums text-[var(--jafa-text-muted)]">
            {{ tour.stepIndex + 1 }} / {{ TOUR_STEPS.length }}
          </span>
          <div class="flex items-center gap-2">
            <Button
              v-if="isFirst"
              label="Skip tour"
              severity="secondary"
              size="small"
              text
              @click="tour.finish"
            />
            <Button v-else label="Back" severity="secondary" size="small" @click="tour.back" />
            <Button
              :label="isFirst ? 'Show me around' : isLast ? 'Finish' : 'Next'"
              size="small"
              @click="tour.next"
            />
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
