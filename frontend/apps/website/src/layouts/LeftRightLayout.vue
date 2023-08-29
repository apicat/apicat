<template>
  <div class="flex h-full">
    <div class="w-246px b-r b-solid b-gray-110 px-20px">
      <slot name="left"></slot>
    </div>
    <div class="flex-1">
      <div class="m-auto pt-22px px-36px" :style="mainViewStyle">
        <slot></slot>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    mainWidth?: number | string
  }>(),
  {
    mainWidth: 776,
  }
)

const mainViewStyle = computed(() => {
  const { mainWidth } = props
  return {
    width: typeof mainWidth !== 'number' ? mainWidth : `${mainWidth}px`,
  }
})
</script>
