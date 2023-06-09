<template>
  <el-dialog v-model="dialogVisible" :title="$t('app.member.tips.editMember')" :width="300" align-center :close-on-click-modal="false">
    <!-- 内容 -->
    <el-form label-position="top" label-width="100px" :model="form" :rules="rules" ref="fromRef" @submit.prevent="handleSubmit(fromRef)">
      <el-form-item :label="$t('app.form.user.username')" prop="username">
        <el-input disabled v-model="form.username" maxlength="30" :placeholder="$t('app.rules.username.required')"></el-input>
      </el-form-item>

      <el-form-item :label="$t('app.form.user.email')" prop="email">
        <el-input v-model="form.email" maxlength="30" :placeholder="$t('app.rules.email.required')" clearable></el-input>
      </el-form-item>

      <el-form-item :label="$t('app.form.user.password')" prop="password">
        <el-input show-password type="password" v-model="form.password" :placeholder="$t('app.rules.password.required')" clearable />
      </el-form-item>

      <el-form-item :label="$t('app.member.form.accountStatus')" prop="role">
        <el-switch
          v-model="form.is_enabled"
          :active-value="1"
          :inactive-value="0"
          :active-text="$t('app.member.form.accountStatusNormal')"
          :inactive-text="$t('app.member.form.accountStatusLock')"
        />
      </el-form-item>
    </el-form>
    <!-- 底部按钮 -->
    <div slot="footer" class="text-right mt-20px">
      <el-button @click="dialogVisible = false">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(fromRef)">
        {{ $t('app.common.save') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { useModal } from '@/hooks'
import { useI18n } from 'vue-i18n'
import type { FormInstance, FormRules } from 'element-plus'
import { updateMember } from '@/api/member'
import { UserInfo } from '@/typings/user'
import { isEmail } from '@apicat/shared'
import useApi from '@/hooks/useApi'

const emits = defineEmits(['ok'])

const { t } = useI18n()

const fromRef = ref<FormInstance>()
const { dialogVisible, showModel, hideModel } = useModal(fromRef as any)
const [isLoading, updateMemberRequest] = useApi(updateMember)

const form = reactive<Partial<UserInfo>>({
  username: '',
  email: '',
  password: '',
})

const rules = reactive<FormRules>({
  email: [
    {
      validator(rule: any, value: any, callback: any) {
        if (value && !isEmail(value)) {
          return callback(new Error(t('app.rules.email.correct')))
        }
        return callback()
      },
      trigger: 'blur',
    },
  ],
  password: [{ min: 8, message: t('app.rules.password.minLength'), trigger: 'blur' }],
})

const handleSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  try {
    const valid = await formEl.validate()
    if (valid) {
      const user: UserInfo = await updateMemberRequest(toRaw(form))
      hideModel()
      emits('ok', user)
    }
  } catch (error) {
    //
  }
}

const show = (user: UserInfo) => {
  form.id = user.id
  form.username = user.username
  form.is_enabled = user.is_enabled
  showModel()
}
defineExpose({
  show,
})
</script>
