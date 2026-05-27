import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useSidebarStore = defineStore('sidebarToggle', () => {
  const stored = typeof localStorage !== 'undefined' ? localStorage.getItem('sidebarToggle') : null
  const isExpanded = ref(stored !== null ? stored === 'true' : true)

  function toggle() {
    isExpanded.value = !isExpanded.value
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('sidebarToggle', String(isExpanded.value))
    }
  }
  return { isExpanded, toggle }
})
