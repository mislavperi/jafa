<script setup lang="ts">
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
</script>

<template>
  <div class="flex items-center justify-center min-h-screen">
    <div class="w-full max-w-sm p-6 border border-surface-200 dark:border-surface-700 rounded-xl shadow-sm">
      <h1 class="text-2xl font-bold mb-6">Sign in to JAFA</h1>

      <Form :resolver="resolver" @submit="handleSubmit" class="flex flex-col gap-4">
        <FormField v-slot="$field" name="username">
          <div class="flex flex-col gap-1">
            <label class="text-sm font-medium">Username</label>
            <InputText
              v-bind="$field"
              placeholder="Your username"
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
              placeholder="Your password"
              :feedback="false"
              autocomplete="current-password"
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

        <Button type="submit" label="Sign in" :loading="isPending" class="w-full" />
      </Form>

      <p class="mt-4 text-sm text-center text-surface-500">
        Don't have an account?
        <RouterLink to="/register" class="text-primary-500 hover:underline font-medium">Register</RouterLink>
      </p>
    </div>
  </div>
</template>
