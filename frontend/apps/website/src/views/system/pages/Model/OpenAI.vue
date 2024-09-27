<script setup lang="ts">
import { ElSelect, type FormInstance, type FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import type { UseCollapse } from '@/components/collapse/useCollapse'
import IconSvg from '@/components/IconSvg.vue'
import { apiUpdateModelOpenAI } from '@/api/system'
import useApi from '@/hooks/useApi'
import { notNullRule } from '@/commons'

const props = defineProps<{
  collapse: UseCollapse
  name: string
  config: Partial<SystemAPI.ModelOpenAI>
  llmModels?: string[]
  embeddingModels?: string[]
}>()

const { t } = useI18n()
const tBase = 'app.system.model.openai'
const formRef = ref<FormInstance>()
const rules: FormRules<typeof props.config> = {
  apiKey: notNullRule(t(`${tBase}.rules.apiKey`)),
  llm: notNullRule(t(`${tBase}.rules.llmName`)),
  embedding: notNullRule(t(`${tBase}.rules.embedding`)),
}
const [submitting, update] = useApi(apiUpdateModelOpenAI)
function submit() {
  formRef.value!.validate((valid) => {
    if (valid)
      update(props.config as SystemAPI.ModelOpenAI)
  })
}
</script>

<template>
  <CollapseCardItem :name="name" :collapse-ctx="collapse">
    <template #title>
      <div class="row-lr">
        <div class="left mr-8px">
          <IconSvg name="ac-openai" width="24" />
        </div>
        <div class="font-bold right">
          {{ $t(`${tBase}.title`) }}
        </div>
      </div>
    </template>
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="props.config" size="large" @submit.prevent="submit">
      <!-- api key -->
      <ElFormItem prop="apiKey" :label="$t(`${tBase}.apiKey`)">
        <ElInput v-model="props.config.apiKey" maxlength="255" />
      </ElFormItem>

      <!-- form email -->
      <ElFormItem prop="organizationID" :label="$t(`${tBase}.organizationID`)">
        <ElInput v-model="props.config.organizationID" maxlength="255" />
      </ElFormItem>

      <!-- api user -->
      <ElFormItem prop="apiBase" :label="$t(`${tBase}.apiBase`)">
        <ElInput v-model="props.config.apiBase" maxlength="255" />
      </ElFormItem>

      <!-- llm name  -->
      <ElFormItem prop="llm" :label="$t(`${tBase}.llmName`)">
        <ElSelect v-model="props.config.llm" class="w-full">
          <ElOption v-for="i in llmModels" :key="i" :label="i" :value="i" />
        </ElSelect>
      </ElFormItem>

      <!-- embedding name  -->
      <ElFormItem prop="embedding" :label="$t(`${tBase}.embedding`)">
        <ElSelect v-model="props.config.embedding" class="w-full">
          <ElOption v-for="i in embeddingModels" :key="i" :label="i" :value="i" />
        </ElSelect>
      </ElFormItem>
    </ElForm>

    <el-button :loading="submitting" type="primary" @click="submit">
      {{ $t('app.common.save') }}
    </el-button>
  </CollapseCardItem>
</template>
