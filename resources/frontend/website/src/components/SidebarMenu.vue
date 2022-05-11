<template>
    <div class="flex flex-col border divide-y bg-white">
        <router-link
            v-for="item in menus"
            :key="item.name"
            :to="{ name: item.name }"
            class="relative h-12 pl-6 flex items-center text-neutral-600 hover:text-neutral-900"
        >
            <span>{{ item.meta.title }}</span>
        </router-link>
    </div>
</template>
<script setup lang="ts">
    import type { PropType } from 'vue'
    import type { ISidebarMenuItem } from '~/index'
    defineProps({
        menus: {
            type: Array as PropType<ISidebarMenuItem[]>,
            default: () => [],
        },
    })
</script>
