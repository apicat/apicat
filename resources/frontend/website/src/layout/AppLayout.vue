<template>
    <AppHeader :is-show-nav="isShowNav" />
    <router-view />
</template>
<script setup lang="ts">
    import AppHeader from '../components/AppHeader.vue'
    import { ref, provide } from 'vue'

    const isShowNav = ref(true)

    const switchNavState = (isShow: boolean) => {
        isShowNav.value = isShow
    }

    provide('switchNavState', switchNavState)
</script>
