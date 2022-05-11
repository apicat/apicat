<template>
    <MainHeader />
    <main class="viewport px-8 py-4">
        <router-view />
    </main>
</template>
<script setup lang="ts">
    import MainHeader from './MainHeader.vue'
</script>
