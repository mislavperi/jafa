<script setup lang="ts">
import SideBar from '@/core/components/SideBar.vue'
</script>

<template class="h-screen">
  <div class="flex flex-row h-screen">
    <SideBar />
    <slot />
  </div>
</template>
