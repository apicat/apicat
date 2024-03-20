<script setup lang="ts">
import { ref } from 'vue'
import type { UseCollapse } from './useCollapse'

const props = defineProps<{
  name: string
  collapseCtx: UseCollapse
}>()

const contentRef = ref<HTMLElement>()
const show = ref(false)
const { open, close } = props.collapseCtx.ctx.register(props.name, {
  open() {
    show.value = true
  },
  close() {
    show.value = false
  },
})
function toggle() {
  return show.value ? close() : open()
}
defineExpose({ open, close, toggle })

function expand() {
  contentRef.value!.style.maxHeight = `${contentRef.value!.scrollHeight}px`
}
function fold() {
  contentRef.value!.style.maxHeight = '0px'
}

const transitions = {
  enter() {
    expand()
  },
  leave() {
    fold()
  },
}

onMounted(() => {
  if (show.value)
    expand()
})
</script>

<template>
  <div class="border border-gray-200 border-solid rounded-lg px-25px card-container">
    <div class="py-20px" @click="toggle">
      <slot name="title" />
    </div>
    <transition @enter="transitions.enter" @leave="transitions.leave">
      <div v-show="show" ref="contentRef" class="card-content">
        <div class="pb-25px">
          <slot />
        </div>
      </div>
    </transition>
  </div>
</template>

<style lang="scss" scoped>
* {
  -webkit-transform: translateZ(0);
  -moz-transform: translateZ(0);
  -ms-transform: translateZ(0);
  -o-transform: translateZ(0);
  transform: translateZ(0);
}
.card-content {
  max-height: 0;
  overflow: hidden;
  will-change: max-height, opacity, display;
  transition:
    max-height 0.3s ease-out,
    opacity 0.3s ease-out;
}
</style>
