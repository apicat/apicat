<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { apiGetService, apiUpdateService } from '@/api/system'
import useApi from '@/hooks/useApi'
import { isUrlRule, notNullRule } from '@/commons'

const { t } = useI18n()
const tBase = 'app.system.service'
const form = ref<SystemAPI.ServiceData>({
  appUrl: '',
  mockUrl: '',
})
const formRef = ref<FormInstance>()
const rules = reactive<FormRules<SystemAPI.ServiceData>>({
  appUrl: [...isUrlRule(t(`${tBase}.rules.appurl`), true, true), ...notNullRule(t(`${tBase}.rules.appurl`))],
  mockUrl: [...isUrlRule(t(`${tBase}.rules.mockurl`), true, true), ...notNullRule(t(`${tBase}.rules.mockurl`))],
})

const [submitting, update] = useApi(apiUpdateService)
async function submit() {
  try {
    await formRef.value!.validate()
    await update(form.value)
  }
  catch (e) {}
}

apiGetService().then((v) => {
  form.value = v
})
</script>

<template>
  <div class="bg-white w-85%">
    <h1>{{ $t('app.system.service.title') }}</h1>
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="form" size="large" @submit.prevent="submit">
      <div>
        <!-- appurl -->
        <ElFormItem prop="appUrl" :label="$t('app.system.service.appurl')">
          <ElInput v-model="form.appUrl" maxlength="255" />
        </ElFormItem>

        <!-- mockurl -->
        <ElFormItem prop="mockUrl" :label="$t('app.system.service.mockurl')">
          <ElInput v-model="form.mockUrl" maxlength="255" />
        </ElFormItem>
      </div>

      <!-- submit -->
      <ElButton :loading="submitting" class="w-full mt-8px" type="primary" @click="submit">
        {{ $t('app.common.save') }}
      </ElButton>
    </ElForm>
  </div>
</template>

<style scoped>
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
