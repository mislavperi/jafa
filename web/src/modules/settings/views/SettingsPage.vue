<script setup lang="ts">
import { ref, computed } from 'vue'
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import ToggleSwitch from 'primevue/toggleswitch'
import Select from 'primevue/select'
import Button from 'primevue/button'
import { useAuthStore } from '@/stores/auth'
import { useDarkModeStore } from '@/stores/darkMode'
import { useThemeStore, ACCENTS, type FontSize } from '@/stores/theme'

const authStore = useAuthStore()

const user = computed(() => authStore.currentUser)

const initial = computed(() => {
  if (user.value?.first_name) return user.value.first_name[0].toUpperCase()
  if (user.value?.username) return user.value.username[0].toUpperCase()
  return '?'
})

const displayName = computed(() => {
  if (user.value?.first_name && user.value?.last_name) {
    return user.value.first_name + ' ' + user.value.last_name
  }
  return user.value?.username ?? ''
})

const currency = ref('EUR')
const weekStart = ref('Monday')

const currencyOptions = ['EUR']
const weekOptions = ['Monday', 'Sunday', 'Saturday']

const notifWeeklySummary = ref(true)
const notifBudgetAlerts = ref(true)
const notifProductUpdates = ref(false)

const darkMode = useDarkModeStore()
const theme = useThemeStore()

function setAccent(id: string) {
  theme.setAccent(id)
  theme.persist(darkMode.isDark)
}

function setFontSize(size: FontSize) {
  theme.setFontSize(size)
  theme.persist(darkMode.isDark)
}

function toggleDark() {
  darkMode.toggle()
  theme.persist(darkMode.isDark)
}

const fontSizes: { label: string; value: FontSize }[] = [
  { label: 'Small', value: 'small' },
  { label: 'Normal', value: 'normal' },
  { label: 'Large', value: 'large' },
]
</script>

