import { ref } from 'vue'
import { defineStore } from 'pinia'
import { ONBOARDING_STORAGE_KEY } from '@/core/constants/storage'

export interface TourStep {
  id: string
  /** Matches a [data-tour="..."] attribute; omitted for centered steps. */
  target?: string
  title: string
  text: string
  placement: 'top' | 'bottom' | 'left' | 'right' | 'center'
}

// The tour runs on the dashboard, so every target must exist there
// (including the sidebar, which is part of the Root layout).
export const TOUR_STEPS: TourStep[] = [
  {
    id: 'welcome',
    title: 'Welcome to JAFA 👋',
    text: 'Just Another Finance App — a simple way to track where your money goes. Here is a quick tour of the essentials; it takes less than a minute.',
    placement: 'center',
  },
  {
    id: 'nav',
    target: 'nav',
    title: 'Find your way around',
    text: 'Everything lives in the sidebar: your Dashboard, the full Expenses list, Reports, Categories and Settings.',
    placement: 'right',
  },
  {
    id: 'add-expense',
    target: 'add-expense',
    title: 'Add your first expense',
    text: 'Click here to record an expense — name, amount, tags, and an optional recurring schedule for things like rent or subscriptions.',
    placement: 'bottom',
  },
  {
    id: 'scan-receipt',
    target: 'scan-receipt',
    title: 'Or scan a receipt',
    text: 'In a hurry? Snap or upload a receipt and JAFA will read the details and fill in the expense for you.',
    placement: 'bottom',
  },
  {
    id: 'stat-cards',
    target: 'stat-cards',
    title: 'Your month at a glance',
    text: 'These cards track what you have spent this month, your budget, and what is left to spend. Set a monthly budget in Settings to unlock all three.',
    placement: 'bottom',
  },
  {
    id: 'recent-expenses',
    target: 'recent-expenses',
    title: 'Recent activity',
    text: 'Your latest expenses show up here. Hover a row to edit or delete it, or jump to the full list with “View all”.',
    placement: 'right',
  },
  {
    id: 'breakdown',
    target: 'breakdown',
    title: 'See where it goes',
    text: 'This chart breaks down the current month by expense, so your biggest spends are easy to spot.',
    placement: 'left',
  },
  {
    id: 'upcoming-bills',
    target: 'upcoming-bills',
    title: 'Never miss a bill',
    text: 'Expenses with a recurring schedule appear here, sorted by due date, with a total for the next 7 days.',
    placement: 'left',
  },
  {
    id: 'settings',
    target: 'nav-settings',
    title: 'Make it yours',
    text: 'Head to Settings to pick a currency, set your monthly budget, and tweak the theme. You can replay this tour from there anytime.',
    placement: 'right',
  },
]

export const useOnboardingStore = defineStore('onboarding', () => {
  const stored =
    typeof localStorage !== 'undefined' ? localStorage.getItem(ONBOARDING_STORAGE_KEY) : null
  const completed = ref(stored === 'true')
  const active = ref(false)
  const stepIndex = ref(0)

  function start() {
    stepIndex.value = 0
    active.value = true
  }

  /** Starts the tour once for users who have never seen (or skipped) it. */
  function maybeStart() {
    if (!completed.value && !active.value) start()
  }

  function next() {
    if (stepIndex.value < TOUR_STEPS.length - 1) stepIndex.value++
    else finish()
  }

  function back() {
    if (stepIndex.value > 0) stepIndex.value--
  }

  function finish() {
    active.value = false
    completed.value = true
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(ONBOARDING_STORAGE_KEY, 'true')
    }
  }

  return { completed, active, stepIndex, start, maybeStart, next, back, finish }
})
