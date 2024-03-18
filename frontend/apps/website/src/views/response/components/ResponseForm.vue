<script setup lang="ts">
import { JSONSchemaEditor, SimpleParamTable, ToggleHeading } from '@apicat/components'
import ResponseExamples from './ResponseExamples.vue'
import { useResponse } from './useResponse'
import { ResponseContentTypesMap } from '@/commons/constant'

export interface ResponseFormProps {
  response: Definition.ResponseDetail
  definitionSchemas: Definition.Schema[]
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
        <JSONSchemaEditor v-if="isJSONSchema" v-model:schema="contentSchema" :definition-schemas="definitionSchemas" />
        <ResponseExamples v-model:examples="examples" :lang="(ResponseContentTypesMap as any)[contentType]" />
      </div>
    </ToggleHeading>
  </div>
</template>