<template>
  <Root>
    <div class="flex flex-col gap-5 h-full min-w-0 p-8 overflow-auto">
      <AppPageHeader title="Settings" subtitle="Manage your account and preferences" />

      <div class="flex flex-col gap-4" style="max-width: 720px">

        <!-- Profile -->
        <div class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 flex flex-col gap-5">
          <p class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.08em] text-[var(--jafa-text-muted)]">Profile</p>
          <div class="flex items-center gap-4">
            <div
              class="w-14 h-14 rounded-full flex items-center justify-center shrink-0 text-xl font-bold text-[var(--jafa-surface)]"
              style="background: var(--jafa-accent)"
            >
              {{ initial }}
            </div>
            <div>
              <p class="text-[var(--jafa-text)] font-semibold text-base">{{ displayName }}</p>
              <p class="text-[var(--jafa-text-muted)] text-sm">{{ user?.email ?? '—' }}</p>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-[var(--jafa-text-muted)] uppercase tracking-wider">Username</label>
              <div class="bg-[var(--jafa-surface-3)] border border-[var(--jafa-border)] rounded-[8px] px-3 py-2 text-[var(--jafa-text)] text-sm">
                {{ user?.username ?? '—' }}
              </div>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-[var(--jafa-text-muted)] uppercase tracking-wider">Email</label>
              <div class="bg-[var(--jafa-surface-3)] border border-[var(--jafa-border)] rounded-[8px] px-3 py-2 text-[var(--jafa-text)] text-sm">
                {{ user?.email ?? '—' }}
              </div>
            </div>
          </div>
        </div>

        <!-- Appearance -->
        <div class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 flex flex-col gap-5">
          <p class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.08em] text-[var(--jafa-text-muted)]">Appearance</p>

          <!-- Theme mode -->
          <div class="flex items-center justify-between">
            <div>
              <p class="text-[var(--jafa-text)] text-sm font-medium">Dark mode</p>
              <p class="text-[var(--jafa-text-muted)] text-xs mt-0.5">Toggle between light and dark theme</p>
            </div>
            <ToggleSwitch :model-value="darkMode.isDark" @update:model-value="toggleDark" />
          </div>

          <!-- Accent color -->
          <div class="flex flex-col gap-2.5">
            <label class="text-xs font-medium text-[var(--jafa-text-muted)] uppercase tracking-wider">Accent color</label>
            <div class="grid grid-cols-6 gap-2">
              <button
                v-for="a in ACCENTS"
                :key="a.id"
                type="button"
                class="aspect-square rounded-lg border-2 flex items-center justify-center transition"
                :style="{
                  background: a.color,
                  borderColor: theme.accentId === a.id ? 'var(--jafa-text)' : 'transparent',
                }"
                :aria-label="a.name"
                @click="setAccent(a.id)"
              >
                <i v-if="theme.accentId === a.id" class="pi pi-check text-[calc(14px*var(--jafa-text-scale,1))]" :style="{ color: a.text }" />
              </button>
            </div>
          </div>

          <!-- Font size -->
          <div class="flex flex-col gap-2.5">
            <label class="text-xs font-medium text-[var(--jafa-text-muted)] uppercase tracking-wider">Font size</label>
            <div class="grid grid-cols-3 gap-2">
              <button
                v-for="f in fontSizes"
                :key="f.value"
                type="button"
                class="px-3 py-2.5 rounded-md border text-[calc(13px*var(--jafa-text-scale,1))] font-medium transition"
                :class="theme.fontSize === f.value
                  ? 'border-[var(--jafa-accent)] bg-[var(--jafa-accent)]/10 text-[var(--jafa-text)]'
                  : 'border-[var(--jafa-border)] bg-[var(--jafa-surface-2)] text-[var(--jafa-text-muted)] hover:text-[var(--jafa-text)]'"
                @click="setFontSize(f.value)"
              >
                <span :class="{
                  'text-[calc(12px*var(--jafa-text-scale,1))]': f.value === 'small',
                  'text-[calc(14px*var(--jafa-text-scale,1))]': f.value === 'normal',
                  'text-[calc(16px*var(--jafa-text-scale,1))]': f.value === 'large',
                }">{{ f.label }}</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Preferences -->
        <div class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 flex flex-col gap-5">
          <p class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.08em] text-[var(--jafa-text-muted)]">Preferences</p>
          <div class="grid grid-cols-2 gap-4">
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-[var(--jafa-text-muted)] uppercase tracking-wider">Currency</label>
              <Select v-model="currency" :options="currencyOptions" class="w-full" disabled />
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-[var(--jafa-text-muted)] uppercase tracking-wider">Week starts on</label>
              <Select v-model="weekStart" :options="weekOptions" class="w-full" />
            </div>
          </div>
        </div>

        <!-- Notifications -->
        <div class="bg-[var(--jafa-surface)] border border-[var(--jafa-border)] rounded-[14px] p-5 flex flex-col gap-4">
          <p class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.08em] text-[var(--jafa-text-muted)]">Notifications</p>
          <div class="flex flex-col divide-y divide-[var(--jafa-border)]">
            <div class="flex items-center justify-between py-3 first:pt-0 last:pb-0">
              <div>
                <p class="text-[var(--jafa-text)] text-sm font-medium">Weekly summary</p>
                <p class="text-[var(--jafa-text-muted)] text-xs mt-0.5">Get a weekly overview of your spending</p>
              </div>
              <ToggleSwitch v-model="notifWeeklySummary" />
            </div>
            <div class="flex items-center justify-between py-3">
              <div>
                <p class="text-[var(--jafa-text)] text-sm font-medium">Budget alerts</p>
                <p class="text-[var(--jafa-text-muted)] text-xs mt-0.5">Alert when you approach your budget limit</p>
              </div>
              <ToggleSwitch v-model="notifBudgetAlerts" />
            </div>
            <div class="flex items-center justify-between py-3 last:pb-0">
              <div>
                <p class="text-[var(--jafa-text)] text-sm font-medium">Product updates</p>
                <p class="text-[var(--jafa-text-muted)] text-xs mt-0.5">Occasional news about new features</p>
              </div>
              <ToggleSwitch v-model="notifProductUpdates" />
            </div>
          </div>
        </div>

        <!-- Danger zone -->
        <div class="bg-[var(--jafa-surface)] border border-red-900/40 rounded-[14px] p-5 flex flex-col gap-4">
          <p class="text-[calc(11px*var(--jafa-text-scale,1))] font-semibold uppercase tracking-[0.08em] text-red-400">Danger Zone</p>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-[var(--jafa-text)] text-sm font-medium">Delete account</p>
              <p class="text-[var(--jafa-text-muted)] text-xs mt-0.5">Permanently delete your account and all data</p>
            </div>
            <Button label="Delete Account" severity="danger" size="small" outlined />
          </div>
        </div>

      </div>
    </div>
  </Root>
</template>
