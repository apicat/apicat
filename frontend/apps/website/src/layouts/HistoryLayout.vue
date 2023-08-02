<template>
  <main :class="ns.b()">
    <div :class="[historyInfo.b()]">
      <div :class="historyInfo.e('img')">
        <a href="javascript:void(0)" @click="router.go(-1)">
          <el-icon :class="historyInfo.e('back')"><ac-icon-ep-arrow-left-bold /></el-icon>
        </a>
      </div>
      <div :class="historyInfo.e('title')" title="projectDetailInfo?.title">历史记录</div>
    </div>

    <slot name="header"></slot>

    <div :class="ns.e('left')">
      <div class="flex flex-col h-full overflow-y-scroll scroll-content">
        <slot name="left"></slot>
      </div>
    </div>
    <div :class="ns.e('right')" class="scroll-content">
      <router-view />
    </div>
  </main>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks/useNamespace'
const ns = useNamespace('doc-layout')
const historyInfo = useNamespace('history-info')
const router = useRouter()
</script>

<style lang="scss" scoped>
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
  }

  @include e(title) {
    @apply truncate text-16px relative pr-20px;
  }
}
</style>
