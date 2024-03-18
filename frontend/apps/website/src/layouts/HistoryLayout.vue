<script setup lang="ts">
import { useHistoryLayoutProvide } from './useHistoryLayoutContext'
import { useNamespace } from '@/hooks/useNamespace'

const props = defineProps<{
  goBack: () => void
}>()

const ns = useNamespace('doc-layout')
const historyInfo = useNamespace('history-info')
const historyLayoutContext = useHistoryLayoutProvide()
historyLayoutContext.goBack = props.goBack
</script>

<template>
  <main :class="ns.b()">
    <div :class="historyInfo.b()">
      <div :class="historyInfo.e('img')">
        <a href="javascript:void(0)" @click="props.goBack()">
          <el-icon :class="historyInfo.e('back')"><ac-icon-ep-arrow-left-bold /></el-icon>
        </a>
      </div>
      <div :class="historyInfo.e('title')">
        {{ $t('app.historyLayout.record') }}
      </div>
    </div>

    <div :class="ns.e('left')">
      <div class="flex flex-col h-full overflow-y-auto scroll-content">
        <slot name="left" />
      </div>
    </div>
    <div :class="ns.e('right')" class="scroll-content">
      <router-view />
    </div>
  </main>
</template>

<style lang="scss">
@use '@/styles/mixins/mixins' as *;
@use '@/styles/variable' as *;

// 项目信息
@include b(history-info) {
  height: $doc-header-height;
  width: $doc-layout-left-width;
  padding: 0 $doc-layout-padding;
  @apply flex items-center fixed left-0 top-0 z-50 bg-gray-100;

  @include e(img) {
    @apply flex-none w-32px h-32px mr-10px cursor-pointer;
  }

  @include e(back) {
    @apply w-32px h-32px rounded-4px  text-12px border-1px border-gray border-solid bg-white hover:bg-gray-100;
    width: 32px !important;
    height: 32px !important;
  }

  @include e(title) {
    @apply truncate text-16px relative pr-20px;
  }
}
</style>
