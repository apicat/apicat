<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue?: string
    placeholder: string
    loading?: boolean
    maxHeight?: number
    minHeight?: number
    prefix?: boolean
    emptyTrigger?: boolean
    autoFocus?: boolean
  }>(),
  {
    modelValue: '',
    placeholder: '',
    maxHeight: 600,
    minHeight: 40,
    loading: false,
    emptyTrigger: false,
  },
)

const emits = defineEmits<{
  (e: 'update:modelValue' | 'submit', value: string): void
}>()

const prompt = useVModel(props, 'modelValue', emits)
const inputRef = ref<HTMLTextAreaElement>()
const maxHeight = computed(() => `${props.maxHeight}px`)
const disabled = computed(() => {
  if (!prompt.value && props.emptyTrigger)
    return false

  return !prompt.value
})

const textareaStyle = computed(() => {
  // 上下间距
  const paddingY = (props.minHeight - 14 * 2) / 2
  const style: Record<string, any> = {
    maxHeight: maxHeight.value,
    paddingTop: `${paddingY}px`,
    paddingBottom: `${paddingY}px`,
  }
  if (props.prefix)
    style.paddingLeft = '40px'
  return style
})

function onInput(e: Event) {
  const target = e.target as HTMLElement
  target.style.overflowY = target.scrollHeight > 600 ? 'auto' : 'hidden'
  target.style.maxHeight = maxHeight.value
  target.style.height = '0'
  target.style.height = `${target.scrollHeight}px`
}

function keyPress(e: KeyboardEvent) {
  if (!e.shiftKey && !e.ctrlKey && !e.altKey) {
    e.preventDefault()
    return onSubmit()
  }
}

function onSubmit() {
  if (!props.loading && (prompt.value || (!prompt.value && props.emptyTrigger)))
    emits('submit', prompt.value)
}

onMounted(() => {
  if (props.autoFocus)
    inputRef.value?.focus()
})

defineExpose({
  focus() {
    setTimeout(() => inputRef.value?.focus(), 0)
  },
})

const containRef = ref<HTMLElement>()
function onFocusIn() {
  containRef.value!.style.border = '2px #98C4FF solid'
}
function onFocusOut() {
  containRef.value?.attributes.removeNamedItem('style')
}
</script>

<template>
  <fieldset ref="containRef" class="test-input">
    <div v-if="prefix" class="absolute flex items-center left-15px top-12px w-16px h-16px">
      <Iconfont class="text-blue text-20px" icon="ac-zhinengyouhua" />
    </div>
    <textarea
      ref="inputRef"
      v-model="prompt"
      :readonly="loading"
      :placeholder="placeholder"
      :style="textareaStyle"
      @keypress.enter="keyPress"
      @focusin="onFocusIn"
      @focusout="onFocusOut"
      @input="onInput"
    />
    <div
      class="absolute items-center right-15px bottom-8px"
      :class="{
        'cursor-pointer': !loading && prompt,
      }"
      @click="onSubmit"
    >
      <Iconfont
        v-if="!loading"
        class="text-primary text-20px transition-all-300"
        :class="disabled ? 'text-gray-300' : 'text-primary'"
        icon="ac-send"
      />
      <ac-icon-eos-icons:three-dots-loading v-else class="text-blue text-20px" />
    </div>
  </fieldset>
</template>

<style lang="scss" scoped>
.test-input:hover {
  border: 2px #dbdbdb solid;
}

.test-input {
  cursor: text;
  overflow: hidden;
  display: flex;
  position: relative;
  border: 2px #e7e7e7 solid;
  border-radius: 10px;
  transition: 0.1s all;
  box-sizing: border-box;

  textarea {
    width: 100%;
    box-sizing: border-box;
    border: none;
    font-size: 14px;
    resize: none;
    overflow-y: hidden;
    height: 40px;
    line-height: 2;
    padding-left: 10px;
    padding-right: 30px;
    outline: none;
    overflow-y: hidden;
    word-break: break-all;

    &:focus {
      outline: none;
      box-shadow: 0 0 0 0 #000;
    }
  }

  .enter-btn {
    position: absolute;
    right: 14px;
    bottom: 12px;
    width: 30px;
    height: 30px;
  }
}
</style>
