<script setup lang="ts">
import type { SchemaTreeNode } from '@apicat/components'
import { JSONSchemaTable, SimpleParamTable, ToggleHeading } from '@apicat/components'
import type { JSONSchema } from '@apicat/editor'
import ResponseExamples from './ResponseExamples.vue'
import { useResponse } from './useResponse'
import { ResponseContentTypesMap } from '@/commons/constant'
import { apiParseSchema } from '@/api/project/definition/schema'

export interface ResponseFormProps {
  response: Definition.ResponseDetail
  definitionSchemas: Definition.Schema[]
  handleIntelligentSchema?: (josnschema: JSONSchema, node: SchemaTreeNode) => Promise<{ nid: string, schema: JSONSchema } | void>
}

export interface ResponseFormEmits {
  (e: 'update:response', response: Definition.ResponseDetail): void
}

const props = withDefaults(defineProps<ResponseFormProps>(), {
  response: () => ({}) as Definition.ResponseDetail,
  definitionSchemas: () => [],
})
const emits = defineEmits<ResponseFormEmits>()

const { examples, headers, contentType, isJSONSchema, contentSchema } = useResponse(props, emits)
</script>

<template>
  <div>
    <ToggleHeading title="Header" :expand="!!headers.length">
      <SimpleParamTable v-model:datas="headers" class="mt-10px" allow-mock />
    </ToggleHeading>

    <ToggleHeading title="Content">
      <template #extra>
        <el-select v-model="contentType">
          <template #prefix>
            Content-Type
          </template>
          <el-option v-for="(key, value) in ResponseContentTypesMap" :key="key" :label="value" :value="value" />
        </el-select>
      </template>

      <div class="">
        <JSONSchemaTable
          v-if="isJSONSchema"
          v-model:schema="contentSchema"
          :definition-schemas="definitionSchemas"
          :handle-parse-schema="apiParseSchema"
          :handle-intelligent-schema="handleIntelligentSchema"
        />
        <ResponseExamples v-model:examples="examples" :lang="(ResponseContentTypesMap as any)[contentType]" />
      </div>
    </ToggleHeading>
  </div>
</template>
