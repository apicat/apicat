<template>
  <div class="relative overflow-hidden" @click="onFileWrapperClick">
    <input v-if="!readonly" type="file" ref="fileInputRef" class="absolute top-0 left-0 w-full h-full opacity-0 -z-1" :accept="accept" @change="onFileChange" />
    <slot :fileName="fileName"></slot>
  </div>
</template>
<script setup lang="ts">
const emits = defineEmits(['change'])
const props = defineProps({
  readonly: {
    type: Boolean,
    default: false,
  },

  accept: {
    type: String,
    default: '*',
  },
})
const fileInputRef = ref<HTMLInputElement | null>(null)
const fileName = ref()
const onFileChange = (event: Event) => {
  const selectedFiles: File[] = Array.from((event.target as HTMLInputElement).files!)
  const file = selectedFiles.length ? selectedFiles[0] : null
  fileName.value = file?.name
  emits('change', file)
}

const onFileWrapperClick = () => {
  if (props.readonly) {
    return
  }
  fileInputRef.value?.click()
}

defineExpose({
  clear: () => {
    fileName.value = ''
  },
})
</script>
