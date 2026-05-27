<script setup lang="ts">
import SideBar from '@/core/components/SideBar.vue'
</script>

<template>
  <div class="flex flex-row h-screen bg-[#0a0a0b]">
    <SideBar />
    <div class="flex-1 min-w-0 overflow-hidden bg-[#0a0a0b]">
      <slot />
    </div>
  </div>
</template>
