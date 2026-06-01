<script setup lang="ts">
import { useId } from 'vue'
import { useReceiptDecor } from '../composables/useReceiptDecor'

const props = withDefaults(
  defineProps<{
    tagline: string
    title: string
    subtitle: string
    thankYou: string
    thankYouSub: string
    barcodeLabel: string
    currency?: string
    barcodeSeed?: number
    floatSeed?: number
    tickerSeed?: number
  }>(),
  { currency: '$', barcodeSeed: 7, floatSeed: 13, tickerSeed: 99 },
)

const { stampDate, stampTime, version, barcodeBars, floatingNums, tickerPath } = useReceiptDecor({
  currency: props.currency,
  barcodeSeed: props.barcodeSeed,
  floatSeed: props.floatSeed,
  tickerSeed: props.tickerSeed,
})

// Unique gradient id so multiple receipts on a page never clash.
const gradientId = `ticker-g-${useId()}`
</script>

<template>
  <div
    class="relative flex flex-col items-center justify-center min-h-screen overflow-hidden px-6 py-12 gap-7 font-mono"
    style="background: radial-gradient(ellipse at 50% 40%, #1a1a1f 0%, #0a0a0b 60%);"
  >
    <!-- Animated backdrop -->
    <div class="absolute inset-0 z-0 pointer-events-none overflow-hidden" aria-hidden="true">
      <div
        v-for="(n, i) in floatingNums"
        :key="i"
        class="absolute drift opacity-25 whitespace-nowrap font-mono"
        :style="{
          left: `${n.left}%`,
          bottom: '-10%',
          animationDelay: `${n.delay}s`,
          animationDuration: `${n.dur}s`,
          fontSize: `${n.size}px`,
          color: n.color,
        }"
      >{{ n.val }}</div>

      <div class="absolute left-0 right-0 h-px" style="bottom: 14%;">
        <svg class="w-full h-[120px] block opacity-30" viewBox="0 0 100 120" preserveAspectRatio="none">
          <defs>
            <linearGradient :id="gradientId" x1="0" x2="1">
              <stop offset="0" stop-color="#f5c518" stop-opacity="0" />
              <stop offset="0.5" stop-color="#f5c518" stop-opacity="0.6" />
              <stop offset="1" stop-color="#f5c518" stop-opacity="0" />
            </linearGradient>
          </defs>
          <path :d="tickerPath" fill="none" :stroke="`url(#${gradientId})`" stroke-width="0.5" vector-effect="non-scaling-stroke" />
        </svg>
      </div>
    </div>

    <!-- Receipt -->
    <div class="relative z-10 w-[360px] max-w-full bg-[#f5f1e6] text-[#1a1a1a] px-7 pt-7 pb-6 rounded-sm shadow-[0_30px_80px_-20px_rgba(0,0,0,0.6),0_8px_24px_-10px_rgba(0,0,0,0.4)] font-mono">
      <!-- Header -->
      <div class="text-center mb-4">
        <div class="inline-flex items-center gap-2 font-bold text-[22px] tracking-[0.12em]">
          <img src="/icon.png" class="w-7 h-7" alt="" />
          <span>JAFA</span>
        </div>
        <div class="mt-2 text-[10px] tracking-[0.18em] uppercase text-[#2d2a22]">
          {{ tagline }}
        </div>
        <div class="mt-3 flex justify-between text-[10.5px] text-[#2d2a22] tracking-wide">
          <span>{{ stampDate }} {{ stampTime }}</span>
          <span>version {{ version }}</span>
        </div>
      </div>

      <hr class="border-0 border-t border-dashed border-[#a8a692] my-2.5" />

      <div class="text-center my-2">
        <div class="text-[18px] font-extrabold tracking-[0.06em] uppercase">{{ title }}</div>
        <div class="mt-1 text-[10.5px] tracking-[0.14em] uppercase text-[#3d3a30]">
          {{ subtitle }}
        </div>
      </div>

      <hr class="border-0 border-t border-dashed border-[#a8a692] my-2.5" />

      <!-- Form body (fields, itemized summary, total, submit) -->
      <slot />

      <hr class="border-0 border-t border-dashed border-[#a8a692] mt-4 mb-2.5" />

      <div class="text-center text-[11px] tracking-[0.16em] uppercase text-[#1a1a1a] font-bold">
        {{ thankYou }}
      </div>
      <div class="text-center text-[9.5px] tracking-[0.1em] text-[#3d3a30] mt-1">
        {{ thankYouSub }}
      </div>

      <!-- Barcode -->
      <div class="flex justify-center gap-[1.5px] my-3 h-[38px] items-stretch">
        <span
          v-for="(w, i) in barcodeBars"
          :key="i"
          class="block bg-[#1a1a1a]"
          :style="{ width: `${w}px` }"
        />
      </div>
      <div class="text-center text-[11px] tracking-[0.3em] font-semibold">
        JAFA · {{ barcodeLabel }} · {{ version }}
      </div>
    </div>

    <slot name="footer" />
  </div>
</template>
