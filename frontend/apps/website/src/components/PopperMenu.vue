<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { Menu } from './typings'
import { useNamespace } from '@/hooks/useNamespace'

export interface Props {
  menus: Menu[] | Record<string, any>
  rowKey?: string
  activeMenuKey?: any
  size?: string
  center?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  rowKey: 'text',
  activeMenuKey: null,
  menus: () => [],
  size: '',
  center: false,
})

const emits = defineEmits(['menuClick'])
const ns = useNamespace('popper-menu')
function onMenuClick(menu: Menu, keyOrIndex: any) {
  menu.onClick && menu.onClick()
  emits('menuClick', menu, keyOrIndex)
}
function isActive(menu: any) {
  return menu[props.rowKey] === props.activeMenuKey
}
const slots = useSlots()
</script>

<template>
  <ul :class="[ns.b(), ns.m(size)]">
    <li
      v-for="(menu, keyOrIndex) in menus"
      :key="menu[rowKey]"
      :class="[
        ns.e('item'),
        {
          'border-t': menu.divided,
          'text-blue-primary': slots.suffix ? false : isActive(menu),
        },
      ]"
      @click="onMenuClick(menu, keyOrIndex)"
    >
      <div v-if="menu.content">
        <component :is="menu.content" />
      </div>
      <div v-else class="row">
        <div
          class="left"
          :class="{
            'text-center': center,
            'items-center': center,
            'content-center': center,
            'justify-center': center,
          }"
          style="overflow: hidden"
          :title="menu.text"
        >
          <slot name="prefix">
            <i v-if="menu.icon" class="mr-1 ac-iconfont" :class="[ns.e('icon'), menu.icon]" />
            <el-icon v-if="menu.elIcon" class="mr-1" :size="menu.size">
              <component :is="menu.elIcon" />
            </el-icon>
            <Icon v-if="menu.iconify" class="mr-1" :icon="menu.iconify" :width="menu.size" />
            <img v-if="menu.image" class="mr-1" :class="ns.e('icon')" :src="menu.image">
          </slot>
          <p class="truncate line-height-16px">
            {{ menu.refText ? menu.refText.value : menu.text }}
          </p>
        </div>
        <div v-if="isActive(menu)" class="right">
          <slot name="suffix" />
        </div>
      </div>
    </li>
  </ul>
</template>

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

.row {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.left,
.right {
  display: flex;
  align-items: center;
}

.left {
  flex-grow: 1;
}

.right {
  justify-content: flex-end;
}
</style>
