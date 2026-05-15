<script setup lang="ts">
import SideBar from '@/core/components/SideBar.vue'
</script>

<template>
  <div class="flex flex-row h-screen">
    <SideBar />
    <div class="flex-1 min-w-0 overflow-hidden bg-zinc-100 dark:bg-zinc-900">
      <slot />
    </div>
  </div>
</template>
