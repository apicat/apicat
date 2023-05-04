<template>
  <el-dialog v-model="dialogVisible" center append-to-body :close-on-click-modal="false" :close-on-press-escape="false" destroy-on-close width="40%">
    <template #header>
      <div class="flex-y-center">
        <el-icon class="mr-5px"><ac-icon-bi-robot /></el-icon>{{ $t('app.schema.common.aiGenerateSchema') }}
      </div>
    </template>

    <el-form class="mt-10px" label-position="top" label-width="0" :model="form" :rules="rules" ref="aiPromptForm" @submit.prevent="handleSubmit(aiPromptForm)">
      <el-form-item prop="name">
        <el-input size="large" v-model="form.name" :placeholder="$t('app.schema.tips.schemaInputTitle')" clearable />
      </el-form-item>
    </el-form>
    <el-button :loading="isLoading" type="primary" @click="handleSubmit(aiPromptForm)">{{ $t('app.common.generate') }}</el-button>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import useApi from '@/hooks/useApi'
import { useParams } from '@/hooks/useParams'
import { FormInstance } from 'element-plus'
import { aiGenerateDefinition } from '@/api/definition'
import { useI18n } from 'vue-i18n'

const emits = defineEmits(['ok'])
const { t } = useI18n()
const aiPromptForm = ref<FormInstance>()

const { dialogVisible, showModel, hideModel } = useModal(aiPromptForm as any)
const [isLoading, aiGenerateDefinitionApi] = useApi(aiGenerateDefinition)()
const { project_id } = useParams()

let otherParams = {}

const form = reactive({
  parent_id: 0,
  name: '',
})

const rules = {
  name: [{ required: true, message: t('app.schema.tips.schemaInputTitle'), trigger: 'blur' }],
}

const show = (params?: any) => {
  otherParams = { ...params }
  showModel()
}

const handleSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return

  try {
    const valid = await formEl.validate()
    if (valid) {
      const data = await aiGenerateDefinitionApi({ project_id, ...toRaw(form), ...otherParams })
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
