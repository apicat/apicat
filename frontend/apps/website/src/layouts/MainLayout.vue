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
<script setup>
import AsideNav from '@/views/component/AsideNav.vue'
</script>
