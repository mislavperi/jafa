<script setup lang="ts">
import { computed } from 'vue'
import { ref } from 'vue'
import { useLogin } from '../composables/useAuth'
import { Form, FormField } from '@primevue/forms'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Message from 'primevue/message'

const showPassword = ref(false)

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

// Receipt metadata
const lastUpdate = new Date(__LAST_UPDATE__)
const stampDate = lastUpdate.toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: '2-digit' })
const stampTime = lastUpdate.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
const version = __APP_VERSION__

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

const floatingNums = computed(() => {
  const samples = ['$12.40', '−$87.40', '$2,148', '$59.99', '+$480', '$10.99', '$1,850', '−$28.50', '$87.40', '$45.00', '−$15.49']
  const out: Array<{ val: string; left: number; delay: number; dur: number; size: number; color: string }> = []
  let seed = 13
  const r = () => { seed = (seed * 9301 + 49297) % 233280; return seed / 233280 }
  for (let i = 0; i < 18; i++) {
    const val = samples[Math.floor(r() * samples.length)]!
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

// PrimeVue passthrough: make InputText look like dashed-underline receipt field
const receiptInputPt = {
  root: {
    class: '!bg-transparent !border-0 !border-b !border-dashed !border-[#8a8878] !rounded-none !text-[#1a1a1a] !font-mono !px-1 !py-1.5 focus:!border-[#1a1a1a] focus:!shadow-none placeholder:!text-[#a8a692]',
  },
}
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
            <linearGradient id="ticker-g" x1="0" x2="1">
              <stop offset="0" stop-color="#f5c518" stop-opacity="0" />
              <stop offset="0.5" stop-color="#f5c518" stop-opacity="0.6" />
              <stop offset="1" stop-color="#f5c518" stop-opacity="0" />
            </linearGradient>
          </defs>
          <path :d="tickerPath" fill="none" stroke="url(#ticker-g)" stroke-width="0.5" vector-effect="non-scaling-stroke" />
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
          ★ EXPENSE TRACKING CO. ★
        </div>
        <div class="mt-3 flex justify-between text-[10.5px] text-[#2d2a22] tracking-wide">
          <span>{{ stampDate }} {{ stampTime }}</span>
          <span>version {{ version }}</span>
        </div>
      </div>

      <hr class="border-0 border-t border-dashed border-[#a8a692] my-2.5" />

      <div class="text-center my-2">
        <div class="text-[18px] font-extrabold tracking-[0.06em] uppercase">Customer Sign-In</div>
        <div class="mt-1 text-[10.5px] tracking-[0.14em] uppercase text-[#3d3a30]">
          Please enter your credentials
        </div>
      </div>

      <hr class="border-0 border-t border-dashed border-[#a8a692] my-2.5" />

      <Form :resolver="resolver" class="flex flex-col gap-3" @submit="handleSubmit">
        <FormField v-slot="$field" name="username">
          <div>
            <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#3d3a30] mb-1 whitespace-nowrap">
              <span>1× Username</span>
              <span class="text-[#1a1a1a] font-bold">REQUIRED</span>
            </label>
            <InputText
              v-bind="$field"
              type="text"
              placeholder="your_username"
              autocomplete="username"
              class="w-full"
              :pt="receiptInputPt"
              unstyled
            />
            <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-700 font-medium">
              {{ $field.error?.message }}
            </div>
          </div>
        </FormField>

        <FormField v-slot="$field" name="password">
          <div>
            <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#3d3a30] mb-1 whitespace-nowrap">
              <span>1× Password</span>
              <span class="text-[#1a1a1a] font-bold">REQUIRED</span>
            </label>
            <div class="relative">
              <InputText
                v-bind="$field"
                :type="showPassword ? 'text' : 'password'"
                placeholder="••••••••"
                autocomplete="current-password"
                class="w-full !pr-8"
                :pt="receiptInputPt"
                unstyled
              />
              <button
                type="button"
                class="absolute right-1 top-1/2 -translate-y-1/2 text-[#2d2a22] hover:text-[#1a1a1a] p-1"
                :aria-label="showPassword ? 'Hide password' : 'Show password'"
                @click="showPassword = !showPassword"
              >
                <i :class="showPassword ? 'pi pi-eye-slash' : 'pi pi-eye'" class="text-[13px]" />
              </button>
            </div>
            <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-700 font-medium">
              {{ $field.error?.message }}
            </div>
          </div>
        </FormField>

        <hr class="border-0 border-t border-dashed border-[#a8a692] my-1" />

        <!-- Itemized summary -->
        <div class="flex flex-col">
          <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
            <span class="text-[#2d2a22] uppercase text-[10.5px] tracking-[0.14em]">Subtotal</span>
            <span class="font-semibold">2 items</span>
          </div>
          <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
            <span class="text-[#2d2a22] uppercase text-[10.5px] tracking-[0.14em]">Security</span>
            <span class="font-semibold">256-bit SSL</span>
          </div>
          <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
            <span class="text-[#2d2a22] uppercase text-[10.5px] tracking-[0.14em]">Session tax</span>
            <span class="font-semibold">€0.00</span>
          </div>
        </div>

        <hr class="border-0 border-t border-dashed border-[#a8a692] my-1" />

        <div class="flex justify-between items-center gap-3 py-1 text-[16px] font-extrabold tracking-[0.04em]">
          <span>TOTAL</span>
          <span>Free Forever</span>
        </div>

        <Message v-if="error" severity="error" :closable="false" class="!mt-1">{{ error.message }}</Message>

        <Button
          type="submit"
          :loading="isPending"
          unstyled
          class="w-full flex items-center justify-center gap-1.5 bg-[#1a1a1a] text-[#f5f1e6] py-3 px-4 font-bold text-[13px] tracking-[0.18em] uppercase hover:opacity-85 disabled:opacity-50 transition mt-1"
        >
          <template v-if="isPending">PROCESSING…</template>
          <template v-else>Sign In <span>→</span></template>
        </Button>
      </Form>

      <hr class="border-0 border-t border-dashed border-[#a8a692] mt-4 mb-2.5" />

      <div class="text-center text-[11px] tracking-[0.16em] uppercase text-[#1a1a1a] font-bold">
        ★ Thank you for budgeting ★
      </div>
      <div class="text-center text-[9.5px] tracking-[0.1em] text-[#3d3a30] mt-1">
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
      <div class="text-center text-[11px] tracking-[0.3em] font-semibold">
        JAFA · 0042 · {{ version }}
      </div>
    </div>

    <div class="relative z-10 text-center text-[11px] tracking-[0.14em] uppercase text-zinc-400">
      New customer?
      <RouterLink to="/register" class="text-[#f5c518] hover:text-[#f97316] font-semibold underline underline-offset-[3px] tracking-[0.04em] uppercase ml-1">
        Open an account →
      </RouterLink>
    </div>
  </div>
</template>

