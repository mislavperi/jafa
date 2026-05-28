import { ref } from 'vue'
import { defineStore } from 'pinia'

const STORAGE_KEY = 'darkMode'

export const useDarkModeStore = defineStore('darkMode', () => {
  const stored = localStorage.getItem(STORAGE_KEY)
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const isDark = ref(stored !== null ? stored === 'true' : prefersDark)

  function apply() {
    if (typeof document === 'undefined') return
    document.documentElement.classList.toggle('p-dark', isDark.value)
  }

  function setDark(v: boolean) {
    isDark.value = v
    localStorage.setItem(STORAGE_KEY, String(v))
    apply()
  }

  function toggle() {
    setDark(!isDark.value)
  }

  apply()

  return { isDark, toggle, setDark }
})
