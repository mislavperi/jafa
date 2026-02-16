import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useDarkModeStore = defineStore('darkMode', () => {
  const stored = localStorage.getItem('darkMode')
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const isDark = ref(stored !== null ? stored === 'true' : prefersDark)

  function apply() {
    document.documentElement.classList.toggle('p-dark', isDark.value)
  }

  function toggle() {
    isDark.value = !isDark.value
    localStorage.setItem('darkMode', String(isDark.value))
    apply()
  }

  apply()

  return { isDark, toggle }
})
