<script setup lang="ts">
import { isEmail } from '@apicat/shared'
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import { storeToRefs } from 'pinia'
import { FORGETPASS_PATH, REGISTER_NAME } from '@/router'
import { useApi } from '@/hooks/useApi'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'
import { popRedirect } from '@/router/filter/auth.filter'

const route = useRoute()
const authForm = shallowRef()
const { t } = useI18n()

const form: SignAPI.RequestLogin = reactive({
  email: '',
  password: '',
  invitationToken: route.query.invitationToken as string,
})

const rules = {
  email: [
    { required: true, message: t('app.rules.email.required'), trigger: 'blur' },
    {
      validator(rule: any, value: any, callback: any) {
        if (!isEmail(value))
          callback(new Error(t('app.rules.email.correct')))
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

const [isLoading, loginRequest] = useApi(useUserStore().login)
const { oAuthURLConfig, isShowGithubOAuth } = storeToRefs(useAppStore())

async function onLoginBtnClick(formIns: FormInstance) {
  // form validate and login
  try {
    await formIns.validate()
    await loginRequest(toRaw(form), popRedirect(route.fullPath))
  }
  catch (error) {
    //
  }
}

async function githubSign() {
  // jump to github oauth
  window.location.href = oAuthURLConfig.value.github()
}
</script>

<template>
  <div class="flex min-h-screen">
    <main class="p-2 ac-login">
      <div class="shadow-xl ac-login__box">
        <div class="ac-login__logo" style="text-align: center">
          <div class="inline-flex items-center">
            <AcLogo style="" />
          </div>
        </div>

        <el-form
          ref="authForm" label-position="top" size="large" :rules="rules" :model="form"
          @keyup.enter="onLoginBtnClick(authForm)" @submit.prevent
        >
          <el-form-item label="" prop="email">
            <div class="ac-login__label">
              <span>{{ $t('app.form.user.email') }}</span>
            </div>
            <el-input
              v-model="form.email" :placeholder="$t('app.rules.email.required')" autocomplete="on"
              maxlength="255"
            />
          </el-form-item>

          <el-form-item label="" prop="password">
            <div class="ac-login__label">
              <span>{{ $t('app.form.user.password') }}</span>
            </div>
            <el-input
              v-model="form.password" maxlength="255" type="password"
              :placeholder="$t('app.rules.password.required')" autocomplete="on" show-password
            />
          </el-form-item>

          <div class="mt-7">
            <el-button :loading="isLoading" class="w-full" type="primary" @click="onLoginBtnClick(authForm)">
              {{ $t('app.common.login') }}
            </el-button>
          </div>
        </el-form>

        <div v-if="isShowGithubOAuth">
          <ElDivider border-style="solid">
            {{ $t('app.common.loginDivider') }}
          </ElDivider>

          <el-button class="w-full b-btn" type="info" @click="githubSign">
            <Icon class="mr-1" icon="mdi:github" height="24" />
            <span class="ml-1">
              {{ $t('app.common.loginGithub') }}
            </span>
          </el-button>
        </div>

        <p class="gap-3 mt-8 text-start">
          <router-link :to="{ name: REGISTER_NAME, query: route.query }" class="mx-1 text-blue-600">
            {{ $t('app.common.registerAccount') }}
          </router-link>
          <router-link :to="FORGETPASS_PATH" class="mx-1 ml-10 text-blue-600">
            {{ $t('app.sign.forgotPass') }}
          </router-link>
        </p>
      </div>
    </main>
  </div>
</template>

<style lang="scss">
.el-divider--horizontal {
  border-top: 1px #d2d2d2 var(--el-border-style);
}

.b-btn {
  height: 40px;
  background-color: white;
  color: black;
  cursor: pointer;
}
</style>
