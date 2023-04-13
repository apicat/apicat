<template>
  <el-input v-model="editingValue" :disabled="disabled" :placeholder="placeholder" ref="inputRef" @focus="openInput" @change="update" />
</template>

<script setup lang="ts">
import { nextTick, ref, shallowRef } from 'vue'

const props = defineProps<{
  value?: string
  disabled?: boolean
  placeholder?: string
}>()

const emits = defineEmits(['change'])
const editingValue = ref(props.value)
const inputRef = shallowRef()
const openInput = () => {
  if (props.disabled) {
    return
  }
  editingValue.value = props.value || ''
  nextTick(() => {
    inputRef.value.focus()
  })
}

const update = () => {
  const v = editingValue.value?.trim()
  if (props.value != v) {
    emits('change', v)
  }
}

defineExpose({
  focus() {
    openInput()
  },
})
</script>
