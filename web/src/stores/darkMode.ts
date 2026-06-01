import { ref } from 'vue'
import { defineStore } from 'pinia'
import { DARK_MODE_STORAGE_KEY as STORAGE_KEY } from '@/core/constants/storage'

export const useDarkModeStore = defineStore('darkMode', () => {
  const hasStorage = typeof localStorage !== 'undefined'
  const stored = hasStorage ? localStorage.getItem(STORAGE_KEY) : null
  const prefersDark =
    typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches
  const isDark = ref(stored !== null ? stored === 'true' : prefersDark)

  function apply() {
    if (typeof document === 'undefined') return
    document.documentElement.classList.toggle('p-dark', isDark.value)
  }

  function setDark(v: boolean) {
    isDark.value = v
    if (hasStorage) localStorage.setItem(STORAGE_KEY, String(v))
    apply()
  }

  function toggle() {
    setDark(!isDark.value)
  }

  apply()

  return { isDark, toggle, setDark }
})
