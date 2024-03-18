<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { FormInstance, FormRules } from 'element-plus'
import { useModal } from '@/hooks'
import { useApi } from '@/hooks/useApi'
import { apiChangeSystemUserPassword } from '@/api/user'

const emits = defineEmits(['success'])

const { t } = useI18n()

const form = ref<{ id: number; password: string;confirmPassword: string }>({
  id: -1,
  password: '',
  confirmPassword: '',
})
const groupFormRef = ref<FormInstance>()
const inputRef = ref()
const { dialogVisible, showModel, hideModel } = useModal(groupFormRef as any)
const [isLoading, changeSystemUserPassword] = useApi(apiChangeSystemUserPassword)

const rules = reactive<FormRules>({
  password: [
    {
      required: true,
      type: 'string',
      min: 8,
      message: t('app.rules.password.minLength'),
      trigger: 'blur',
    },
  ],
  confirmPassword: [
    {
      required: true,
      validator(__: any, _: any, callback: any) {
        if (form.value.password !== form.value.confirmPassword)
          callback(new Error(t('app.rules.password.noMatch')))
        else if (form.value.confirmPassword.length === 0)
          callback(new Error(t('app.rules.password.minLength')))
        else callback()
      },
      trigger: 'blur',
    },
  ],
})

async function handleSubmit(formEl: FormInstance | undefined) {
  if (!formEl)
    return

  try {
    await formEl.validate()
    await changeSystemUserPassword(toRaw(form.value))
    hideModel()
    emits('success')
  }
  catch (error) {

  }
}

async function show(user: UserAPI.ResponseUserInfo) {
  showModel()
  await nextTick()
  groupFormRef.value?.clearValidate()
  setTimeout(() => inputRef.value?.focus(), 0)
  form.value = { id: user.id, password: '', confirmPassword: '' }
}

defineExpose({
  show,
})
</script>

<template>
  <el-dialog v-model="dialogVisible" :title="t('app.system.users.updatePasswordTitle')" :width="348" align-center :close-on-click-modal="false">
    <!-- 内容 -->
    <el-form
      ref="groupFormRef"
      label-position="top"
      label-width="100px"
      :model="form"
      :rules="rules"
      @submit.prevent="handleSubmit(groupFormRef)"
    >
      <!-- new pass -->
      <ElFormItem prop="password" :label="$t('app.user.password.new')">
        <ElInput v-model="form.password" maxlength="255" class="h-40px" type="password" />
      </ElFormItem>

      <!-- confirm pass -->
      <ElFormItem prop="confirmPassword" :label="$t('app.user.password.confirm')">
        <ElInput v-model="form.confirmPassword" maxlength="255" class="h-40px" type="password" />
      </ElFormItem>
    </el-form>
    <!-- 底部按钮 -->
    <div class="text-right -mb-10px">
      <el-button @click="dialogVisible = false">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(groupFormRef)">
        {{ $t('app.common.update') }}
      </el-button>
    </div>
  </el-dialog>
</template>
