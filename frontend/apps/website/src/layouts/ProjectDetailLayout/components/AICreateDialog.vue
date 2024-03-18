<script setup lang="ts">
import { useI18n } from 'vue-i18n'

export interface Options {
  type: 'collection' | 'schema'
  onOk?: (args: { showLoading: () => void; hideLoading: () => void; prompt: string }) => Promise<void> | void
  onCancel?: () => void
}

const { t } = useI18n()
const promptTextareaRef = ref()
const visible = ref(false)
const loading = ref(false)
const prompt = ref('')
const typeRef = ref<'collection' | 'schema'>('collection')
const isCollection = computed(() => typeRef.value === 'collection')
const title = computed(() => (isCollection.value ? t('app.project.collection.ai.title') : t('app.schema.ai.title')))
const placeholder = computed(() =>
  isCollection.value ? t('app.project.collection.ai.aiPromptPlaceholder') : t('app.schema.ai.aiPromptPlaceholder'),
)

const defaultOption: Options = {
  type: 'collection',
  onOk: () => {},
  onCancel: () => {},
}

function onEnter() {
  defaultOption.onOk?.({
    showLoading,
    hideLoading,
    prompt: prompt.value,
  })
}

function showLoading() {
  loading.value = true
}

function hideLoading() {
  loading.value = false
}

function reset() {
  prompt.value = ''
  defaultOption.type = 'collection'
  defaultOption.onOk = undefined
  defaultOption.onCancel = undefined
}

watch(visible, (val) => {
  if (!val) {
    defaultOption.onCancel?.()
    reset()
  }
})

defineExpose({
  show(option: Options = defaultOption) {
    const { type } = option
    typeRef.value = type
    visible.value = true

    defaultOption.onOk = option.onOk
    defaultOption.onCancel = option.onCancel
    nextTick(() => promptTextareaRef.value?.focus())
  },
  hide() {
    visible.value = false
  },
})
</script>

<template>
  <el-dialog
    v-model="visible"
    :lock-scroll="false"
    :close-on-click-modal="false"
    destroy-on-close
    class="ai-dialog"
    width="800px">
    <template #header>
      <p class="flex items-center">
        <Iconfont class="text-blue" icon="ac-zhinengyouhua" />
        <span class="font-500 ml-10px text-16px">
          {{ title }}
        </span>
      </p>
    </template>
    <PromptTextarea
      ref="promptTextareaRef"
      v-model="prompt"
      :placeholder="placeholder"
      :loading="loading"
      @submit="onEnter" />
  </el-dialog>
</template>

<style lang="scss">
.ai-dialog {
  .el-dialog__header {
    padding-bottom: 20px;
  }

  .el-dialog__body {
    padding-left: 25px;
    padding-right: 25px;
  }
}
</style>
