import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useDarkModeStore = defineStore('darkMode', () => {
  const hasWindow = typeof window !== 'undefined'
  const stored = hasWindow ? localStorage.getItem('darkMode') : null
  const prefersDark = hasWindow && window.matchMedia('(prefers-color-scheme: dark)').matches
  const isDark = ref(stored !== null ? stored === 'true' : prefersDark)

  function apply() {
    if (typeof document === 'undefined') return
    document.documentElement.classList.toggle('p-dark', isDark.value)
  }

  function toggle() {
    isDark.value = !isDark.value
    if (hasWindow) localStorage.setItem('darkMode', String(isDark.value))
    apply()
  }

  apply()

  return { isDark, toggle }
})
