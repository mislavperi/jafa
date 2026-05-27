<script setup lang="ts">
import { computed } from 'vue'
import { useLogin } from '../composables/useAuth'

const { mutate: login, isPending, error } = useLogin()

function resolver({ values }: { values: Record<string, string> }) {
  const errors: Record<string, { message: string }[]> = {}
  if (!values.username) errors.username = [{ message: 'Username is required' }]
  if (!values.password) errors.password = [{ message: 'Password is required' }]
  return { values, errors }
}

function handleSubmit({ valid, values }: { valid: boolean; values: Record<string, string> }) {
  if (valid) {
    login({ username: values.username, password: values.password })
  }
}

// Stable receipt metadata for this session
const now = new Date()
const stampDate = now
  .toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: '2-digit' })
  .replace(/\//g, '/')
const stampTime = now.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
const txnId = 'TXN-' + Math.random().toString(36).slice(2, 8).toUpperCase()

// Deterministic barcode bars
const barcodeBars = computed(() => {
  const arr: number[] = []
  let seed = 7
  for (let i = 0; i < 52; i++) {
    seed = (seed * 9301 + 49297) % 233280
    const r = seed / 233280
    arr.push(r < 0.5 ? 1 : r < 0.85 ? 2 : 3)
  }
  return arr
})

// Floating background numbers
const floatingNums = computed(() => {
  const samples = ['$12.40', '−$87.40', '$2,148', '$59.99', '+$480', '$10.99', '$1,850', '−$28.50', '$87.40', '$45.00', '−$15.49']
  const out: Array<{ val: string; left: number; delay: number; dur: number; size: number; color: string }> = []
  let seed = 13
  const r = () => { seed = (seed * 9301 + 49297) % 233280; return seed / 233280 }
  for (let i = 0; i < 18; i++) {
    const val = samples[Math.floor(r() * samples.length)]
    out.push({
      val,
      left: r() * 100,
      delay: r() * -30,
      dur: 20 + r() * 25,
      size: 11 + r() * 4,
      color: val.startsWith('+')
        ? 'rgba(34, 197, 94, 0.45)'
        : val.startsWith('−')
          ? 'rgba(239, 68, 68, 0.4)'
          : 'rgba(245, 197, 24, 0.4)',
    })
  }
  return out
})

// Subtle ticker line path
const tickerPath = computed(() => {
  const points: [number, number][] = []
  let seed = 99
  let y = 60
  for (let x = 0; x <= 100; x += 4) {
    seed = (seed * 9301 + 49297) % 233280
    y += (seed / 233280 - 0.5) * 16
    y = Math.max(20, Math.min(100, y))
    points.push([x, y])
  }
  return 'M ' + points.map(p => `${p[0]} ${p[1]}`).join(' L ')
})
</script>

<template>
  <div
    class="relative flex flex-col items-center justify-center min-h-screen overflow-hidden px-6 py-12 gap-7"
    style="background: radial-gradient(ellipse at 50% 40%, #1a1a1f 0%, #0a0a0b 60%);"
  >
    <!-- Animated backdrop -->
    <div class="absolute inset-0 z-0 pointer-events-none overflow-hidden" aria-hidden="true">
      <div
        v-for="(n, i) in floatingNums"
        :key="i"
        class="absolute drift opacity-25 whitespace-nowrap"
        :style="{
          left: `${n.left}%`,
          bottom: '-10%',
          animationDelay: `${n.delay}s`,
          animationDuration: `${n.dur}s`,
          fontSize: `${n.size}px`,
          color: n.color,
          fontFamily: `'Geist Mono', monospace`,
        }"
      >{{ n.val }}</div>

      <div class="absolute left-0 right-0 h-px" style="bottom: 14%;">
        <svg class="w-full h-[120px] block opacity-30" viewBox="0 0 100 120" preserveAspectRatio="none">
          <defs>
            <linearGradient id="ticker-g" x1="0" x2="1">
              <stop offset="0" stop-color="#f5c518" stop-opacity="0"/>
              <stop offset="0.5" stop-color="#f5c518" stop-opacity="0.6"/>
              <stop offset="1" stop-color="#f5c518" stop-opacity="0"/>
            </linearGradient>
          </defs>
          <path :d="tickerPath" fill="none" stroke="url(#ticker-g)" stroke-width="0.5" vector-effect="non-scaling-stroke"/>
        </svg>
      </div>
    </div>

    <!-- The receipt -->
    <Form :resolver="resolver" @submit="handleSubmit" class="receipt-shell relative z-10">
      <!-- Header -->
      <div class="text-center mb-4">
        <div class="inline-flex items-center gap-2 font-bold text-[22px] tracking-[0.12em] text-[#1a1a1a]">
          <img src="/icon.png" class="w-7 h-7" alt=""/>
          <span>JAFA</span>
        </div>
        <div class="mt-2 text-[10px] tracking-[0.18em] uppercase text-[#6b6b66]">
          ★ EXPENSE TRACKING CO. ★
        </div>
        <div class="mt-3 flex justify-between text-[10.5px] text-[#6b6b66] tracking-wide">
          <span>{{ stampDate }}  {{ stampTime }}</span>
          <span>#{{ txnId }}</span>
        </div>
      </div>

      <hr class="receipt-divider"/>

      <div class="text-center my-2">
        <div class="text-[18px] font-extrabold tracking-[0.06em] uppercase">Customer Sign-In</div>
        <div class="mt-1 text-[10.5px] tracking-[0.14em] uppercase text-[#8a8878]">
          Please enter your credentials
        </div>
      </div>

      <hr class="receipt-divider"/>

      <!-- Username field -->
      <FormField v-slot="$field" name="username">
        <div class="mb-3">
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#8a8878] mb-1 whitespace-nowrap">
            <span>1× Username</span>
            <span class="text-[#c7c5b8]">REQ</span>
          </label>
          <input
            v-bind="$field"
            type="text"
            class="receipt-input"
            placeholder="your_username"
            autocomplete="username"
          />
          <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-600 font-medium">
            {{ $field.error?.message }}
          </div>
        </div>
      </FormField>

      <!-- Password field -->
      <FormField v-slot="$field" name="password">
        <div class="mb-3">
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#8a8878] mb-1 whitespace-nowrap">
            <span>1× Password</span>
            <span class="text-[#c7c5b8]">REQ</span>
          </label>
          <input
            v-bind="$field"
            type="password"
            class="receipt-input"
            placeholder="••••••••"
            autocomplete="current-password"
          />
          <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-600 font-medium">
            {{ $field.error?.message }}
          </div>
        </div>
      </FormField>

      <hr class="receipt-divider"/>

      <!-- Itemized "summary" -->
      <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px] whitespace-nowrap">
        <span class="text-[#6b6b66] uppercase text-[10.5px] tracking-[0.14em]">Subtotal</span>
        <span class="font-semibold text-[#1a1a1a]">2 items</span>
      </div>
      <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px] whitespace-nowrap">
        <span class="text-[#6b6b66] uppercase text-[10.5px] tracking-[0.14em]">Security</span>
        <span class="font-semibold text-[#1a1a1a]">256-bit SSL</span>
      </div>
      <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px] whitespace-nowrap">
        <span class="text-[#6b6b66] uppercase text-[10.5px] tracking-[0.14em]">Session tax</span>
        <span class="font-semibold text-[#1a1a1a]">$0.00</span>
      </div>

      <hr class="receipt-divider"/>

      <!-- TOTAL -->
      <div class="flex justify-between items-center gap-3 py-1 text-[16px] font-extrabold tracking-[0.04em] whitespace-nowrap">
        <span>TOTAL</span>
        <span>Free Forever</span>
      </div>

      <!-- Error -->
      <div v-if="error" class="mt-2 px-3 py-2 bg-red-50 border border-red-200 text-red-700 text-[11.5px] leading-snug">
        {{ error.message }}
      </div>

      <!-- Submit -->
      <button type="submit" class="receipt-btn mt-2" :disabled="isPending">
        <template v-if="isPending">PROCESSING…</template>
        <template v-else>Sign In <span class="ml-1">→</span></template>
      </button>

      <hr class="receipt-divider" style="margin-top: 18px;"/>

      <div class="text-center text-[11px] tracking-[0.16em] uppercase text-[#4a4a44] font-bold">
        ★ Thank you for budgeting ★
      </div>
      <div class="text-center text-[9.5px] tracking-[0.1em] text-[#8a8878] mt-1">
        Every dollar tracked is a dollar earned
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
      <div class="text-center text-[11px] tracking-[0.3em] text-[#1a1a1a] font-semibold">
        JAFA · 0042 · {{ txnId.split('-')[1] }}
      </div>
    </Form>

    <!-- Sign-up nudge -->
    <div class="relative z-10 text-center text-[11px] tracking-[0.14em] uppercase text-zinc-400 whitespace-nowrap">
      New customer?
      <RouterLink to="/register" class="receipt-link ml-1">Open an account →</RouterLink>
    </div>
  </div>
</template>
