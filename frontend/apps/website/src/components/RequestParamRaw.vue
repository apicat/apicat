<template>
  <ToggleHeading title="请求参数" v-if="headers.length || cookies.length || queries.length || paths.length || hasBody">
    <ToggleHeading title="Header" v-if="headers.length">
      <SimpleParameterEditor :readonly="true" v-model="headers" />
    </ToggleHeading>
    <ToggleHeading title="Cookie" v-if="cookies.length">
      <SimpleParameterEditor :readonly="true" v-model="cookies" />
    </ToggleHeading>
    <ToggleHeading title="Query" v-if="queries.length">
      <SimpleParameterEditor :readonly="true" v-model="queries" />
    </ToggleHeading>
    <ToggleHeading title="Path" v-if="paths.length">
      <SimpleParameterEditor :readonly="true" v-model="paths" />
    </ToggleHeading>
    <ToggleHeading :title="bodyTitle" v-if="hasBody">
      <div v-if="contentType === RequestContentTypesMap['form-data'] || contentType === RequestContentTypesMap['x-www-form-urlencoded']">
        <SimpleSchemaEditor :readonly="true" v-model="contentVal.schema" />
      </div>

      <div v-if="isJsonSchema">
        <JSONSchemaEditor :readonly="true" v-model="contentVal.schema" :definitions="definitions" />
      </div>

      <div v-if="contentType === RequestContentTypesMap.raw">
        <CodeEditor v-model="contentVal.schema.example" />
      </div>

      <div v-if="contentType === RequestContentTypesMap.binary">
        <FileUploaderWrapper :readonly="true" class="flex items-center justify-center bg-gray-100 border border-gray-200 border-solid rounded cursor-pointer h-30px">
          <span class="truncate px-6px"> {{ contentVal.schema.example }}</span>
        </FileUploaderWrapper>
      </div>
    </ToggleHeading>
  </ToggleHeading>
</template>
<script setup lang="ts">
import SimpleParameterEditor from '@/components/APIEditor/SimpleEditor.vue'
import SimpleSchemaEditor from '@/components/APIEditor/SimpleSchemaEditor.vue'
import JSONSchemaEditor from '@/components/APIEditor/Editor.vue'
import CodeEditor from '@/components/APIEditor/CodeEditor.vue'
import { RequestContentTypesMap } from '@/commons'
import { useNodeAttrs, HTTP_REQUEST_NODE_KEY } from '@/hooks/useNodeAttrs'
import { Definition } from './APIEditor/types'
import { HttpDocument } from '@/typings'

const props = defineProps<{ doc: HttpDocument; definitions: Definition[] }>()
const request = useNodeAttrs(props, HTTP_REQUEST_NODE_KEY, 'doc')

const headers = computed(() => request.value.parameters.header || [])
const cookies = computed(() => request.value.parameters.cookie || [])
const queries = computed(() => request.value.parameters.query || [])
const paths = computed(() => request.value.parameters.path || [])

const contentType = computed(() => {
  if (!request.value.content) {
    return RequestContentTypesMap.none
  }
  const keys = Object.keys(request.value.content)
  if (keys && keys.length) {
    return keys[0]
  }
  return RequestContentTypesMap.none
})

const isJsonSchema = computed(() => {
  return contentType.value == RequestContentTypesMap.json || contentType.value == RequestContentTypesMap.xml
})
const bodyTitle = computed(() => `Body (${contentType.value})`)
const hasBody = computed(() => contentType.value && contentType.value !== RequestContentTypesMap.none)
const contentVal = computed(() => request.value.content[contentType.value])
</script>
