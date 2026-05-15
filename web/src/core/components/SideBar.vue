<script setup lang="ts">
import { computed } from 'vue'
import Menu from 'primevue/menu'
import { ref } from 'vue'
import { useDarkModeStore } from '@/stores/darkMode'
import { useAuthStore } from '@/stores/auth'
import { useSidebarStore } from '@/stores/sidebarToggle'
import { useLogout } from '@/modules/auth/composables/useAuth'

const darkMode = useDarkModeStore()
const authStore = useAuthStore()
const sidebarToggle = useSidebarStore()
const { mutate: logout } = useLogout()

const avatarLabel = computed(() => authStore.currentUser?.username?.charAt(0).toUpperCase() ?? 'A')

const menu = ref()
const menuItems = ref([
  {
    label: 'Edit Profile',
    icon: 'pi pi-user-edit',
    command: () => {},
  },
  {
    label: 'Logout',
    icon: 'pi pi-sign-out',
    command: () => logout(),
  },
])

const navItems = [
  { icon: 'pi pi-home', label: 'Dashboard', to: '/' },
  { icon: 'pi pi-wallet', label: 'Expenses', to: '/expenses' },
  { icon: 'pi pi-chart-bar', label: 'Reports', to: '/reports' },
  { icon: 'pi pi-tags', label: 'Categories', to: '/categories' },
  { icon: 'pi pi-cog', label: 'Settings', to: '/settings' },
]

function toggleMenu(event: Event) {
  menu.value.toggle(event)
}
</script>

<template>
  <nav
    class="flex flex-col h-full bg-zinc-100 dark:bg-[#252830] border-r border-zinc-200 dark:border-white/5 transition-all duration-300 overflow-hidden shrink-0"
    :class="sidebarToggle.isExpanded ? 'w-52' : 'w-14'"
  >
    <!-- Logo + title -->
    <div class="flex items-center gap-2.5 px-3 py-3 border-b border-zinc-200 dark:border-white/5">
      <img src="../../../public/icon.png" class="w-7 h-7 shrink-0" />
      <span v-if="sidebarToggle.isExpanded" class="text-zinc-800 dark:text-white font-semibold text-sm tracking-wide truncate">JAFA</span>
    </div>

    <!-- Nav items -->
    <div class="flex flex-col gap-0.5 p-2 mt-1 flex-1">
      <div
        v-for="item in navItems"
        :key="item.label"
        class="flex items-center gap-2.5 px-2.5 py-2 rounded-lg cursor-pointer transition-colors text-zinc-500 hover:text-zinc-900 hover:bg-zinc-200 dark:text-white/50 dark:hover:text-white dark:hover:bg-white/5"
        :class="sidebarToggle.isExpanded ? '' : 'justify-center'"
      >
        <i :class="item.icon" class="text-sm shrink-0" />
        <span v-if="sidebarToggle.isExpanded" class="text-sm">{{ item.label }}</span>
      </div>
    </div>

    <!-- Bottom controls -->
    <div class="flex flex-col gap-0.5 p-2 border-t border-zinc-200 dark:border-white/5">
      <div
        class="flex items-center gap-2.5 px-2.5 py-2 rounded-lg cursor-pointer transition-colors text-zinc-500 hover:text-zinc-900 hover:bg-zinc-200 dark:text-white/50 dark:hover:text-white dark:hover:bg-white/5"
        :class="sidebarToggle.isExpanded ? '' : 'justify-center'"
        @click="darkMode.toggle"
      >
        <i :class="darkMode.isDark ? 'pi pi-sun' : 'pi pi-moon'" class="text-sm shrink-0" />
        <span v-if="sidebarToggle.isExpanded" class="text-sm">{{ darkMode.isDark ? 'Light Mode' : 'Dark Mode' }}</span>
      </div>

      <div
        class="flex items-center gap-2.5 px-2.5 py-2 rounded-lg cursor-pointer transition-colors text-zinc-500 hover:text-zinc-900 hover:bg-zinc-200 dark:text-white/50 dark:hover:text-white dark:hover:bg-white/5"
        :class="sidebarToggle.isExpanded ? '' : 'justify-center'"
        @click="sidebarToggle.toggle"
      >
        <i :class="sidebarToggle.isExpanded ? 'pi pi-arrow-left' : 'pi pi-arrow-right'" class="text-sm shrink-0" />
        <span v-if="sidebarToggle.isExpanded" class="text-sm">Collapse</span>
      </div>

      <div
        class="flex items-center gap-2.5 px-2.5 py-2 rounded-lg cursor-pointer transition-colors text-zinc-500 hover:text-zinc-900 hover:bg-zinc-200 dark:text-white/50 dark:hover:text-white dark:hover:bg-white/5"
        :class="sidebarToggle.isExpanded ? '' : 'justify-center'"
        @click="toggleMenu"
      >
        <img
          v-if="authStore.currentUser?.avatar_url"
          :src="authStore.currentUser.avatar_url"
          class="w-6 h-6 rounded-full object-cover shrink-0"
        />
        <div
          v-else
          class="w-6 h-6 rounded-full bg-amber-400 flex items-center justify-center shrink-0"
        >
          <span class="text-xs font-bold text-black">{{ avatarLabel }}</span>
        </div>
        <span v-if="sidebarToggle.isExpanded" class="text-sm truncate text-zinc-700 dark:text-white/70">
          {{ authStore.currentUser?.username ?? 'Account' }}
        </span>
      </div>
      <Menu ref="menu" :model="menuItems" :popup="true" />
    </div>
  </nav>
</template>
