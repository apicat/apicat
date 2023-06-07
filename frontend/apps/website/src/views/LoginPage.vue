<template>
  <div class="flex min-h-screen">
    <main class="p-2 ac-login">
      <div class="shadow-xl ac-login__box">
        <div class="ac-login__logo">
          <router-link class="inline-flex items-center" to="/"> <AcLogo /> </router-link>
        </div>
        <h2 class="mb-4 text-xl font-medium text-zinc-800">{{ $t('app.common.login') }}</h2>

        <el-form label-position="top" size="large" :rules="rules" ref="authForm" :model="form" @keyup.enter="onLoginBtnClick(authForm)" @submit.prevent="onLoginBtnClick(authForm)">
          <el-form-item label="" prop="email">
            <div class="ac-login__label">
              <span>{{ $t('app.form.user.email') }}</span>
            </div>
            <el-input v-model="form.email" :placeholder="$t('app.rules.email.required')" autocomplete="on" />
          </el-form-item>

          <el-form-item label="" prop="password">
            <div class="ac-login__label">
              <span>{{ $t('app.form.user.password') }}</span>
            </div>
            <el-input type="password" v-model="form.password" :placeholder="$t('app.rules.password.required')" autocomplete="on" show-password />
          </el-form-item>

          <div class="mt-7">
            <el-button :loading="isLoading" @click="onLoginBtnClick(authForm)" class="w-full" type="primary">{{ $t('app.common.login') }}</el-button>
          </div>
        </el-form>

        <p class="mt-8 text-center">
          <router-link :to="REGISTER_PATH" class="mx-1 text-blue-600">{{ $t('app.common.registerAccount') }}</router-link>
        </p>
      </div>
    </main>
  </div>
</template>
<script setup lang="ts">
import { REGISTER_PATH } from '@/router'
import { useApi } from '@/hooks/useApi'
import { useUserStore } from '@/store/user'
import { isEmail } from '@apicat/shared'
import { FormInstance } from 'element-plus'
import { UserInfo } from '@/typings/user'
import { useI18n } from 'vue-i18n'

const authForm = shallowRef()

const { t } = useI18n()
const form: Partial<UserInfo> = reactive({
  email: '',
  password: '',
})

const rules = {
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
  password: [{ type: 'string', min: 8, message: t('app.rules.password.minLength'), trigger: 'blur' }],
} as any

const [isLoading, loginRequest] = useApi(useUserStore().login)

const onLoginBtnClick = async (formIns: FormInstance) => {
  try {
    await formIns.validate()
    await loginRequest(toRaw(form))
  } catch (error) {
    //
  }
}
</script>
