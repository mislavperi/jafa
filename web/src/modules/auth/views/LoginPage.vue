<script setup lang="ts">
import { ref } from 'vue'
import { useLogin } from '../composables/useAuth'
import { RECEIPT_INPUT_PT } from '../composables/useReceiptDecor'
import AuthReceiptLayout from '../components/AuthReceiptLayout.vue'
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
</script>

<template>
  <AuthReceiptLayout
    tagline="★ EXPENSE TRACKING CO. ★"
    title="Customer Sign-In"
    subtitle="Please enter your credentials"
    thank-you="★ Thank you for budgeting ★"
    thank-you-sub="Every dollar tracked is a dollar earned"
    barcode-label="0042"
    currency="$"
    :barcode-seed="7"
    :float-seed="13"
    :ticker-seed="99"
  >
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
              autocomplete="current-password"
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

    <template #footer>
      <div class="relative z-10 text-center text-[11px] tracking-[0.14em] uppercase text-zinc-400">
        New customer?
        <RouterLink to="/register" class="text-[#f5c518] hover:text-[#f97316] font-semibold underline underline-offset-[3px] tracking-[0.04em] uppercase ml-1">
          Open an account →
        </RouterLink>
      </div>
    </template>
  </AuthReceiptLayout>
</template>
