<template>
  <el-form @submit.prevent="onSubmit" @keyup.enter="onSubmit" ref="formRef" :model="userInfo" :rules="rules" label-position="top" class="max-w-sm">
    <el-form-item :label="$t('app.form.user.username')" prop="username">
      <el-input v-model="userInfo.username" maxlength="30" :placeholder="$t('app.rules.username.required')"></el-input>
    </el-form-item>

    <el-form-item :label="$t('app.form.user.email')" prop="email">
      <el-input v-model="userInfo.email" maxlength="30" :placeholder="$t('app.rules.email.required')"></el-input>
    </el-form-item>

    <el-form-item>
      <el-button type="primary" @click="onSubmit" :loading="isLoading">{{ $t('app.common.save') }}</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/store/user'
import useApi from '@/hooks/useApi'
import { UserInfo } from '@/typings/user'
import { isEmail } from '@apicat/shared'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const formRef = ref()

const userStore = useUserStore()

const userInfo: Ref<UserInfo> = ref({ ...userStore.userInfo })

const rules = {
  username: [{ required: true, message: t('app.rules.username.required'), trigger: 'blur' }],
  email: [
    { required: true, message: t('app.rules.email.required'), trigger: 'blur' },
    {
      validator(rule: any, value: any, callback: any) {
        if (!isEmail(value)) {
          callback(new Error(t('app.rules.email.correct')))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
} as any

const [isLoading, modifyUserInfoRequest] = useApi(userStore.modifyUserInfo)

const onSubmit = async () => {
  try {
    const valid = await formRef.value?.validate()
    valid && modifyUserInfoRequest(toRaw(userInfo))
  } catch (error) {
    //
  }
}
</script>
