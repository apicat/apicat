<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { FormInstance, FormRules } from 'element-plus'
import { useModal } from '@/hooks'
import useProjectGroupStore from '@/store/projectGroup'

const emits = defineEmits(['success'])
let successCallback: null | ((id: number) => void) = null

const { t } = useI18n()

const form = ref<ProjectAPI.ResponseGroup>({
  name: '',
  id: 0,
})

const isEdit = computed(() => form.value.id)

const titleRef = computed(() => isEdit.value ? t('app.project.groups.editGroup') : t('app.project.groups.createGroup'))

const groupFormRef = ref<FormInstance>()
const inputRef = ref()
const isLoading = ref(false)
const { dialogVisible, showModel, hideModel } = useModal(groupFormRef as any)
const projectGroupStore = useProjectGroupStore()

const rules = reactive<FormRules>({
  name: [
    {
      required: true,
      message: t('app.project.groups.inputGroupNameTip'),
      trigger: 'change',
    },
  ],
})

async function handleSubmit(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  try {
    const valid = await formEl.validate()
    if (valid) {
      isLoading.value = true
      const group = await projectGroupStore.createOrUpdateProjectGroup({ ...toRaw(form.value) })
      emits('success')
      successCallback && successCallback(group?.id as number)
      hideModel()
      successCallback = null
    }
  }
  catch (error) {
    //
  }
  finally {
    isLoading.value = false
  }
}

async function show(group: ProjectAPI.ResponseGroup = { name: '', id: 0 }) {
  showModel()
  await nextTick()
  groupFormRef.value?.clearValidate()
  setTimeout(() => inputRef.value?.focus(), 0)
  form.value = { ...group }
}

function showWithCallback(callback: (group_id: number) => void) {
  successCallback = callback
  show()
}

defineExpose({
  show,
  showWithCallback,
})
</script>

<template>
  <el-dialog v-model="dialogVisible" :title="titleRef" :width="348" align-center :close-on-click-modal="false">
    <!-- 内容 -->
    <el-form
      ref="groupFormRef"
      label-position="top"
      label-width="100px"
      :model="form"
      :rules="rules"
      @submit.prevent="handleSubmit(groupFormRef)"
    >
      <el-form-item label="" prop="name">
        <el-input
          ref="inputRef"
          v-model="form.name"
          :placeholder="$t('app.project.groups.inputGroupName')"
          clearable
          maxlength="255"
        />
      </el-form-item>
    </el-form>
    <!-- 底部按钮 -->
    <div class="text-right -mb-10px">
      <el-button @click="dialogVisible = false">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(groupFormRef)">
        {{ isEdit ? $t('app.common.update') : $t('app.common.create') }}
      </el-button>
    </div>
  </el-dialog>
</template>
