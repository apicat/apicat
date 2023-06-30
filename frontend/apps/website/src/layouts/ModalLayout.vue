<template>
  <div class="flex h-700px">
    <div class="bg-gray-100 pt-20px pl-30px" :style="{ width: `${slideWidth}px` }">
      <slot name="nav"></slot>
    </div>
    <div class="relative flex-1 h-full overflow-y-scroll px-30px pb-20px">
      <p class="sticky top-0 z-10 bg-white text-16px text-gray-950 font-500 py-20px">
        <slot name="title"></slot>
      </p>
      <slot></slot>
    </div>
  </div>
</template>
<script setup>
defineProps({
  slideWidth: {
    type: Number,
    default: 300,
  },
})
</script>
