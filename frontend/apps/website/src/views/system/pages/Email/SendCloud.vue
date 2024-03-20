<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import { type UseCollapse } from '@/components/collapse/useCollapse'
import IconSvg from '@/components/IconSvg.vue'
import { apiUpdateEmailSendCloud } from '@/api/system'
import { notNullRule } from '@/commons'
import useApi from '@/hooks/useApi'

const props = defineProps<{ collapse: UseCollapse; name: string; config: Partial<SystemAPI.EmailSendCloud> }>()
const { t } = useI18n()
const tBase = 'app.system.email.sendcloud'
const formRef = ref<FormInstance>()
const rules: FormRules<typeof props.config> = {
  apiUser: notNullRule(t(`${tBase}.rules.apiUser`)),
  apiKey: notNullRule(t(`${tBase}.rules.apiKey`)),
  fromName: notNullRule(t(`${tBase}.rules.fromName`)),
  fromEmail: notNullRule(t(`${tBase}.rules.fromEmail`)),
}
const [submitting, update] = useApi(apiUpdateEmailSendCloud)
function submit() {
  formRef.value!.validate((valid) => {
    if (valid)
      update(props.config as SystemAPI.EmailSendCloud)
  })
}
</script>

<template>
  <CollapseCardItem :name="name" :collapse-ctx="collapse">
    <template #title>
      <div class="row-lr">
        <div class="left mr-8px">
          <IconSvg name="ac-SendCloud" width="24" />
        </div>
        <div class="font-bold right">
          {{ $t(`${tBase}.title`) }}
        </div>
      </div>
    </template>
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="props.config" @submit.prevent="submit">
      <!-- api user -->
      <ElFormItem prop="apiUser" :label="$t(`${tBase}.apiUser`)">
        <ElInput v-model="props.config.apiUser" maxlength="255" />
      </ElFormItem>

      <!-- api key -->
      <ElFormItem prop="apiKey" :label="$t(`${tBase}.apiKey`)">
        <ElInput v-model="props.config.apiKey" maxlength="255" />
      </ElFormItem>

      <!-- form email -->
      <ElFormItem prop="fromEmail" :label="$t(`${tBase}.fromEmail`)">
        <ElInput v-model="props.config.fromEmail" maxlength="255" />
      </ElFormItem>

      <!-- form name -->
      <ElFormItem prop="fromName" :label="$t(`${tBase}.fromName`)">
        <ElInput v-model="props.config.fromName" maxlength="255" />
      </ElFormItem>
    </ElForm>

    <el-button :loading="submitting" type="primary" @click="submit">
      {{ $t('app.common.update') }}
    </el-button>
  </CollapseCardItem>
</template>

<style scoped></style>
