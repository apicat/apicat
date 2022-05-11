<template>
    <SidebarLayout :menus="menus" />
</template>
<script setup lang="ts">
    import SidebarLayout from '@/layout/SidebarLayout.vue'
    import { ref } from 'vue'
    import { UserRoutes } from '@/router/user.router'

    const menus = ref(UserRoutes)
</script>
