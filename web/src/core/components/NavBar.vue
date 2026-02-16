<script setup lang="ts">
import { ref } from 'vue'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import OverlayBadge from 'primevue/overlaybadge'
import Menu from 'primevue/menu'
import { useDarkModeStore } from '@/stores/darkMode'

const darkMode = useDarkModeStore()
const menu = ref()
const menuItems = ref([
  {
    label: 'Notifications',
    icon: 'pi pi-bell',
    command: () => {}
  },
  {
    label: 'Edit Profile',
    icon: 'pi pi-user-edit',
    command: () => {}
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
      <OverlayBadge value="4" severity="danger" class="inline-flex cursor-pointer" @click="toggleMenu">
        <Avatar label="A" size="large" shape="circle" />
      </OverlayBadge>
      <Menu ref="menu" :model="menuItems" :popup="true" />
    </div>
  </nav>
</template>
