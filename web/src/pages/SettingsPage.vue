<script setup lang="ts">
import { ref, computed } from 'vue'
import Root from '@/core/views/Root.vue'
import AppPageHeader from '@/core/components/AppPageHeader.vue'
import ToggleSwitch from 'primevue/toggleswitch'
import Select from 'primevue/select'
import Button from 'primevue/button'
import { useAuthStore } from '@/stores/auth'

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

const currency = ref('USD')
const weekStart = ref('Monday')

const currencyOptions = ['USD', 'EUR', 'GBP', 'CAD', 'AUD', 'JPY']
const weekOptions = ['Monday', 'Sunday', 'Saturday']

const notifWeeklySummary = ref(true)
const notifBudgetAlerts = ref(true)
const notifProductUpdates = ref(false)
</script>

<template>
  <Root>
    <div class="flex flex-col gap-5 h-full min-w-0 p-8 overflow-auto">
      <AppPageHeader title="Settings" subtitle="Manage your account and preferences" />

      <div class="flex flex-col gap-4" style="max-width: 720px">

        <!-- Profile -->
        <div class="bg-[#131316] border border-[#26262c] rounded-[14px] p-5 flex flex-col gap-5">
          <p class="text-[11px] font-semibold uppercase tracking-[0.08em] text-zinc-400">Profile</p>
          <div class="flex items-center gap-4">
            <div
              class="w-14 h-14 rounded-full flex items-center justify-center shrink-0 text-xl font-bold text-[#131316]"
              style="background: linear-gradient(135deg, #f5c518, #f97316)"
            >
              {{ initial }}
            </div>
            <div>
              <p class="text-white font-semibold text-base">{{ displayName }}</p>
              <p class="text-zinc-400 text-sm">{{ user?.email ?? '—' }}</p>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-zinc-400 uppercase tracking-wider">Username</label>
              <div class="bg-[#1f1f24] border border-[#26262c] rounded-[8px] px-3 py-2 text-white text-sm">
                {{ user?.username ?? '—' }}
              </div>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-zinc-400 uppercase tracking-wider">Email</label>
              <div class="bg-[#1f1f24] border border-[#26262c] rounded-[8px] px-3 py-2 text-white text-sm">
                {{ user?.email ?? '—' }}
              </div>
            </div>
          </div>
        </div>

        <!-- Preferences -->
        <div class="bg-[#131316] border border-[#26262c] rounded-[14px] p-5 flex flex-col gap-5">
          <p class="text-[11px] font-semibold uppercase tracking-[0.08em] text-zinc-400">Preferences</p>
          <div class="grid grid-cols-2 gap-4">
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-zinc-400 uppercase tracking-wider">Currency</label>
              <Select v-model="currency" :options="currencyOptions" class="w-full" />
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-zinc-400 uppercase tracking-wider">Week starts on</label>
              <Select v-model="weekStart" :options="weekOptions" class="w-full" />
            </div>
          </div>
        </div>

        <!-- Notifications -->
        <div class="bg-[#131316] border border-[#26262c] rounded-[14px] p-5 flex flex-col gap-4">
          <p class="text-[11px] font-semibold uppercase tracking-[0.08em] text-zinc-400">Notifications</p>
          <div class="flex flex-col divide-y divide-[#26262c]">
            <div class="flex items-center justify-between py-3 first:pt-0 last:pb-0">
              <div>
                <p class="text-white text-sm font-medium">Weekly summary</p>
                <p class="text-zinc-500 text-xs mt-0.5">Get a weekly overview of your spending</p>
              </div>
              <ToggleSwitch v-model="notifWeeklySummary" />
            </div>
            <div class="flex items-center justify-between py-3">
              <div>
                <p class="text-white text-sm font-medium">Budget alerts</p>
                <p class="text-zinc-500 text-xs mt-0.5">Alert when you approach your budget limit</p>
              </div>
              <ToggleSwitch v-model="notifBudgetAlerts" />
            </div>
            <div class="flex items-center justify-between py-3 last:pb-0">
              <div>
                <p class="text-white text-sm font-medium">Product updates</p>
                <p class="text-zinc-500 text-xs mt-0.5">Occasional news about new features</p>
              </div>
              <ToggleSwitch v-model="notifProductUpdates" />
            </div>
          </div>
        </div>

        <!-- Danger zone -->
        <div class="bg-[#131316] border border-red-900/40 rounded-[14px] p-5 flex flex-col gap-4">
          <p class="text-[11px] font-semibold uppercase tracking-[0.08em] text-red-400">Danger Zone</p>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-white text-sm font-medium">Delete account</p>
              <p class="text-zinc-500 text-xs mt-0.5">Permanently delete your account and all data</p>
            </div>
            <Button label="Delete Account" severity="danger" size="small" outlined />
          </div>
        </div>

      </div>
    </div>
  </Root>
</template>
