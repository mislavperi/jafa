<script setup lang="ts">
import { ref, computed } from 'vue'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import { useDarkModeStore } from '@/stores/darkMode'
import { useAuthStore } from '@/stores/auth'
import { useLogout } from '@/modules/auth/composables/useAuth'

const darkMode = useDarkModeStore()
const authStore = useAuthStore()
const { mutate: logout } = useLogout()

const avatarLabel = computed(() =>
  authStore.currentUser?.username?.charAt(0).toUpperCase() ?? 'A'
)

const menu = ref()
const menuItems = ref([
  {
    label: 'Edit Profile',
    icon: 'pi pi-user-edit',
    command: () => {}
  },
  {
    label: 'Logout',
    icon: 'pi pi-sign-out',
    command: () => logout()
  }
])

function toggleMenu(event: Event) {
  menu.value.toggle(event)
}
</script>

<template>
  <nav class="flex items-center justify-between px-4 py-2 border-b border-surface sm:px-6">
    <div class="text-xl font-bold">JAFA</div>
    <div class="flex items-center gap-2">
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
      <Avatar v-else :label="avatarLabel" size="normal" shape="circle" class="cursor-pointer" @click="toggleMenu" />
      <Menu ref="menu" :model="menuItems" :popup="true" />
    </div>
  </nav>
</template>
