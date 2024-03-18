<script setup lang="ts">
import { isEmail } from '@apicat/shared'
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { LOGIN_NAME } from '@/router'
import { useApi } from '@/hooks/useApi'
import { useUserStore } from '@/store/user'
import { DEFAULT_LANGUAGE } from '@/commons'
import { popRedirect } from '@/router/filter/auth.filter'

const route = useRoute()
const { t } = useI18n()

const authForm = shallowRef()
const form = reactive<SignAPI.RequestRegister>({
  name: '',
  email: '',
  password: '',
  invitationToken: route.query.invitationToken as string,
  language: DEFAULT_LANGUAGE,
})

const rules: FormRules = {
  name: [
    {
      required: true,
      message: t('app.rules.username.required'),
      trigger: 'blur',
    },
    {
      validator(_: any, value: string, callback: any) {
        if (value.length < 2 || value.length > 64) callback(new Error(t('app.rules.username.wrongLength')))
        else callback()
      },
      trigger: 'blur',
    },
  ],
  email: [
    { required: true, message: t('app.rules.email.required'), trigger: 'blur' },
    {
      validator(_: any, value: any, callback: any) {
        if (!isEmail(value)) callback(new Error(t('app.rules.email.correct')))
        else callback()
      },
      trigger: 'blur',
    },
  ],
  password: [
    {
      required: true,
      type: 'string',
      min: 8,
      message: t('app.rules.password.minLength'),
      trigger: 'blur',
    },
  ],
} as any

const [isLoading, userRegisterRequest] = useApi(useUserStore().register)
async function onRegistBtnClick(formIns: FormInstance) {
  try {
    await formIns.validate()
    await userRegisterRequest(toRaw(form), popRedirect(route.fullPath))
  } catch (error) {
    //
  }
}
</script>

<template>
  <div class="flex min-h-screen">
    <main class="p-2 ac-login">
      <div class="shadow-xl ac-login__box">
        <div class="ac-login__logo" style="text-align: center">
          <router-link class="inline-flex items-center" to="/">
            <AcLogo style="" />
          </router-link>
        </div>

        <el-form
          ref="authForm"
          label-position="top"
          size="large"
          :rules="rules"
          :model="form"
          @keyup.enter="onRegistBtnClick(authForm)"
          @submit.prevent>
          <el-form-item label="" prop="name">
            <div class="ac-login__label">
              <span>{{ $t('app.form.user.username') }}</span>
            </div>
            <el-input
              v-model="form.name"
              maxlength="64"
              autocomplete="off"
              :placeholder="$t('app.rules.username.required')" />
          </el-form-item>

          <el-form-item label="" prop="email">
            <div class="ac-login__label">
              <span>{{ $t('app.form.user.email') }}</span>
            </div>
            <el-input
              v-model="form.email"
              maxlength="255"
              autocomplete="off"
              :placeholder="$t('app.rules.email.required')" />
          </el-form-item>

          <el-form-item label="" prop="password">
            <div class="ac-login__label">
              <span>{{ $t('app.form.user.password') }}</span>
            </div>
            <el-input
              v-model="form.password"
              type="password"
              show-password
              maxlength="255"
              autocomplete="off"
              :placeholder="$t('app.rules.password.required')" />
          </el-form-item>

          <div class="mt-7">
            <el-button :loading="isLoading" class="w-full" type="primary" @click="onRegistBtnClick(authForm)">
              {{ $t('app.common.register') }}
            </el-button>
          </div>
        </el-form>

        <p class="gap-3 mt-8 text-start">
          <span>{{ $t('app.common.loginTip') }}</span>
          <router-link :to="{ name: LOGIN_NAME, query: route.query }" class="mx-1 text-blue-600">
            {{ $t('app.common.login') }}
          </router-link>
        </p>
      </div>
    </main>
  </div>
</template>
