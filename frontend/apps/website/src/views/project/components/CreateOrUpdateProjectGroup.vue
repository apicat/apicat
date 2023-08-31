<template>
  <el-dialog v-model="dialogVisible" :title="titleRef" :width="348" align-center :close-on-click-modal="false">
    <!-- 内容 -->
    <el-form label-position="top" label-width="100px" :model="form" :rules="rules" ref="groupFormRef" @submit.prevent="handleSubmit(groupFormRef)">
      <el-form-item label="" prop="name">
        <el-input ref="inputRef" v-model="form.name" placeholder="请输入分组名称" clearable maxlength="255" />
      </el-form-item>
    </el-form>
    <!-- 底部按钮 -->
    <div class="text-right -mb-10px">
      <el-button @click="dialogVisible = false">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(groupFormRef)">
        {{ $t('app.common.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { useModal } from '@/hooks'
import { useI18n } from 'vue-i18n'
import type { FormInstance, FormRules } from 'element-plus'
import { ProjectGroup } from '@/typings'
import useProjectGroupStore from '@/store/projectGroup'

const emits = defineEmits(['success'])

const { t } = useI18n()

const form = ref<ProjectGroup>({
  name: '',
})

const titleRef = computed(() => (form.value.id ? '编辑分组' : '创建分组'))

const groupFormRef = ref<FormInstance>()
const inputRef = ref()
const isLoading = ref(false)
const { dialogVisible, showModel, hideModel } = useModal(groupFormRef as any)
const projectGroupStore = useProjectGroupStore()

const rules = reactive<FormRules>({
  name: [{ required: true, message: '请输入分组名称', trigger: 'change' }],
})

const handleSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  try {
    const valid = await formEl.validate()
    if (valid) {
      isLoading.value = true
      await projectGroupStore.createOrUpdateProjectGroup(form.value)
      emits('success')
      hideModel()
    }
  } catch (error) {
    //
  } finally {
    isLoading.value = false
  }
}

const show = async (group: ProjectGroup = { name: '' }) => {
  showModel()
  await nextTick()
  groupFormRef.value?.clearValidate()
  setTimeout(() => inputRef.value?.focus(), 0)
  form.value = { ...group }
}

defineExpose({
  show,
})
</script>
