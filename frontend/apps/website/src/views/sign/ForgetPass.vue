<script setup lang="ts">
import { ElMessage, type FormInstance } from 'element-plus'
import { isEmail } from '@apicat/shared'
import { useI18n } from 'vue-i18n'
import { LOGIN_PATH } from '@/router'
import AheadNav from '@/views/component/AheadNav.vue'
import useApi from '@/hooks/useApi'
import { apiForgotSendEmail } from '@/api/sign/forgot'

const { t } = useI18n()

const authForm = shallowRef()
const form = reactive<{ email: string }>({
  email: '',
})

const rules = {
  email: [
    { required: true, message: t('app.rules.email.required'), trigger: 'blur' },
    {
      validator(rule: any, value: any, callback: any) {
        if (!isEmail(value)) callback(new Error(t('app.rules.email.correct')))
        else callback()
      },
      trigger: 'blur',
    },
  ],
} as any

const [isLoading, sendEmail] = useApi(apiForgotSendEmail)

async function onForgotSubmit(formIns: FormInstance) {
  try {
    await formIns.validate()
    await sendEmail(form.email)
    ElMessage.success(t('app.user.password.success'))
  } catch (err) {}
}
</script>

<template>
  <AheadNav>
    <div class="flex items-center justify-center w-full text-center">
      <div class="mt-20">
        <div>
          <h1 class="font-500 text-24px">
            {{ $t('app.sign.forgotPass') }}
          </h1>
        </div>
        <div class="mt-10">
          <el-form
            ref="authForm"
            label-position="top"
            size="large"
            :rules="rules"
            :model="form"
            @keyup.enter="onForgotSubmit(authForm)"
            @submit.prevent>
            <el-form-item label="" prop="email">
              <div class="ac-login__label">
                <span>{{ $t('app.sign.forgotPassEmail') }}</span>
              </div>
              <el-input maxlength="255" v-model="form.email" autocomplete="on" />
            </el-form-item>

            <div class="mt-7">
              <el-button :loading="isLoading" class="w-full" type="primary" @click="onForgotSubmit(authForm)">
                {{ $t('app.sign.forgotPassSend') }}
              </el-button>
            </div>
          </el-form>

          <p class="gap-3 mt-8 text-center">
            <router-link :to="LOGIN_PATH" class="mx-1 text-blue-600">
              {{ $t('app.common.login') }}
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </AheadNav>
</template>
