import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import { login, logout, register } from '../api/auth'
import { useAuthStore } from '@/stores/auth'
import type { LoginRequest, RegisterRequest } from '../models/auth'

export function useLogin() {
  const authStore = useAuthStore()
  const router = useRouter()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: LoginRequest) => login(payload),
    onSuccess: (user) => {
      authStore.setUser(user)
      queryClient.setQueryData(['auth', 'me'], user)
      const redirect = router.currentRoute.value.query.redirect
      router.push(typeof redirect === 'string' ? redirect : '/')
    },
  })
}

export function useLogout() {
  const authStore = useAuthStore()
  const router = useRouter()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: logout,
    onSuccess: () => {
      authStore.clearUser()
      queryClient.removeQueries({ queryKey: ['auth'] })
      router.push('/login')
    },
  })
}

export function useRegister() {
  const authStore = useAuthStore()
  const router = useRouter()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: RegisterRequest) => register(payload),
    onSuccess: (user) => {
      authStore.setUser(user)
      queryClient.setQueryData(['auth', 'me'], user)
      const redirect = router.currentRoute.value.query.redirect
      router.push(typeof redirect === 'string' ? redirect : '/')
    },
  })
}
