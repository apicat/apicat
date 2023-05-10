<template>
  <div class="flex min-h-screen">
    <main class="p-2 ac-login">
      <div class="shadow-xl ac-login__box">
        <div class="ac-login__logo">
          <router-link class="inline-flex items-center" to="/"> <AcLogo /> </router-link>
        </div>
        <h2 class="mb-4 text-xl font-medium text-zinc-800">{{ $t('app.common.register') }}</h2>

        <el-form
          @keyup.enter="onRegistBtnClick(authForm)"
          @submit.prevent="onRegistBtnClick(authForm)"
          label-position="top"
          size="large"
          :rules="rules"
          ref="authForm"
          :model="form"
        >
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
            <el-input type="password" v-model="form.password" :placeholder="$t('app.rules.password.required')" show-password autocomplete="on" />
          </el-form-item>

          <div class="mt-7">
            <el-button :loading="isLoading" @click="onRegistBtnClick(authForm)" class="w-full" type="primary">{{ $t('app.common.register') }}</el-button>
          </div>
        </el-form>

        <p class="mt-8 text-center">
          <router-link :to="LOGIN_PATH" class="mx-1 text-blue-600">{{ $t('app.common.login') }}</router-link>
        </p>
      </div>
    </main>
  </div>
</template>
<script setup lang="ts">
import { LOGIN_PATH } from '@/router'
import { useApi } from '@/hooks/useApi'
import { useUserStore } from '@/store/user'
import { isEmail } from '@apicat/shared'
import { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authForm = shallowRef()

const form = reactive({
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

const [isLoading, userRegisterRequest] = useApi(useUserStore().register)

const onRegistBtnClick = async (formIns: FormInstance) => {
  try {
    await formIns.validate()
    await userRegisterRequest(toRaw(form))
  } catch (error) {
    //
  }
}
</script>
