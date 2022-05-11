<template>
    <PreviewHeader v-if="isShowHeader" :is-show-search="isShowSearch" />
    <router-view />
    <ac-backtop :bottom="100" :right="100" />
</template>
<script setup lang="ts">
    import { ref, provide } from 'vue'
    import PreviewHeader from '../views/preview/components/PreviewHeader.vue'
    const isShowSearch = ref(true)
    const isShowHeader = ref(true)

    const showSearchInput = (isShow: boolean) => {
        isShowSearch.value = isShow
    }

    const showHeader = (isShow: boolean) => {
        isShowHeader.value = isShow
    }

    provide('showSearchInput', showSearchInput)
    provide('showHeader', showHeader)
</script>
