<template>
  <div v-if="param.header && param.header.length">
    <p class="mb-5px text-16px">Header</p>
    <SimpleParameterEditor :readonly="true" v-model="param.header" allow-mock />
  </div>

  <div v-for="(_, contentTypeKey) in param.content" :key="contentTypeKey">
    <p class="mb-10px text-16px mt-14px">Content({{ contentTypeKey }})</p>
    <JSONSchemaEditor :readonly="true" :definitions="definitions" v-model="param.content[contentTypeKey].schema" v-if="isJsonschema" />
    <p class="my-10px" v-if="param.content[contentTypeKey].schema.example">
      {{ $t('app.response.tips.responseExample') }}
    </p>
    <CodeEditor
      :readonly="true"
      v-if="param.content[contentTypeKey].schema.example"
      class="mt-14px"
      v-model="param.content[contentTypeKey].schema.example"
      :lang="contentTypes[contentTypeKey]"
    />
    <ResponseExamplesForm :examples="examples" readonly :lang="contentTypes[contentTypeKey]" />
  </div>
</template>

<script lang="ts">
export declare interface APICatResponse {
  id?: number | string
  code: number
  description: string
  content: Record<string, { schema: JSONSchema }>
  header?: APICatSchemaObject[]
}

export const contentTypes: Record<string, string> = {
  'application/json': 'json',
  'application/xml': 'xml',
  'text/html': 'html',
  'text/plain': 'raw',
  'application/octet-stream': 'raw',
}
</script>

<script setup lang="ts">
import SimpleParameterEditor from '@/components/APIEditor/SimpleEditor.vue'
import JSONSchemaEditor from '@/components/APIEditor/Editor.vue'
import CodeEditor from '@/components/APIEditor/CodeEditor.vue'
import { APICatSchemaObject, DefinitionSchema, JSONSchema } from './APIEditor/types'
import ResponseExamplesForm from '@/views/component/ResponseExamples.vue'
import { isEmpty } from 'lodash-es'

const props = defineProps<{
  param: APICatResponse | any
  definitions: DefinitionSchema[]
}>()

const contentTypes: Record<string, string> = {
  'application/json': 'json',
  'application/xml': 'xml',
  'text/html': 'html',
  'text/plain': 'raw',
  'application/octet-stream': 'raw',
}

const contentDefaultType = computed(() => {
  for (let x in props.param.content) {
    return x
  }
  return 'application/json'
})

const isJsonschema = computed(() => contentDefaultType.value == 'application/json' || contentDefaultType.value == 'application/xml')
const examples = computed(()=>props.param.content[contentDefaultType.value].examples || {})
</script>
