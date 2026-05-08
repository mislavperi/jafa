<script setup lang="ts">
import { ref, computed } from 'vue'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import { useDarkModeStore } from '@/stores/darkMode'
import { useAuthStore } from '@/stores/auth'
import { useSidebarStore } from '@/stores/sidebarToggle'
import { useLogout } from '@/modules/auth/composables/useAuth'

const darkMode = useDarkModeStore()
const authStore = useAuthStore()
const sidebarToggle = useSidebarStore()
const { mutate: logout } = useLogout()

const avatarLabel = computed(() => authStore.currentUser?.username?.charAt(0).toUpperCase() ?? 'A')

const expanded = ref()

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

function toggleMenu(event: Event) {
  menu.value.toggle(event)
}
</script>

<template>
  <nav
    class="flex flex-col py-2 border-r border-surface h-full"
    :class="[
      sidebarToggle.isExpanded ? 'w-48' : 'w-14',
      'transition-all duration-300 flex flex-col py-2 border-r border-surface h-full overflow-hidden',
    ]"
  >
    <div
      class="flex flex-col items-center"
      :class="[sidebarToggle.isExpanded ? 'items-start' : 'items-center']"
    >
      <img src="../../../public/icon.png" class="w-16" />
      <Button
        :icon="sidebarToggle.isExpanded ? 'pi pi-arrow-left' : 'pi pi-arrow-right'"
        severity="secondary"
        text
        class="border border-surface"
        @click="sidebarToggle.toggle"
      />
    </div>
    <div
      class="flex flex-col items-center gap-2"
      :class="[sidebarToggle.isExpanded ? 'items-start' : 'items-center']"
    >
      <Button
        :icon="darkMode.isDark ? 'pi pi-sun' : 'pi pi-moon'"
        severity="secondary"
        text
        rounded
        @click="darkMode.toggle"
      />
      <img
        v-if="authStore.currentUser?.avatar_url"
        :src="authStore.currentUser.avatar_url"
        class="w-9 h-9 rounded-full cursor-pointer object-cover"
        @click="toggleMenu"
      />
      <Avatar
        v-else
        :label="avatarLabel"
        size="normal"
        shape="circle"
        class="cursor-pointer"
        @click="toggleMenu"
      />
      <Menu ref="menu" :model="menuItems" :popup="true" />
    </div>
  </nav>
</template>
