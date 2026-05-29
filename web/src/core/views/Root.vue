<script setup lang="ts">
import SideBar from '@/core/components/SideBar.vue'
</script>

<template>
  <div class="flex flex-row h-screen bg-[var(--jafa-bg)]">
    <SideBar />
    <div class="flex-1 min-w-0 overflow-hidden bg-[var(--jafa-bg)]">
      <slot />
    </div>
  </div>
</template>
