<script setup lang="ts">
import { computed } from 'vue'
import { useRegister } from '../composables/useAuth'

const { mutate: register, isPending, error } = useRegister()

function resolver({ values }: { values: Record<string, string> }) {
  const errors: Record<string, { message: string }[]> = {}
  if (!values.username) errors.username = [{ message: 'Username is required' }]
  if (!values.password) errors.password = [{ message: 'Password is required' }]
  if (values.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(values.email)) {
    errors.email = [{ message: 'Please enter a valid email address' }]
  }
  return { values, errors }
}

function handleSubmit({ valid, values }: { valid: boolean; values: Record<string, string> }) {
  if (valid) {
    register({
      username: values.username,
      password: values.password,
      first_name: values.first_name,
      last_name: values.last_name,
      email: values.email,
    })
  }
}

const now = new Date()
const stampDate = now
  .toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: '2-digit' })
  .replace(/\//g, '/')
const stampTime = now.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
const txnId = 'TXN-' + Math.random().toString(36).slice(2, 8).toUpperCase()

const barcodeBars = computed(() => {
  const arr: number[] = []
  let seed = 31
  for (let i = 0; i < 52; i++) {
    seed = (seed * 9301 + 49297) % 233280
    const r = seed / 233280
    arr.push(r < 0.5 ? 1 : r < 0.85 ? 2 : 3)
  }
  return arr
})
</script>

<template>
  <div
    class="relative flex flex-col items-center justify-center min-h-screen overflow-hidden px-6 py-12 gap-7"
    style="background: radial-gradient(ellipse at 50% 40%, #1a1a1f 0%, #0a0a0b 60%);"
  >
    <Form :resolver="resolver" @submit="handleSubmit" class="receipt-shell relative z-10">
      <!-- Header -->
      <div class="text-center mb-4">
        <div class="inline-flex items-center gap-2 font-bold text-[22px] tracking-[0.12em] text-[#1a1a1a]">
          <img src="/icon.png" class="w-7 h-7" alt=""/>
          <span>JAFA</span>
        </div>
        <div class="mt-2 text-[10px] tracking-[0.18em] uppercase text-[#6b6b66]">
          ★ NEW MEMBER APPLICATION ★
        </div>
        <div class="mt-3 flex justify-between text-[10.5px] text-[#6b6b66] tracking-wide">
          <span>{{ stampDate }}  {{ stampTime }}</span>
          <span>#{{ txnId }}</span>
        </div>
      </div>

      <hr class="receipt-divider"/>

      <div class="text-center my-2">
        <div class="text-[18px] font-extrabold tracking-[0.06em] uppercase">Open an Account</div>
        <div class="mt-1 text-[10.5px] tracking-[0.14em] uppercase text-[#8a8878]">
          Free · No card · Track from day one
        </div>
      </div>

      <hr class="receipt-divider"/>

      <!-- First / last name (two columns) -->
      <div class="grid grid-cols-2 gap-3 mb-1">
        <FormField v-slot="$field" name="first_name">
          <div>
            <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#8a8878] mb-1 whitespace-nowrap">
              <span>1× First</span>
              <span class="text-[#c7c5b8]">OPT</span>
            </label>
            <input
              v-bind="$field"
              type="text"
              class="receipt-input"
              placeholder="Alex"
              autocomplete="given-name"
            />
          </div>
        </FormField>

        <FormField v-slot="$field" name="last_name">
          <div>
            <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#8a8878] mb-1 whitespace-nowrap">
              <span>1× Last</span>
              <span class="text-[#c7c5b8]">OPT</span>
            </label>
            <input
              v-bind="$field"
              type="text"
              class="receipt-input"
              placeholder="Chen"
              autocomplete="family-name"
            />
          </div>
        </FormField>
      </div>

      <!-- Email -->
      <FormField v-slot="$field" name="email">
        <div class="mb-3 mt-1">
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#8a8878] mb-1 whitespace-nowrap">
            <span>1× Email</span>
            <span class="text-[#c7c5b8]">OPT</span>
          </label>
          <input
            v-bind="$field"
            type="email"
            class="receipt-input"
            placeholder="you@example.com"
            autocomplete="email"
          />
          <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-600 font-medium">
            {{ $field.error?.message }}
          </div>
        </div>
      </FormField>

      <!-- Username -->
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
            placeholder="choose_a_handle"
            autocomplete="username"
          />
          <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-600 font-medium">
            {{ $field.error?.message }}
          </div>
        </div>
      </FormField>

      <!-- Password -->
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
            autocomplete="new-password"
          />
          <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-600 font-medium">
            {{ $field.error?.message }}
          </div>
        </div>
      </FormField>

      <hr class="receipt-divider"/>

      <!-- Itemized -->
      <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px] whitespace-nowrap">
        <span class="text-[#6b6b66] uppercase text-[10.5px] tracking-[0.14em]">Setup fee</span>
        <span class="font-semibold text-[#1a1a1a]">$0.00</span>
      </div>
      <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px] whitespace-nowrap">
        <span class="text-[#6b6b66] uppercase text-[10.5px] tracking-[0.14em]">Monthly</span>
        <span class="font-semibold text-[#1a1a1a]">$0.00</span>
      </div>
      <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px] whitespace-nowrap">
        <span class="text-[#6b6b66] uppercase text-[10.5px] tracking-[0.14em]">Trial period</span>
        <span class="font-semibold text-[#1a1a1a]">Forever</span>
      </div>

      <hr class="receipt-divider"/>

      <div class="flex justify-between items-center gap-3 py-1 text-[16px] font-extrabold tracking-[0.04em] whitespace-nowrap">
        <span>TOTAL DUE</span>
        <span>$0.00</span>
      </div>

      <div v-if="error" class="mt-2 px-3 py-2 bg-red-50 border border-red-200 text-red-700 text-[11.5px] leading-snug">
        {{ error.message }}
      </div>

      <button type="submit" class="receipt-btn mt-2" :disabled="isPending">
        <template v-if="isPending">PROCESSING…</template>
        <template v-else>Create Account <span class="ml-1">→</span></template>
      </button>

      <hr class="receipt-divider" style="margin-top: 18px;"/>

      <div class="text-center text-[11px] tracking-[0.16em] uppercase text-[#4a4a44] font-bold">
        ★ Welcome aboard ★
      </div>
      <div class="text-center text-[9.5px] tracking-[0.1em] text-[#8a8878] mt-1">
        Please retain this slip for your records
      </div>

      <div class="flex justify-center gap-[1.5px] my-3 h-[38px] items-stretch">
        <span
          v-for="(w, i) in barcodeBars"
          :key="i"
          class="block bg-[#1a1a1a]"
          :style="{ width: `${w}px` }"
        />
      </div>
      <div class="text-center text-[11px] tracking-[0.3em] text-[#1a1a1a] font-semibold">
        JAFA · NEW · {{ txnId.split('-')[1] }}
      </div>
    </Form>

    <div class="relative z-10 text-center text-[11px] tracking-[0.14em] uppercase text-zinc-400 whitespace-nowrap">
      Existing customer?
      <RouterLink to="/login" class="receipt-link ml-1">Sign in →</RouterLink>
    </div>
  </div>
</template>
