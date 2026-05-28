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

interface PrefsPayload {
  accentId: string
  fontSize: FontSize
  darkMode: boolean
}

import { apiFetch } from '@/core/api'
import { isAuthRequiredError } from '@/core/auth-error'

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

  async function load(isDark: boolean) {
    const prefs = await fetchPrefs()
    if (prefs) {
      suppressSave = true
      accentId.value = prefs.accentId
      fontSize.value = prefs.fontSize
      suppressSave = false
    } else {
      // No prefs yet — seed defaults from current dark-mode state
      suppressSave = true
      void isDark
      suppressSave = false
    }
    loaded.value = true
    apply()
  }

  function persist(darkMode: boolean) {
    if (suppressSave || !loaded.value) return
    void savePrefs({ accentId: accentId.value, fontSize: fontSize.value, darkMode })
  }

  watch([accentId, fontSize], () => {
    apply()
  })

  apply()

  return {
    accentId,
    fontSize,
    loaded,
    setAccent,
    setFontSize,
    currentAccent,
    load,
    persist,
  }
})
