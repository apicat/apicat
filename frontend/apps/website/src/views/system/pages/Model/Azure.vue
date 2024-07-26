<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import type { UseCollapse } from '@/components/collapse/useCollapse'
import IconSvg from '@/components/IconSvg.vue'
import { apiUpdateModelAzure } from '@/api/system'
import useApi from '@/hooks/useApi'
import { notNullRule } from '@/commons'

const props = defineProps<{
  collapse: UseCollapse
  name: string
  config: Partial<SystemAPI.ModelAzure>
  llmModels?: string[]
  embeddingModels?: string[]
}>()
const { t } = useI18n()
const tBase = 'app.system.model.azure'
const formRef = ref<FormInstance>()
const rules: FormRules<typeof props.config> = {
  apiKey: notNullRule(t(`${tBase}.rules.apiKey`)),
  endpoint: notNullRule(t(`${tBase}.rules.endpoint`)),
  llm: notNullRule(t(`${tBase}.rules.llmName`)),
  embedding: notNullRule(t(`${tBase}.rules.embedding`)),
}

const [submitting, update] = useApi(apiUpdateModelAzure)

const config = ref({
  apiKey: '',
  endpoint: '',
  llm: '',
  llmDeployName: '',
  embedding: '',
  embeddingDeployName: '',
  ...props.config,
})

// sync config
watch(() => props.config, val => Object.assign(config.value, val))

function submit() {
  formRef.value!.validate((valid) => {
    if (valid)
      update(config.value as SystemAPI.ModelAzure)
  })
}
</script>

<template>
  <CollapseCardItem :name="name" :collapse-ctx="collapse">
    <template #title>
      <div class="row-lr">
        <div class="left mr-8px">
          <IconSvg name="ac-azure" width="24" />
        </div>
        <div class="font-bold right">
          {{ $t(`${tBase}.title`) }}
        </div>
      </div>
    </template>
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="config" size="large" @submit.prevent="submit">
      <!-- api key -->
      <ElFormItem prop="apiKey" :label="$t(`${tBase}.apiKey`)">
        <ElInput v-model="config.apiKey" maxlength="255" />
      </ElFormItem>

      <!-- endpoint  -->
      <ElFormItem prop="endpoint" :label="$t(`${tBase}.endpoint`)">
        <ElInput v-model="config.endpoint" maxlength="255" />
      </ElFormItem>

      <!-- llm name  -->
      <el-form-item label="LLM deployment name">
        <el-col :span="11">
          <el-form-item prop="llmDeployName">
            <ElInput v-model="config.llmDeployName" maxlength="255" />
          </el-form-item>
        </el-col>
        <el-col class="text-center" :span="2">
          <span class="text-gray-500" />
        </el-col>
        <el-col :span="11">
          <el-form-item prop="llm">
            <ElSelect v-model="config.llm" class="w-full">
              <ElOption v-for="i in llmModels" :key="i" :label="i" :value="i" />
            </ElSelect>
          </el-form-item>
        </el-col>
      </el-form-item>

      <!-- embedding  -->
      <el-form-item label="Embedding model deployment name">
        <el-col :span="11">
          <el-form-item prop="embeddingDeployName">
            <ElInput v-model="config.embeddingDeployName" maxlength="255" />
          </el-form-item>
        </el-col>
        <el-col class="text-center" :span="2">
          <span class="text-gray-500" />
        </el-col>
        <el-col :span="11">
          <el-form-item prop="llm">
            <ElSelect v-model="config.embedding" class="w-full">
              <ElOption v-for="i in embeddingModels" :key="i" :label="i" :value="i" />
            </ElSelect>
          </el-form-item>
        </el-col>
      </el-form-item>
    </ElForm>

    <el-button :loading="submitting" type="primary" @click="submit">
      {{ $t('app.common.save') }}
    </el-button>
  </CollapseCardItem>
</template>
