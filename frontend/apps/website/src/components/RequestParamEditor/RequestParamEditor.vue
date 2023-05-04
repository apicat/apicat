<template>
  <div>
    <h2 class="text-16px font-500">{{ $t('app.request.title') }}</h2>
    <el-tabs>
      <el-tab-pane label="Header">
        <template #label>
          Header
          <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px" v-if="headersCount">{{ headersCount }}</span>
        </template>
        <SimpleParameterEditor v-model="headers">
          <template #before>
            <GlobalParameter :data="globalHeaders" :onSwitch="switchGlobalHeader" />
          </template>
        </SimpleParameterEditor>
      </el-tab-pane>

      <el-tab-pane label="Cookie">
        <template #label>
          Cookie
          <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px" v-if="cookiesCount">{{ cookiesCount }}</span>
        </template>
        <SimpleParameterEditor v-model="cookies">
          <template #before>
            <GlobalParameter :data="globalCookies" :onSwitch="switchGlobalCookie" />
          </template>
        </SimpleParameterEditor>
      </el-tab-pane>

      <el-tab-pane label="Query">
        <template #label>
          Query
          <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px" v-if="queriesCount">{{ queriesCount }}</span>
        </template>
        <SimpleParameterEditor v-model="queries">
          <template #before>
            <GlobalParameter :data="globalQueries" :onSwitch="switchGlobalQuery" />
          </template>
        </SimpleParameterEditor>
      </el-tab-pane>

      <el-tab-pane label="Body">
        <template #label>
          Body
          <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px" v-if="bodyCount">{{ bodyCount }}</span>
        </template>

        <div class="mb-16px">
          <el-radio-group v-model="currentContentTypeRef">
            <el-radio :label="value" v-for="(value, key) in RequestContentTypesMap">{{ key }}</el-radio>
          </el-radio-group>
        </div>

        <div v-if="currentContentTypeRef === RequestContentTypesMap.none" class="text-center">{{ $t('app.request.tips.noRequestBody') }}</div>

        <div v-if="currentContentTypeRef === RequestContentTypesMap['form-data']">
          <SimpleSchemaEditor v-model="contentValues[RequestContentTypesMap['form-data']].schema" :definitions="definitions" has-file />
        </div>

        <div v-if="currentContentTypeRef === RequestContentTypesMap['x-www-form-urlencoded']">
          <SimpleSchemaEditor v-model="contentValues[RequestContentTypesMap['x-www-form-urlencoded']].schema" :definitions="definitions" />
        </div>

        <div v-if="currentContentTypeRef === RequestContentTypesMap.json">
          <JSONSchemaEditor v-model="contentValues[RequestContentTypesMap.json].schema" :definitions="definitions" />
        </div>

        <div v-if="currentContentTypeRef === RequestContentTypesMap.xml">
          <JSONSchemaEditor v-model="contentValues[RequestContentTypesMap.xml].schema" :definitions="definitions" />
        </div>

        <div v-if="currentContentTypeRef === RequestContentTypesMap.raw">
          <CodeEditor v-model="contentValues[RequestContentTypesMap.raw].schema.example" />
        </div>

        <div v-if="currentContentTypeRef === RequestContentTypesMap.binary">
          <FileUploaderWrapper class="flex items-center border border-gray-200 border-solid rounded cursor-pointer h-30px" v-slot="{ fileName }" @change="handleChooseFile">
            <label class="flex items-center h-full bg-gray-200 border-r border-gray-200 border-solid px-6px font-500">{{ RequestContentTypesMap.binary }}</label>
            <span class="flex-1 truncate px-6px"> {{ fileName ?? $t('app.request.tips.selectFile') }}</span>
          </FileUploaderWrapper>
        </div>
      </el-tab-pane>

      <el-tab-pane label="Path">
        <template #label>
          Path
          <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px" v-if="pathsCount">{{ pathsCount }}</span>
        </template>
        <SimpleParameterEditor v-model="paths">
          <template #before>
            <GlobalParameter :data="globalPaths" :onSwitch="switchGlobalPath" />
          </template>
        </SimpleParameterEditor>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
<script setup lang="ts">
import SimpleParameterEditor from '@/components/APIEditor/SimpleEditor.vue'
import SimpleSchemaEditor from '@/components/APIEditor/SimpleSchemaEditor.vue'
import JSONSchemaEditor from '@/components/APIEditor/Editor.vue'
import CodeEditor from '@/components/APIEditor/CodeEditor.vue'
import { HttpDocument } from '@/typings'
import { useParameter } from './useParameter'
import { useContentType } from './useContentType'
import { Definition } from '../APIEditor/types'
import GlobalParameter from './GlobalParameter.vue'

const props = defineProps<{ modelValue: HttpDocument; definitions?: Definition[] }>()

const {
  headers,
  cookies,
  queries,
  paths,
  headersCount,
  cookiesCount,
  queriesCount,
  pathsCount,
  globalHeaders,
  globalCookies,
  globalPaths,
  globalQueries,
  switchGlobalHeader,
  switchGlobalCookie,
  switchGlobalPath,
  switchGlobalQuery,
} = useParameter(props)
const { RequestContentTypesMap, currentContentTypeRef, contentValues, bodyCount, handleChooseFile } = useContentType(props)
</script>
