<template>
  <el-dialog v-model="dialogVisible" title="项目分组" :width="348" align-center :close-on-click-modal="false">
    <!-- 内容 -->
    <el-form label-position="top" label-width="100px" :model="form" ref="formRef" @submit.prevent="handleSubmit(formRef)">
      <el-form-item label="" prop="id">
        <el-select v-model="form.id" class="w-full">
          <el-option v-for="item in groupsForOptions" :key="item.id" :label="item.name" :value="item.id!" />
        </el-select>
      </el-form-item>
    </el-form>
    <!-- 底部按钮 -->
    <div class="text-right -mb-10px">
      <el-button @click="dialogVisible = false">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(formRef)">
        {{ $t('app.common.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { useModal } from '@/hooks'
import { useI18n } from 'vue-i18n'
import type { FormInstance } from 'element-plus'
import useProjectGroupStore from '@/store/projectGroup'
import { storeToRefs } from 'pinia'
import { settingProjectGroup } from '@/api/project'
import { ProjectInfo } from '@/typings'

type FormState = { id: number; projectId: string }

const emits = defineEmits(['success'])

const { t } = useI18n()

const form = reactive<FormState>({
  id: 0,
  projectId: '',
})

const formRef = ref<FormInstance>()
const isLoading = ref(false)
const { dialogVisible, showModel, hideModel } = useModal(formRef as any)
const projectGroupStore = useProjectGroupStore()
const { groupsForOptions } = storeToRefs(projectGroupStore)

const handleSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  try {
    const valid = await formEl.validate()
    if (valid) {
      isLoading.value = true
      const { projectId, id } = form
      await settingProjectGroup(projectId, id)
      emits('success')
      hideModel()
    }
  } catch (error) {
    //
  } finally {
    isLoading.value = false
  }
}

const show = async (projectInfo: ProjectInfo) => {
  showModel()
  await nextTick()
  formRef.value?.clearValidate()
  Object.assign(form, {
    id: projectInfo.group_id || 0,
    projectId: projectInfo.id,
  })
}

defineExpose({
  show,
})
</script>
