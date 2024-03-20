<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { ElInput } from 'element-plus'
import { useI18n } from 'vue-i18n'
import type { UseCollapse } from '@/components/collapse/useCollapse'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import IconSvg from '@/components/IconSvg.vue'
import useApi from '@/hooks/useApi'
import { isUrlRule, notNullRule } from '@/commons'
import { apiUpdateStroageCF } from '@/api/system'

const props = defineProps<{ collapse: UseCollapse; name: string; config: Partial<SystemAPI.StorageCF> }>()
const { t } = useI18n()
const tBase = 'app.system.storage.cf'
const formRef = ref<FormInstance>()
const rules: FormRules<typeof props.config> = {
  accountID: notNullRule(t(`${tBase}.rules.accountID`)),
  accessKeyID: notNullRule(t(`${tBase}.rules.accessKeyID`)),
  accessKeySecret: notNullRule(t(`${tBase}.rules.accessKeySecret`)),
  bucketName: notNullRule(t(`${tBase}.rules.bucketName`)),
  bucketUrl: [...isUrlRule(t(`${tBase}.rules.bucketUrl`)), ...notNullRule(t(`${tBase}.rules.bucketUrl`))],
}

const [submitting, update] = useApi(apiUpdateStroageCF)
function submit() {
  formRef.value!.validate((valid) => {
    if (valid)
      update(props.config as SystemAPI.StorageCF)
  })
}
</script>

<template>
  <CollapseCardItem :name="name" :collapse-ctx="collapse">
    <template #title>
      <div class="row-lr">
        <div class="left mr-8px">
          <IconSvg name="ac-cloudflare" />
        </div>
        <div class="right font-bold">
          {{ $t(`${tBase}.title`) }}
        </div>
      </div>
    </template>
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="props.config" @submit.prevent="submit">
      <ElFormItem prop="accountID" :label="$t(`${tBase}.accountID`)">
        <ElInput v-model="props.config.accountID" maxlength="255" />
      </ElFormItem>
      <ElFormItem prop="accessKeyID" :label="$t(`${tBase}.accessKeyID`)">
        <ElInput v-model="props.config.accessKeyID" maxlength="255" />
      </ElFormItem>
      <ElFormItem prop="accessKeySecret" :label="$t(`${tBase}.accessKeySecret`)">
        <ElInput v-model="props.config.accessKeySecret" maxlength="255" />
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
