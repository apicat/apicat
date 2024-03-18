<script setup>
defineProps({
  slideWidth: {
    type: Number,
    default: 260,
  },
})
</script>

<template>
  <div class="flex overflow-hidden rounded h-700px">
    <div class="bg-gray-100 pt-20px pl-30px" :style="{ width: `${slideWidth}px` }">
      <slot name="nav" />
    </div>
    <div class="relative flex-1 h-full overflow-y-auto scroll-content px-30px pb-20px">
      <p class="sticky top-0 z-10 bg-white text-16px text-gray-950 font-500 py-20px">
        <slot name="title" />
      </p>
      <slot />
    </div>
  </div>
</template>
