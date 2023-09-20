<template>
  <el-dialog v-model="dialogVisible" center append-to-body :close-on-click-modal="false" :close-on-press-escape="false" destroy-on-close width="40%">
    <template #header>
      <div class="flex-y-center"><Iconfont class="mr-5px" icon="ac-zhinengyouhua" /> {{ $t('app.interface.common.aiGenerateInterface') }}</div>
    </template>

    <el-form class="mt-10px" label-position="top" label-width="0" :model="form" :rules="rules" ref="aiPromptForm" @submit.prevent="handleSubmit(aiPromptForm)">
      <el-form-item prop="title">
        <el-input size="large" ref="titleInputRef" v-model="form.title" :placeholder="$t('app.interface.form.modalTitle')" clearable />
      </el-form-item>
    </el-form>
    <el-button :loading="isLoading" type="primary" @click="handleSubmit(aiPromptForm)">生成</el-button>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import { createCollectionByAI } from '@/api/collection'
import useApi from '@/hooks/useApi'
import { useParams } from '@/hooks/useParams'
import { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'

const emits = defineEmits(['ok'])
const { t } = useI18n()

const aiPromptForm = ref<FormInstance>()
const titleInputRef = ref()

const { dialogVisible, showModel, hideModel } = useModal(aiPromptForm as any)
const [isLoading, createCollectionByAIApi] = useApi(createCollectionByAI)
const { project_id, iteration_id } = useParams()

let otherParams = {}

const form = reactive({
  title: '',
})

const rules = {
  title: [{ required: true, message: t('app.interface.form.modalTitle'), trigger: 'blur' }],
}

const show = (params?: any) => {
  otherParams = { ...params }
  showModel()
  setTimeout(() => titleInputRef.value.focus(), 0)
}

const handleSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  try {
    const valid = await formEl.validate()
    if (valid) {
      const data = await createCollectionByAIApi({ project_id, iteration_id, title: form.title, ...otherParams })

      emits('ok', data.id)
      reset()
      hideModel()
    }
  } catch (error) {
    //
  }
}

const reset = () => {
  otherParams = {}
}

defineExpose({
  show,
})
</script>
