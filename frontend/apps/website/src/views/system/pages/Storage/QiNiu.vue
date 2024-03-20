<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { ElInput } from 'element-plus'
import { useI18n } from 'vue-i18n'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import { type UseCollapse } from '@/components/collapse/useCollapse'
import IconSvg from '@/components/IconSvg.vue'
import { isUrlRule, notNullRule } from '@/commons'
import useApi from '@/hooks/useApi'
import { apiUpdateStroageQiniu } from '@/api/system'

const props = defineProps<{ collapse: UseCollapse; name: string; config: Partial<SystemAPI.StorageQiniu> }>()
const { t } = useI18n()
const tBase = 'app.system.storage.qiniu'
const formRef = ref<FormInstance>()
const rules: FormRules<typeof props.config> = {
  accessKey: notNullRule(t(`${tBase}.rules.accessKey`)),
  secretKey: notNullRule(t(`${tBase}.rules.secretKey`)),
  bucketName: notNullRule(t(`${tBase}.rules.bucketName`)),
  bucketUrl: [...isUrlRule(t(`${tBase}.rules.bucketUrl`)), ...notNullRule(t(`${tBase}.rules.bucketUrl`))],
}

const [submitting, update] = useApi(apiUpdateStroageQiniu)
function submit() {
  formRef.value!.validate((valid) => {
    if (valid)
      update(props.config as SystemAPI.StorageQiniu)
  })
}
</script>

<template>
  <CollapseCardItem :name="name" :collapse-ctx="collapse">
    <template #title>
      <div class="row-lr">
        <div class="left mr-8px">
          <IconSvg name="ac-qiniu" width="60" height="24" />
        </div>
        <div class="right font-bold">
          {{ $t(`${tBase}.title`) }}
        </div>
      </div>
    </template>
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="props.config" @submit.prevent="submit">
      <ElFormItem prop="accessKey" :label="$t(`${tBase}.accessKey`)">
        <ElInput v-model="props.config.accessKey" maxlength="255" />
      </ElFormItem>
      <ElFormItem prop="secretKey" :label="$t(`${tBase}.secretKey`)">
        <ElInput v-model="props.config.secretKey" maxlength="255" />
      </ElFormItem>
      <ElFormItem prop="bucketName" :label="$t(`${tBase}.bucketName`)">
        <ElInput v-model="props.config.bucketName" maxlength="255" />
      </ElFormItem>
      <ElFormItem prop="bucketUrl" :label="$t(`${tBase}.bucketUrl`)">
        <ElInput v-model="props.config.bucketUrl" maxlength="255" />
      </ElFormItem>
    </ElForm>

    <el-button :loading="submitting" type="primary" @click="submit">
      {{ $t('app.common.update') }}
    </el-button>
  </CollapseCardItem>
</template>

<style scoped></style>
