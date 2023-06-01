<template>
  <ToggleHeading :title="$t('app.request.title')" v-if="hasHeader || hasCookie || hasQuery || hasPath || hasBody">
    <ToggleHeading title="Header" v-if="hasHeader">
      <SimpleParameterEditor readonly v-model="headers">
        <template #before>
          <GlobalParameter :data="globalHeaders" readonly />
        </template>
      </SimpleParameterEditor>
    </ToggleHeading>
    <ToggleHeading title="Cookie" v-if="hasCookie">
      <SimpleParameterEditor readonly v-model="cookies">
        <template #before>
          <GlobalParameter :data="globalCookies" readonly />
        </template>
      </SimpleParameterEditor>
    </ToggleHeading>
    <ToggleHeading title="Query" v-if="hasQuery">
      <SimpleParameterEditor readonly v-model="queries">
        <template #before>
          <GlobalParameter :data="globalQueries" readonly />
        </template>
      </SimpleParameterEditor>
    </ToggleHeading>
    <ToggleHeading title="Path" v-if="hasPath">
      <SimpleParameterEditor readonly v-model="paths">
        <template #before>
          <GlobalParameter :data="globalPaths" readonly />
        </template>
      </SimpleParameterEditor>
    </ToggleHeading>
    <ToggleHeading :title="bodyTitle" v-if="hasBody">
      <div v-if="contentType === RequestContentTypesMap['form-data'] || contentType === RequestContentTypesMap['x-www-form-urlencoded']">
        <SimpleSchemaEditor readonly v-model="contentVal.schema" />
      </div>

      <div v-if="isJsonSchema">
        <JSONSchemaEditor readonly v-model="contentVal.schema" :definitions="definitions" />
      </div>

      <div v-if="contentType === RequestContentTypesMap.raw">
        <CodeEditor v-model="contentVal.schema.example" />
      </div>

      <div v-if="contentType === RequestContentTypesMap.binary">
        <FileUploaderWrapper readonly class="flex items-center justify-center bg-gray-100 border border-gray-200 border-solid rounded cursor-pointer h-30px">
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
import GlobalParameter from './RequestParamEditor/GlobalParameter.vue'
import { RequestContentTypesMap } from '@/commons'
import { useNodeAttrs, HTTP_REQUEST_NODE_KEY } from '@/hooks/useNodeAttrs'
import { DefinitionSchema } from './APIEditor/types'
import { HttpDocument } from '@/typings'
import uesGlobalParametersStore from '@/store/globalParameters'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps<{ doc: HttpDocument; definitions: DefinitionSchema[] }>()
const request = useNodeAttrs(props, HTTP_REQUEST_NODE_KEY, 'doc')
const globalParametersStore = uesGlobalParametersStore()
const { parameters: globalParameters } = storeToRefs(globalParametersStore)

const headers = computed(() => request.value.parameters.header || [])
const cookies = computed(() => request.value.parameters.cookie || [])
const queries = computed(() => request.value.parameters.query || [])
const paths = computed(() => request.value.parameters.path || [])

const globalHeaders = computed(() =>
  globalParameters.value.header
    .filter((param) => !(request.value.globalExcepts.header || []).includes(param.id))
    .map((param) => ({ ...param, required: param.required ? t('app.table.yes') : t('app.table.no') }))
)

const globalCookies = computed(() =>
  globalParameters.value.cookie
    .filter((param) => !(request.value.globalExcepts.cookie || []).includes(param.id))
    .map((param) => ({ ...param, required: param.required ? t('app.table.yes') : t('app.table.no') }))
)

const globalQueries = computed(() =>
  globalParameters.value.query
    .filter((param) => !(request.value.globalExcepts.query || []).includes(param.id))
    .map((param) => ({ ...param, required: param.required ? t('app.table.yes') : t('app.table.no') }))
)

const globalPaths = computed(() =>
  globalParameters.value.path
    .filter((param) => !(request.value.globalExcepts.path || []).includes(param.id))
    .map((param) => ({ ...param, required: param.required ? t('app.table.yes') : t('app.table.no') }))
)

// todo 公参数暂未实现，{$ref:'#definition/parameters/xxx'}
const hasHeader = computed(() => (headers.value || []).filter((item: any) => !item.$ref).length || globalHeaders.value.length)
const hasCookie = computed(() => (cookies.value || []).filter((item: any) => !item.$ref).length || globalCookies.value.length)
const hasQuery = computed(() => (queries.value || []).filter((item: any) => !item.$ref).length || globalQueries.value.length)
const hasPath = computed(() => (paths.value || []).filter((item: any) => !item.$ref).length || globalPaths.value.length)

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

const isJsonSchema = computed(() => contentType.value == RequestContentTypesMap.json || contentType.value == RequestContentTypesMap.xml)
const bodyTitle = computed(() => `Body (${contentType.value})`)
const hasBody = computed(() => contentType.value && contentType.value !== RequestContentTypesMap.none)
const contentVal = computed(() => request.value.content[contentType.value])
</script>
