<template>
  <el-space direction="vertical" fill warp class="w-full">
    <el-form :inline="true">
      <el-form-item :label="$t('app.response.table.code')">
        <el-select v-model="model.data.code" :placeholder="$t('app.response.table.code')" filterable>
          <el-option v-for="code in HttpCodeList" :label="code.code + ' ' + code.desc" :value="code.code" />
        </el-select>
      </el-form-item>

      <el-form-item :label="$t('app.response.table.name')">
        <el-input disabled :model-value="model.responseRefObject.name" maxlength="50" />
      </el-form-item>

      <el-form-item :label="$t('app.response.table.desc')">
        <el-input disabled :model-value="model.responseRefObject.description" maxlength="300" />
      </el-form-item>
    </el-form>
  </el-space>
  <DefinitionResponseRaw :response="model.responseRefObject" :definition-schemas="definitionSchemas" />
</template>

<script setup lang="ts">
import DefinitionResponseRaw from '@/components/DefinitionResponse/DefinitionResponseRaw.vue'
import { HttpCodeList } from '@apicat/shared'
import { DefinitionResponse } from '@/typings'
import { DefinitionSchema } from './APIEditor/types'

const props = defineProps<{
  response: any
  // 引用模型的集合
  definitionResponses?: DefinitionResponse[]
  definitionSchemas?: DefinitionSchema[]
}>()

const model = useVModel(props, 'response', undefined, { passive: true })
</script>
