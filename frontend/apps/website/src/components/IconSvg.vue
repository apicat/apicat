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

const style = computed(() => {
  const a: any = {
    fontSize: props.width + 'px',
  }
  if (props.height) {
    a.height = props.height + 'px'
  }
  return a
})
</script>

<template>
  <svg class="icon" :style="style" aria-hidden="true">
    <use :xlink:href="`#${name}`"></use>
  </svg>
</template>

<style scoped></style>
