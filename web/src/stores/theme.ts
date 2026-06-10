import { ref, watch } from 'vue'
import { defineStore } from 'pinia'

export interface AccentSwatch {
  id: string
  name: string
  color: string
  text: string
}

export const ACCENTS: AccentSwatch[] = [
  { id: 'amber',   name: 'Amber',   color: '#f5c518', text: '#1a1a1a' },
  { id: 'orange',  name: 'Orange',  color: '#f97316', text: '#1a1a1a' },
  { id: 'emerald', name: 'Emerald', color: '#10b981', text: '#ffffff' },
  { id: 'sky',     name: 'Sky',     color: '#0ea5e9', text: '#ffffff' },
  { id: 'violet',  name: 'Violet',  color: '#8b5cf6', text: '#ffffff' },
  { id: 'rose',    name: 'Rose',    color: '#f43f5e', text: '#ffffff' },
]

export type FontSize = 'small' | 'normal' | 'large'

const FONT_SCALE: Record<FontSize, string> = {
  small: '0.875',
  normal: '1',
  large: '1.15',
}

export type WeekStart = 'Monday' | 'Sunday' | 'Saturday'

interface PrefsPayload {
  accentId: string
  fontSize: FontSize
  darkMode: boolean
  currency: string
  weekStart: WeekStart
  monthlyBudget: number
  notifyWeeklySummary: boolean
  notifyBudgetAlerts: boolean
  notifyProductUpdates: boolean
}

import { apiFetch } from '@/core/api'
import { isAuthRequiredError } from '@/core/auth-error'
import { useDarkModeStore } from '@/stores/darkMode'

async function fetchPrefs(): Promise<PrefsPayload | null> {
  try {
    const r = await apiFetch('/api/preferences')
    if (!r.ok) return null
    return await r.json()
  } catch (e) {
    if (isAuthRequiredError(e)) throw e
    return null
  }
}

async function savePrefs(p: PrefsPayload): Promise<void> {
  try {
    await apiFetch('/api/preferences', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(p),
    })
  } catch (e) {
    if (isAuthRequiredError(e)) throw e
    // ignore network errors
  }
}

export const useThemeStore = defineStore('theme', () => {
  const accentId = ref('amber')
  const fontSize = ref<FontSize>('normal')
  const currency = ref('EUR')
  const weekStart = ref<WeekStart>('Monday')
  const monthlyBudget = ref(0)
  const notifyWeeklySummary = ref(true)
  const notifyBudgetAlerts = ref(true)
  const notifyProductUpdates = ref(false)
  const loaded = ref(false)
  let suppressSave = false

  function currentAccent(): AccentSwatch {
    return ACCENTS.find((a) => a.id === accentId.value) ?? ACCENTS[0]!
  }

  function apply() {
    const a = currentAccent()
    const root = document.documentElement
    root.style.setProperty('--jafa-accent', a.color)
    root.style.setProperty('--jafa-accent-text', a.text)
    root.style.setProperty('--jafa-text-scale', FONT_SCALE[fontSize.value])
  }

  function setAccent(id: string) {
    accentId.value = id
  }

  function setFontSize(size: FontSize) {
    fontSize.value = size
  }

  function setCurrency(code: string) {
    currency.value = code
  }

  function setWeekStart(day: WeekStart) {
    weekStart.value = day
  }

  function setMonthlyBudget(amount: number) {
    monthlyBudget.value = Math.max(0, amount)
  }

  function setNotification(key: 'weeklySummary' | 'budgetAlerts' | 'productUpdates', value: boolean) {
    if (key === 'weeklySummary') notifyWeeklySummary.value = value
    else if (key === 'budgetAlerts') notifyBudgetAlerts.value = value
    else notifyProductUpdates.value = value
  }

  async function load(isDark: boolean) {
    const prefs = await fetchPrefs()
    if (prefs) {
      suppressSave = true
      accentId.value = prefs.accentId
      fontSize.value = prefs.fontSize
      if (prefs.currency) currency.value = prefs.currency
      if (prefs.weekStart) weekStart.value = prefs.weekStart
      if (typeof prefs.monthlyBudget === 'number') monthlyBudget.value = prefs.monthlyBudget
      if (typeof prefs.notifyWeeklySummary === 'boolean') notifyWeeklySummary.value = prefs.notifyWeeklySummary
      if (typeof prefs.notifyBudgetAlerts === 'boolean') notifyBudgetAlerts.value = prefs.notifyBudgetAlerts
      if (typeof prefs.notifyProductUpdates === 'boolean') notifyProductUpdates.value = prefs.notifyProductUpdates
      // Restore the saved dark-mode choice so it syncs across devices, not just
      // the value cached in localStorage.
      if (typeof prefs.darkMode === 'boolean') {
        useDarkModeStore().setDark(prefs.darkMode)
      }
      suppressSave = false
    } else {
      // No prefs row yet — keep the current dark-mode state (already mirrored to
      // the backend the first time the user changes it).
      void isDark
    }
    loaded.value = true
    apply()
  }

  function persist(darkMode: boolean) {
    if (suppressSave || !loaded.value) return
    void savePrefs({
      accentId: accentId.value,
      fontSize: fontSize.value,
      darkMode,
      currency: currency.value,
      weekStart: weekStart.value,
      monthlyBudget: monthlyBudget.value,
      notifyWeeklySummary: notifyWeeklySummary.value,
      notifyBudgetAlerts: notifyBudgetAlerts.value,
      notifyProductUpdates: notifyProductUpdates.value,
    })
  }

  watch([accentId, fontSize], () => {
    apply()
  })

  apply()

  return {
    accentId,
    fontSize,
    currency,
    weekStart,
    monthlyBudget,
    notifyWeeklySummary,
    notifyBudgetAlerts,
    notifyProductUpdates,
    loaded,
    setAccent,
    setFontSize,
    setCurrency,
    setWeekStart,
    setMonthlyBudget,
    setNotification,
    currentAccent,
    load,
    persist,
  }
})
