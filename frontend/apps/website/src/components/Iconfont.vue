<script setup lang="ts">
const props = defineProps({
  icon: {
    type: String,
    default: '',
  },
  size: {
    type: Number,
    default: 16,
  },
})

const style = computed(() => ({
  fontSize: `${props.size}px`,
}))
</script>

<template>
  <i class="leading-none ac-iconfont" :class="icon" :style="style" />
</template>
