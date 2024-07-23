<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import { type UseCollapse } from '@/components/collapse/useCollapse'
import IconSvg from '@/components/IconSvg.vue'
import { apiUpdateModelAzure } from '@/api/system'
import useApi from '@/hooks/useApi'
import { notNullRule } from '@/commons'

const props = defineProps<{ collapse: UseCollapse; name: string; config: Partial<SystemAPI.ModelAzure> }>()
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
  embedding: '',
  ...props.config
})

// sync config
watch(() => props.config, (val) => Object.assign(config.value, val))

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
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="config" @submit.prevent="submit">
      <!-- api key -->
      <ElFormItem prop="apiKey" :label="$t(`${tBase}.apiKey`)">
        <ElInput v-model="config.apiKey" maxlength="255" />
      </ElFormItem>

      <!-- endpoint  -->
      <ElFormItem prop="endpoint" :label="$t(`${tBase}.endpoint`)">
        <ElInput v-model="config.endpoint" maxlength="255" />
      </ElFormItem>

      <!-- llm name  -->
      <ElFormItem prop="llm" :label="$t(`${tBase}.llmName`)">
        <ElInput v-model="config.llm" maxlength="255" />
      </ElFormItem>

      <!-- embedding  -->
      <ElFormItem prop="embedding" :label="$t(`${tBase}.embedding`)">
        <ElInput v-model="config.embedding" maxlength="255" />
      </ElFormItem>
    </ElForm>

    <el-button :loading="submitting" type="primary" @click="submit">
      {{ $t('app.common.save') }}
    </el-button>
  </CollapseCardItem>
</template>
