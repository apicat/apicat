<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

import { useI18n } from 'vue-i18n'
import { apiUpdateEmail } from '@/api/user'
import useApi from '@/hooks/useApi'

const { t } = useI18n()

const form = ref<{ email: string }>({
  email: '',
})
const formRef = ref<FormInstance>()
const rules = reactive<FormRules>({
  email: [
    {
      validator(rule: any, value: any, callback: any) {
        function test() {
          return /^[A-Z0-9._+-]+@[A-Z0-9][A-Z0-9.-]*\.[A-Z]{2,63}$/i.test(value)
        }
        // if (value && !isEmail(value))
        if (value && !test()) return callback(new Error(t('app.rules.email.correct')))

        return callback()
      },
      trigger: 'blur',
    },
    {
      required: true,
      message: t('app.rules.email.required'),
      trigger: 'blur',
    },
  ],
})

const [submitting, updateEmail] = useApi(apiUpdateEmail)
async function submit() {
  try {
    await formRef.value!.validate()
    await updateEmail(form.value.email)
    ElMessage.success(t('app.user.email.success'))
  } catch (e) {}
}
</script>

<template>
  <div class="flex flex-col justify-center mx-auto px-36px" style="align-items: center">
    <div style="width: 40vw; align-items: start" class="text-start">
      <div style="width: 450px; background-color: white">
        <h1>{{ $t('app.user.email.title') }}</h1>
        <ElForm
          ref="formRef"
          class="content"
          label-position="top"
          :rules="rules"
          :model="form"
          @submit.prevent="submit">
          <div style="margin-top: 40px">
            <!-- username -->
            <ElFormItem prop="email" :label="$t('app.user.email.email')">
              <ElInput maxlength="255" v-model="form.email" class="h-40px" />
            </ElFormItem>
          </div>

          <!-- submit -->
          <ElButton :loading="submitting" class="w-full" type="primary" @click="submit">
            {{ $t('app.user.email.send') }}
          </ElButton>
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
