<script setup lang="ts">
import { computed } from 'vue'
import Menu from 'primevue/menu'
import { ref } from 'vue'
import { useDarkModeStore } from '@/stores/darkMode'
import { useAuthStore } from '@/stores/auth'
import { useSidebarStore } from '@/stores/sidebarToggle'
import { useLogout } from '@/modules/auth/composables/useAuth'
import { useRoute } from 'vue-router'

const darkMode = useDarkModeStore()
const authStore = useAuthStore()
const sidebarToggle = useSidebarStore()
const { mutate: logout } = useLogout()
const route = useRoute()

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
  { icon: 'pi pi-th-large', label: 'Dashboard', to: '/' },
  { icon: 'pi pi-list', label: 'Expenses', to: '/expenses' },
  { icon: 'pi pi-chart-bar', label: 'Reports', to: '/reports' },
  { icon: 'pi pi-tags', label: 'Categories', to: '/categories' },
  { icon: 'pi pi-cog', label: 'Settings', to: '/settings' },
]

function isActive(to: string) {
  if (to === '/') return route.path === '/'
  return route.path.startsWith(to)
}

function toggleMenu(event: Event) {
  menu.value.toggle(event)
}
</script>

<template>
  <nav
    class="flex flex-col h-full bg-[#0a0a0b] border-r border-[#26262c] transition-all duration-200 overflow-hidden shrink-0"
    :class="sidebarToggle.isExpanded ? 'w-[248px]' : 'w-[68px]'"
  >
    <!-- Brand -->
    <div class="flex items-center gap-2.5 px-4 py-4">
      <img src="/icon.png" class="w-7 h-7 shrink-0" alt=""/>
      <span
        v-if="sidebarToggle.isExpanded"
        class="text-white font-bold text-[15px] tracking-[0.12em] uppercase whitespace-nowrap"
      >JAFA</span>
    </div>

    <!-- Nav items -->
    <div class="flex flex-col gap-0.5 px-3 mt-1 flex-1">
      <RouterLink
        v-for="item in navItems"
        :key="item.label"
        :to="item.to"
        class="flex items-center gap-3 px-2.5 py-2.5 rounded-lg cursor-pointer transition-colors text-[13px] whitespace-nowrap"
        :class="[
          isActive(item.to)
            ? 'bg-[#1f1f24] text-white'
            : 'text-zinc-400 hover:text-white hover:bg-[#18181c]',
          sidebarToggle.isExpanded ? '' : 'justify-center',
        ]"
      >
        <i :class="item.icon" class="text-[15px] shrink-0" style="opacity: 0.9" />
        <span v-if="sidebarToggle.isExpanded">{{ item.label }}</span>
      </RouterLink>
    </div>

    <!-- Bottom -->
    <div class="flex flex-col gap-0.5 px-3 py-3 border-t border-[#26262c]">
      <button
        class="flex items-center gap-3 px-2.5 py-2.5 rounded-lg cursor-pointer transition-colors text-[13px] whitespace-nowrap text-zinc-400 hover:text-white hover:bg-[#18181c]"
        :class="sidebarToggle.isExpanded ? '' : 'justify-center'"
        @click="darkMode.toggle"
      >
        <i :class="darkMode.isDark ? 'pi pi-sun' : 'pi pi-moon'" class="text-[15px] shrink-0"/>
        <span v-if="sidebarToggle.isExpanded">
          {{ darkMode.isDark ? 'Light mode' : 'Dark mode' }}
        </span>
      </button>

      <button
        class="flex items-center gap-3 px-2.5 py-2.5 rounded-lg cursor-pointer transition-colors text-[13px] whitespace-nowrap text-zinc-400 hover:text-white hover:bg-[#18181c]"
        :class="sidebarToggle.isExpanded ? '' : 'justify-center'"
        @click="sidebarToggle.toggle"
      >
        <i
          :class="sidebarToggle.isExpanded ? 'pi pi-angle-double-left' : 'pi pi-angle-double-right'"
          class="text-[15px] shrink-0"
        />
        <span v-if="sidebarToggle.isExpanded">Collapse</span>
      </button>

      <button
        class="flex items-center gap-3 px-2 py-2 mt-1 rounded-lg cursor-pointer transition-colors whitespace-nowrap hover:bg-[#18181c]"
        :class="sidebarToggle.isExpanded ? '' : 'justify-center'"
        @click="toggleMenu"
      >
        <img
          v-if="authStore.currentUser?.avatar_url"
          :src="authStore.currentUser.avatar_url"
          class="w-7 h-7 rounded-full object-cover shrink-0"
        />
        <div
          v-else
          class="w-7 h-7 rounded-full flex items-center justify-center shrink-0"
          style="background: linear-gradient(135deg, #f5c518, #f97316);"
        >
          <span class="text-[11px] font-bold text-[#1a1a1a]">{{ avatarLabel }}</span>
        </div>
        <div v-if="sidebarToggle.isExpanded" class="flex flex-col items-start leading-tight overflow-hidden">
          <span class="text-[12.5px] text-white font-medium truncate max-w-[150px]">
            {{ authStore.currentUser?.username ?? 'Account' }}
          </span>
          <span class="text-[11px] text-zinc-500 truncate max-w-[150px]">
            {{ authStore.currentUser?.email ?? 'Sign in' }}
          </span>
        </div>
      </button>
      <Menu ref="menu" :model="menuItems" :popup="true" />
    </div>
  </nav>
</template>
