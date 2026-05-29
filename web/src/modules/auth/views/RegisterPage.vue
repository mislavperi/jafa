<script setup lang="ts">
import { ref } from 'vue'
import { useRegister } from '../composables/useAuth'
import { RECEIPT_INPUT_PT } from '../composables/useReceiptDecor'
import AuthReceiptLayout from '../components/AuthReceiptLayout.vue'
import { Form, FormField } from '@primevue/forms'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Message from 'primevue/message'

const { mutate: register, isPending, error } = useRegister()
const showPassword = ref(false)

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
</script>

<template>
  <AuthReceiptLayout
    tagline="★ NEW MEMBER APPLICATION ★"
    title="Open an Account"
    subtitle="Free · No card · Track from day one"
    thank-you="★ Welcome aboard ★"
    thank-you-sub="Please retain this slip for your records"
    barcode-label="NEW"
    currency="€"
    :barcode-seed="31"
    :float-seed="41"
    :ticker-seed="77"
  >
    <Form :resolver="resolver" class="flex flex-col gap-3" @submit="handleSubmit">
      <!-- First / last name -->
      <div class="grid grid-cols-2 gap-3">
        <FormField v-slot="$field" name="first_name">
          <div>
            <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#3d3a30] mb-1 whitespace-nowrap">
              <span>1× First</span>
              <span class="text-[#1a1a1a] font-bold">OPTIONAL</span>
            </label>
            <InputText
              v-bind="$field"
              type="text"
              placeholder="Alex"
              autocomplete="given-name"
              class="w-full"
              :pt="RECEIPT_INPUT_PT"
              unstyled
            />
          </div>
        </FormField>

        <FormField v-slot="$field" name="last_name">
          <div>
            <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#3d3a30] mb-1 whitespace-nowrap">
              <span>1× Last</span>
              <span class="text-[#1a1a1a] font-bold">OPTIONAL</span>
            </label>
            <InputText
              v-bind="$field"
              type="text"
              placeholder="Chen"
              autocomplete="family-name"
              class="w-full"
              :pt="RECEIPT_INPUT_PT"
              unstyled
            />
          </div>
        </FormField>
      </div>

      <FormField v-slot="$field" name="email">
        <div>
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#3d3a30] mb-1 whitespace-nowrap">
            <span>1× Email</span>
            <span class="text-[#1a1a1a] font-bold">OPTIONAL</span>
          </label>
          <InputText
            v-bind="$field"
            type="email"
            placeholder="you@example.com"
            autocomplete="email"
            class="w-full"
            :pt="RECEIPT_INPUT_PT"
            unstyled
          />
          <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-700 font-medium">
            {{ $field.error?.message }}
          </div>
        </div>
      </FormField>

      <FormField v-slot="$field" name="username">
        <div>
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase text-[#3d3a30] mb-1 whitespace-nowrap">
            <span>1× Username</span>
            <span class="text-[#1a1a1a] font-bold">REQUIRED</span>
          </label>
          <InputText
            v-bind="$field"
            type="text"
            placeholder="choose_a_handle"
            autocomplete="username"
            class="w-full"
            :pt="RECEIPT_INPUT_PT"
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
              autocomplete="new-password"
              class="w-full !pr-8"
              :pt="RECEIPT_INPUT_PT"
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

      <!-- Itemized -->
      <div class="flex flex-col">
        <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
          <span class="text-[#2d2a22] uppercase text-[10.5px] tracking-[0.14em]">Setup fee</span>
          <span class="font-semibold">€0.00</span>
        </div>
        <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
          <span class="text-[#2d2a22] uppercase text-[10.5px] tracking-[0.14em]">Monthly</span>
          <span class="font-semibold">€0.00</span>
        </div>
        <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
          <span class="text-[#2d2a22] uppercase text-[10.5px] tracking-[0.14em]">Trial period</span>
          <span class="font-semibold">Forever</span>
        </div>
      </div>

      <hr class="border-0 border-t border-dashed border-[#a8a692] my-1" />

      <div class="flex justify-between items-center gap-3 py-1 text-[16px] font-extrabold tracking-[0.04em]">
        <span>TOTAL DUE</span>
        <span>€0.00</span>
      </div>

      <Message v-if="error" severity="error" :closable="false" class="!mt-1">{{ error.message }}</Message>

      <Button
        type="submit"
        :loading="isPending"
        unstyled
        class="w-full flex items-center justify-center gap-1.5 bg-[#1a1a1a] text-[#f5f1e6] py-3 px-4 font-bold text-[13px] tracking-[0.18em] uppercase hover:opacity-85 disabled:opacity-50 transition mt-1"
      >
        <template v-if="isPending">PROCESSING…</template>
        <template v-else>Create Account <span>→</span></template>
      </Button>
    </Form>

    <template #footer>
      <div class="relative z-10 text-center text-[11px] tracking-[0.14em] uppercase text-zinc-400">
        Existing customer?
        <RouterLink to="/login" class="text-[#f5c518] hover:text-[#f97316] font-semibold underline underline-offset-[3px] tracking-[0.04em] uppercase ml-1">
          Sign in →
        </RouterLink>
      </div>
    </template>
  </AuthReceiptLayout>
</template>
