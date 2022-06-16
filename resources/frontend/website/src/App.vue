<script setup lang="ts">
    import { storeToRefs } from 'pinia'
    import { RouterView } from 'vue-router'
    import { useAppStore } from './stores/app'
    import AcLoading from './components/AcLoading.vue'

    const store = useAppStore()
    const { isShowLoading } = storeToRefs(store)
</script>

<template>
    <AcLoading v-if="isShowLoading" />
    <RouterView />
</template>
