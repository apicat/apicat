<script setup lang="ts">
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'
import AheadNav from '@/views/component/AheadNav.vue'
import useApi from '@/hooks/useApi'
import { apiResetPass, checkResetPassCodeIsExpired } from '@/api/sign/forgot'
import { LOGIN_PATH } from '@/router'
import { BadRequestError } from '@/api/error'
import { useInitedPageWithGlobalLoading } from '@/hooks'

const { t } = useI18n()
const route = useRoute()
const authForm = shallowRef()
const isExpriredCode = ref(false)
const isShowSuccessTip = ref(false)
const code = route.params.token as string
const form = reactive<{ password: string; re_password: string }>({ password: '', re_password: '' })
const failMsg = ref<CommonResponseMessageForMessageTemplate>()

const rules = {
  password: [
    {
      required: true,
      type: 'string',
      min: 8,
      message: t('app.rules.password.minLength'),
      trigger: 'blur',
    },
  ],
  re_password: [
    {
      required: true,
      validator(__: any, _: any, callback: any) {
        if (form.password !== form.re_password) callback(new Error(t('app.rules.password.noMatch')))
        else if (form.re_password.length === 0) callback(new Error(t('app.rules.password.minLength')))
        else callback()
      },
      trigger: 'blur',
    },
  ],
} as any

const [isLoading, resetPass] = useApi(apiResetPass)

async function onResetSubmit(formIns: FormInstance) {
  try {
    await formIns.validate()
    await resetPass(form, code)
    isShowSuccessTip.value = true
  } catch (err) {
    //
  }
}
onBeforeMount(
  useInitedPageWithGlobalLoading(async () => {
    try {
      await checkResetPassCodeIsExpired(code)
      isExpriredCode.value = false
    } catch (error) {
      isExpriredCode.value = true
      if (error instanceof BadRequestError) {
        const { response } = error as BadRequestError<CommonResponseMessageForMessageTemplate>
        failMsg.value = response || {}
      }
    }
  }),
)
</script>

<template>
  <AheadNav v-if="!isExpriredCode">
    <div v-if="!isShowSuccessTip" class="flex items-center justify-center w-full text-center">
      <div class="mt-20">
        <div>
          <h1 class="font-500 text-24px">
            {{ $t('app.sign.resetPass') }}
          </h1>
        </div>
        <div class="mt-10">
          <el-form
            ref="authForm"
            label-position="top"
            size="large"
            :rules="rules"
            :model="form"
            @keyup.enter="onResetSubmit(authForm)"
            @submit.prevent>
            <el-form-item label="" prop="password">
              <div class="ac-login__label">
                <span>{{ $t('app.sign.resetPassNew') }}</span>
              </div>
              <el-input maxlength="255" v-model="form.password" type="password" autocomplete="off" />
            </el-form-item>

            <el-form-item label="" prop="re_password">
              <div class="ac-login__label">
                <span>{{ $t('app.sign.resetPassRepeat') }}</span>
              </div>
              <el-input v-model="form.re_password" type="password" autocomplete="off" />
            </el-form-item>

            <div class="mt-7">
              <el-button
                maxlength="255"
                class="w-full"
                type="primary"
                :loading="isLoading"
                @click="onResetSubmit(authForm)">
                {{ $t('app.sign.resetPassSend') }}
              </el-button>
            </div>
          </el-form>
        </div>
      </div>
    </div>

    <MessageTemplate v-else hide-header title="Password reset successful" emoji="🎉">
      <template #description>
        <AcCountDown v-slot="{ seconds }" :time="5" :link="LOGIN_PATH" auto-jump>
          <p>
            The system will automatically jump to the login page after {{ seconds }} seconds. If not, please
            <RouterLink class="text-primary" :replace="true" :to="LOGIN_PATH"> click here </RouterLink>.
          </p>
        </AcCountDown>
      </template>
    </MessageTemplate>
  </AheadNav>
  <MessageTemplate v-else v-bind="failMsg" />
</template>
