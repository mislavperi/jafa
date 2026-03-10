import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '@/modules/auth/models/auth'

export const useAuthStore = defineStore('auth', () => {
  const currentUser = ref<User | null>(null)
  const bootstrapped = ref(false)

  const isAuthenticated = computed(() => currentUser.value !== null)

  function setUser(user: User | null) {
    currentUser.value = user
  }

  function clearUser() {
    currentUser.value = null
    bootstrapped.value = false
  }

  function setBootstrapped() {
    bootstrapped.value = true
  }

  return { currentUser, isAuthenticated, bootstrapped, setUser, clearUser, setBootstrapped }
})
