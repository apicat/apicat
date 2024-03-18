<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { FORGETPASS_PATH } from '@/router'
import { apiResetPassword } from '@/api/user'

const { t } = useI18n()
const formRef = ref<FormInstance>()
const form = ref<UserAPI.RequestResetPassword>({
  password: '',
  newPassword: '',
  reNewPassword: '',
})
const rules: FormRules<typeof form> = {
  password: [
    {
      required: true,
      type: 'string',
      min: 8,
      message: t('app.rules.password.minLength'),
      trigger: 'blur',
    },
  ],
  newPassword: [
    {
      required: true,
      type: 'string',
      min: 8,
      message: t('app.rules.password.minLength'),
      trigger: 'blur',
    },
  ],
  reNewPassword: [
    {
      required: true,
      validator(__: any, _: any, callback: any) {
        if (form.value.newPassword !== form.value.reNewPassword) callback(new Error(t('app.rules.password.noMatch')))
        else if (form.value.reNewPassword.length === 0) callback(new Error(t('app.rules.password.minLength')))
        else callback()
      },
      trigger: 'blur',
    },
  ],
}

async function submit() {
  await formRef.value?.validate()
  await apiResetPassword(form.value)
  ElMessage.success(t('app.user.password.resetSuccess'))
}
</script>

<template>
  <div class="flex flex-col justify-center mx-auto px-36px" style="align-items: center">
    <div style="width: 40vw; align-items: start" class="text-start">
      <div style="width: 450px; background-color: white">
        <h1>{{ $t('app.user.password.title') }}</h1>
        <ElForm ref="formRef" :rules="rules" :model="form" class="content" label-position="top">
          <div style="margin-top: 40px">
            <!-- old pass -->
            <ElFormItem prop="password" :label="$t('app.user.password.old')">
              <ElInput maxlength="255" v-model="form.password" class="h-40px" type="password" />
            </ElFormItem>

            <!-- new pass -->
            <ElFormItem prop="newPassword" :label="$t('app.user.password.new')">
              <ElInput maxlength="255" v-model="form.newPassword" class="h-40px" type="password" />
            </ElFormItem>

            <!-- confirm pass -->
            <ElFormItem prop="reNewPassword" :label="$t('app.user.password.confirm')">
              <ElInput maxlength="255" v-model="form.reNewPassword" class="h-40px" type="password" />
            </ElFormItem>
          </div>

          <!-- submit -->
          <ElButton class="w-full mt-5" type="primary" @click="submit">
            {{ $t('app.user.password.update') }}
          </ElButton>

          <!-- forgot -->
          <div class="w-full mt-5 text-center">
            <router-link :to="FORGETPASS_PATH" class="mx-1 text-blue-600">
              {{ $t('app.user.password.forgot') }}
            </router-link>
          </div>
        </ElForm>
      </div>
    </div>
  </div>
</template>

<style scoped>
h1 {
  font-size: 30px;
}

:deep(.el-select .el-input) {
  height: 40px;
}

:deep(.el-button) {
  height: 40px;
}

.row {
  margin-top: 1em;
  margin-bottom: 1em;
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.left,
.right {
  display: flex;
  align-items: center;
}

.left {
  justify-content: flex-start;
  /* flex-grow: 1; */
}

.right {
  /* justify-content: flex-end; */
  flex-grow: 1;
}

.content {
  margin-top: 40px;
}

/* el-upload */
:deep(.content .el-upload) {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

/* el-image */
.content .block {
  padding: 30px 0;
  text-align: center;
  border-right: solid 1px var(--el-border-color);
  display: inline-block;
  width: 49%;
  box-sizing: border-box;
  vertical-align: top;
}
.content .demonstration {
  display: block;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 20px;
}
.content .el-image {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

.content .image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-secondary);
  font-size: 30px;
}
</style>
