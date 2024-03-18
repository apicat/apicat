<script setup lang="ts">
import { JSONSchemaEditor, SimpleParamTable, ToggleHeading } from '@apicat/components'
import ResponseExamples from './ResponseExamples.vue'
import { useResponse } from './useResponse'
import { ResponseContentTypesMap } from '@/commons/constant'

export interface ResponseFormProps {
  response: Definition.ResponseDetail
  definitionSchemas: Definition.Schema[]
}

const props = withDefaults(defineProps<ResponseFormProps>(), {
  isRef: false,
  definitionResponses: () => [],
  definitionSchemas: () => [],
})
const { examples, headers, contentType, isJSONSchema, contentSchema } = useResponse(props)
</script>

<template>
  <div>
    <p v-if="headers.length" class="mb-5px text-16px mt-10px">
      Header
    </p>
    <SimpleParamTable v-if="headers.length" :datas="headers" allow-mock readonly />
    <p class="mb-10px text-16px mt-14px">
      Content({{ contentType }})
    </p>
    <JSONSchemaEditor v-if="isJSONSchema" :schema="contentSchema" :definition-schemas="definitionSchemas" readonly />
    <ResponseExamples :examples="examples" readonly :lang="(ResponseContentTypesMap as any)[contentType]" />
  </div>
</template>
