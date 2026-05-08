import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useSidebarStore = defineStore('sidebarToggle', () => {
  const stored = localStorage.getItem('sidebarToggle')
  const isExpanded = ref(stored !== null ? stored === 'true' : stored === 'false')

  function toggle() {
    isExpanded.value = !isExpanded.value
    localStorage.setItem('sidebarToggle', String(isExpanded.value))
  }
  return { isExpanded, toggle }
})
