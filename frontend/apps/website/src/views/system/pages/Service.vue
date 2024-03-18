<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { apiGetService, apiUpdateService } from '@/api/system'
import useApi from '@/hooks/useApi'
import { isUrlRule, notNullRule } from '@/commons'

const { t } = useI18n()
const tBase = 'app.system.service'
const form = ref<SystemAPI.ServiceData>({
  appName: '',
  appUrl: '',
  appServerBind: '',
  mockUrl: '',
  mockServerBind: '',
})
const formRef = ref<FormInstance>()
const rules = reactive<FormRules<SystemAPI.ServiceData>>({
  appName: notNullRule(t(`${tBase}.rules.appname`)),
  appUrl: [...isUrlRule(t(`${tBase}.rules.appurl`), true, true), ...notNullRule(t(`${tBase}.rules.appurl`))],
  appServerBind: [
    ...isUrlRule(t(`${tBase}.rules.appserver`), true, false),
    ...notNullRule(t(`${tBase}.rules.appserver`)),
  ],
  mockUrl: [...isUrlRule(t(`${tBase}.rules.mockurl`), true, true), ...notNullRule(t(`${tBase}.rules.mockurl`))],
  mockServerBind: [
    ...isUrlRule(t(`${tBase}.rules.mockserver`), true, false),
    ...notNullRule(t(`${tBase}.rules.mockserver`)),
  ],
})

const [submitting, update] = useApi(apiUpdateService)
async function submit() {
  try {
    await formRef.value!.validate()
    await update(form.value)
  } catch (e) {}
}

apiGetService().then((v) => {
  form.value = v
})
</script>

<template>
  <div class="bg-white w-450px">
    <h1>{{ $t('app.system.service.title') }}</h1>
    <ElForm ref="formRef" class="content" label-position="top" :rules="rules" :model="form" @submit.prevent="submit">
      <div style="margin-top: 40px">
        <!-- appname -->
        <ElFormItem prop="appName" :label="$t('app.system.service.appname')">
          <ElInput v-model="form.appName" maxlength="255" />
        </ElFormItem>

        <!-- appurl -->
        <ElFormItem prop="appUrl" :label="$t('app.system.service.appurl')">
          <ElInput v-model="form.appUrl" maxlength="255" />
        </ElFormItem>

        <!-- server -->
        <ElFormItem prop="appServerBind" :label="$t('app.system.service.server')">
          <ElInput v-model="form.appServerBind" maxlength="255" />
        </ElFormItem>

        <!-- mockurl -->
        <ElFormItem prop="mockUrl" :label="$t('app.system.service.mockurl')">
          <ElInput v-model="form.mockUrl" maxlength="255" />
        </ElFormItem>

        <!-- mockserver -->
        <ElFormItem prop="mockServerBind" :label="$t('app.system.service.mockserver')">
          <ElInput v-model="form.mockServerBind" maxlength="255" />
        </ElFormItem>
      </div>

      <!-- submit -->
      <ElButton :loading="submitting" class="w-full" type="primary" @click="submit">
        {{ $t('app.common.update') }}
      </ElButton>
    </ElForm>
    <p class="text-gray-helper mt-10px">
      {{ $t('app.system.service.tip') }}
    </p>
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
