<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { Menu } from './typings'

defineProps<{
  menu: Menu
}>()
</script>

<template>
  <div class="menu">
    <i v-if="menu.icon" class="mr-1 ac-iconfont" :class="menu.icon" />
    <el-icon v-if="menu.elIcon" class="mr-1" :size="menu.size">
      <component :is="menu.elIcon" />
    </el-icon>
    <Icon v-if="menu.iconify" class="mr-1" :icon="menu.iconify" :width="menu.size" />
    <img v-if="menu.image" class="mr-1" :src="menu.image">
    {{ menu.text }}
  </div>
</template>

<style scoped lang="scss">
.menu{
    @apply flex flex-y-center;

    > * {
        @apply cursor-pointer;
        font-size: 14px;
        width: 14px;
    }
}
</style>
