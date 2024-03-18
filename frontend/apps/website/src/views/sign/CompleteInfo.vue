<script setup lang="ts">
import { isEmail } from '@apicat/shared'
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'
import AheadNav from '@/views/component/AheadNav.vue'
import { useUserStore } from '@/store/user'
import useApi from '@/hooks/useApi'
import { DEFAULT_LANGUAGE } from '@/commons'

const { t } = useI18n()
const route = useRoute()
const authForm = shallowRef()
const form = reactive<SignAPI.RequestRegister>({
  name: (route.query.name as string) || '',
  email: (route.query.email as string) || '',
  password: '',
  invitationToken: route.query.invitationToken as string,
  language: DEFAULT_LANGUAGE,
})

const rules = {
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

async function onCompleteSubmit(formIns: FormInstance) {
  try {
    await formIns.validate()
    await userRegisterRequest({
      ...form,
      bind: {
        oauthUserID: route.query.oauthUserID as string,
        type: route.params.type as string,
      },
    })
  } catch (err) {}
}
</script>

<template>
  <AheadNav>
    <div class="flex items-center justify-center w-full text-center">
      <div class="mt-20">
        <div>
          <h1 class="font-500 text-24px">
            {{ $t('app.sign.completeInfo') }}
          </h1>
        </div>
        <div class="mt-10">
          <el-form
            ref="authForm"
            label-position="top"
            size="large"
            :rules="rules"
            :model="form"
            @keyup.enter="onCompleteSubmit(authForm)"
            @submit.prevent>
            <el-form-item label="" prop="email">
              <div class="ac-login__label">
                <span>{{ $t('app.form.user.email') }}</span>
              </div>
              <el-input v-model="form.email" maxlength="255" autocomplete="off" />
            </el-form-item>

            <el-form-item label="" prop="password">
              <div class="ac-login__label">
                <span>{{ $t('app.form.user.password') }}</span>
              </div>
              <el-input v-model="form.password" maxlength="255" type="password" autocomplete="off" />
            </el-form-item>

            <div class="mt-7">
              <el-button class="w-full" type="primary" :loading="isLoading" @click="onCompleteSubmit(authForm)">
                {{ $t('app.sign.completeInfoSend') }}
              </el-button>
            </div>
          </el-form>
        </div>
      </div>
    </div>
  </AheadNav>
</template>
