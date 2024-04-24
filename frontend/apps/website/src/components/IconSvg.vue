<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    name: string
    width?: string | number
    height?: string | number
  }>(),
  {
    width: '24',
  },
)

const style = computed(() => ({
  fontSize: `${props.width}px`,
  height: `${props.height ?? props.width}px`,
}))
</script>

<template>
  <svg class="icon" :style="style" aria-hidden="true" :width="width">
    <use :xlink:href="`#${name}`" />
  </svg>
</template>

<style scoped></style>
