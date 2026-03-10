<script setup lang="ts">
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
</script>

<template>
  <div class="flex items-center justify-center min-h-screen">
    <div class="w-full max-w-sm p-6 border border-surface-200 dark:border-surface-700 rounded-xl shadow-sm">
      <h1 class="text-2xl font-bold mb-6">Create an account</h1>

      <Form :resolver="resolver" @submit="handleSubmit" class="flex flex-col gap-4">
        <div class="flex gap-3">
          <FormField v-slot="$field" name="first_name" class="flex-1">
            <div class="flex flex-col gap-1">
              <label class="text-sm font-medium">First name</label>
              <InputText
                v-bind="$field"
                placeholder="First name"
                autocomplete="given-name"
                class="w-full"
              />
            </div>
          </FormField>
          <FormField v-slot="$field" name="last_name" class="flex-1">
            <div class="flex flex-col gap-1">
              <label class="text-sm font-medium">Last name</label>
              <InputText
                v-bind="$field"
                placeholder="Last name"
                autocomplete="family-name"
                class="w-full"
              />
            </div>
          </FormField>
        </div>

        <FormField v-slot="$field" name="email">
          <div class="flex flex-col gap-1">
            <label class="text-sm font-medium">Email</label>
            <InputText
              v-bind="$field"
              type="email"
              placeholder="you@example.com"
              autocomplete="email"
              class="w-full"
              :invalid="$field.invalid"
            />
            <Message v-if="$field.invalid" severity="error" size="small" variant="simple">
              {{ $field.error?.message }}
            </Message>
          </div>
        </FormField>

        <FormField v-slot="$field" name="username">
          <div class="flex flex-col gap-1">
            <label class="text-sm font-medium">Username</label>
            <InputText
              v-bind="$field"
              placeholder="Choose a username"
              autocomplete="username"
              class="w-full"
              :invalid="$field.invalid"
            />
            <Message v-if="$field.invalid" severity="error" size="small" variant="simple">
              {{ $field.error?.message }}
            </Message>
          </div>
        </FormField>

        <FormField v-slot="$field" name="password">
          <div class="flex flex-col gap-1">
            <label class="text-sm font-medium">Password</label>
            <Password
              v-bind="$field"
              placeholder="Choose a password"
              autocomplete="new-password"
              input-class="w-full"
              class="w-full"
              :invalid="$field.invalid"
            />
            <Message v-if="$field.invalid" severity="error" size="small" variant="simple">
              {{ $field.error?.message }}
            </Message>
          </div>
        </FormField>

        <Message v-if="error" severity="error" :closable="false">{{ error.message }}</Message>

        <Button type="submit" label="Create account" :loading="isPending" class="w-full" />
      </Form>

      <p class="mt-4 text-sm text-center text-surface-500">
        Already have an account?
        <RouterLink to="/login" class="text-primary-500 hover:underline font-medium">Sign in</RouterLink>
      </p>
    </div>
  </div>
</template>
