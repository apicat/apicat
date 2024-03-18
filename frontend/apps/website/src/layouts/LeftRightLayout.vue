<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    mainWidth?: number | string
  }>(),
  {
    mainWidth: 750,
  },
)

const mainViewStyle = computed(() => {
  const { mainWidth } = props
  return {
    width: typeof mainWidth !== 'number' ? mainWidth : `${mainWidth}px`,
  }
})
</script>

<template>
  <div class="flex h-full">
    <div class="w-246px b-r b-solid b-gray-110 px-20px">
      <slot name="left" />
    </div>

    <div class="fixed right-0 overflow-auto left-326px h-100vh">
      <div class="m-auto pt-22px" :style="mainViewStyle">
        <slot />
      </div>
    </div>
  </div>
</template>
