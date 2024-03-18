<script setup>
import AsideNav from '@/views/component/AsideNav/index.vue'
</script>

<template>
  <div class="flex w-full h-screen">
    <div class="h-full">
      <AsideNav />
    </div>
    <div class="flex-1 overflow-auto">
      <router-view />
    </div>
  </div>
</template>
