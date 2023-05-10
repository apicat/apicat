<template>
  <el-form @submit.prevent="onSubmit" @keyup.enter="onSubmit" ref="formRef" :model="form" :rules="rules" label-position="top" class="max-w-sm">
    <el-form-item :label="$t('app.form.user.oldPassword')" prop="password">
      <el-input type="password" v-model="form.password" :placeholder="$t('app.rules.password.requiredOld')" autocomplete="on" />
    </el-form-item>

    <el-form-item :label="$t('app.form.user.newPassword')" prop="new_password">
      <el-input type="password" v-model="form.new_password" :placeholder="$t('app.rules.password.requiredNew')" autocomplete="on" />
    </el-form-item>

    <el-form-item :label="$t('app.form.user.confirmNewPassword')" prop="confirm_new_password">
      <el-input type="password" v-model="form.confirm_new_password" :placeholder="$t('app.rules.password.requiredConfirm')" autocomplete="on" />
    </el-form-item>

    <el-form-item>
      <el-button type="primary" @click="onSubmit" :loading="isLoading">{{ $t('app.common.save') }}</el-button>
    </el-form-item>
  </el-form>
</template>
<script setup lang="ts">
import useApi from '@/hooks/useApi'
import { useUserStore } from '@/store/user'
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const userStore = useUserStore()
const formRef = ref()

const form = reactive({
  password: '',
  new_password: '',
  confirm_new_password: '',
})

const minLengthRule = { type: 'string', min: 8, message: t('app.rules.password.minLength'), trigger: 'blur' }

const rules = {
  password: [{ required: true, message: t('app.rules.password.requiredOld'), trigger: 'blur' }, minLengthRule],
  new_password: [{ required: true, message: t('app.rules.password.requiredNew'), trigger: 'blur' }, minLengthRule],
  confirm_new_password: [
    {
      required: true,
      message: t('app.rules.password.requiredConfirm'),
      trigger: 'blur',
    },
    minLengthRule,
    {
      validator: (rule: any, value: string, callback: any) => {
        if (value !== form.new_password) {
          callback(new Error(t('app.rules.password.noMatch')))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
} as any

const [isLoading, modifyPasswordRequest] = useApi(userStore.modifyUserPassword)

const onSubmit = async () => {
  try {
    const valid = await formRef?.value.validate()
    valid && (await modifyPasswordRequest(form))
  } catch (error) {
    //
  }
}
</script>
