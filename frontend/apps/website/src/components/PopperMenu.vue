<template>
  <ul :class="[ns.b(), ns.m(size)]">
    <li
      v-for="(menu, keyOrIndex) in menus"
      :key="menu[rowKey]"
      :class="[ns.e('item'), { 'border-t': menu.divided, 'text-blue-primary': menu[rowKey] === activeMenuKey }]"
      @click="onMenuClick(menu, keyOrIndex)"
    >
      <i v-if="menu.icon" class="mr-1 ac-iconfont" :class="[ns.e('icon'), menu.icon]"></i>
      <el-icon v-if="menu.elIcon" class="mr-1" :size="menu.size"><component :is="menu.elIcon" /></el-icon>
      <img v-if="menu.image" class="mr-1" :class="ns.e('icon')" :src="menu.image" />
      {{ menu.text }}
    </li>
  </ul>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks/useNamespace'
import { Menu } from './typings'

interface Props {
  menus: Menu[] | Record<string, any>
  rowKey?: string
  activeMenuKey?: any
  size?: string
}

const props = withDefaults(defineProps<Props>(), {
  rowKey: 'text',
  activeMenuKey: null,
  menus: () => [],
  size: '',
})

const emits = defineEmits(['menu-click'])
const ns = useNamespace('popper-menu')
const onMenuClick = (menu: Menu, keyOrIndex: any) => {
  menu.onClick && menu.onClick()
  emits('menu-click', menu, keyOrIndex)
}
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(popper-menu) {
  padding: 5px 0;
  min-width: auto;

  @include m(thin) {
    padding: 4px 0;
    @include e(item) {
      padding: 5px 16px;
    }
  }
  @include m(small) {
    @include e(item) {
      padding: 8px 16px;
    }
  }

  @include m(large) {
    padding: 10px 0;
    @include e(item) {
      padding: 10px 34px;
    }
  }

  @include m(bg-white) {
    @include e(item) {
      @apply hover:bg-white border-b border-gray-6;

      &:last-child {
        border-bottom: none;
      }
    }
  }

  @include e(item) {
    @apply flex items-center py-2.5 px-6 cursor-pointer hover:bg-neutral-100 truncate;
    &:hover {
      color: var(--title-color);
    }
  }

  @include e(icon) {
    font-size: 16px;
    width: 16px;
    height: 16px;
    margin-right: 4px;
    line-height: 1;
  }
}
</style>
