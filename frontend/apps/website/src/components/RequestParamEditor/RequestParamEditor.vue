<template>
  <div>
    <h2 class="text-16px font-500">请求参数</h2>
    <el-tabs>
      <el-tab-pane label="Header">
        <template #label>
          Header
          <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px" v-if="headersCount">{{ headersCount }}</span>
        </template>
        <SimpleParameterEditor v-model="headers" />
      </el-tab-pane>
      <el-tab-pane label="Cookie">
        <template #label>
          Cookie
          <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px" v-if="cookiesCount">{{ cookiesCount }}</span>
        </template>
        <SimpleParameterEditor v-model="cookies" />
      </el-tab-pane>
      <el-tab-pane label="Query">
        <template #label>
          Query
          <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px" v-if="queriesCount">{{ queriesCount }}</span>
        </template>
        <SimpleParameterEditor v-model="queries" />
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

        <div v-show="currentContentTypeRef === RequestContentTypesMap.none" class="text-center">该请求没有Body体</div>

        <div v-show="currentContentTypeRef === RequestContentTypesMap['form-data']">
          <SimpleSchemaEditor v-model="contentValues[RequestContentTypesMap['form-data']].schema" :definitions="definitions" has-file />
        </div>

        <div v-show="currentContentTypeRef === RequestContentTypesMap['x-www-form-urlencoded']">
          <SimpleSchemaEditor v-model="contentValues[RequestContentTypesMap['x-www-form-urlencoded']].schema" :definitions="definitions" />
        </div>

        <div v-if="currentContentTypeRef === RequestContentTypesMap.json">
          <JSONSchemaEditor v-model="contentValues[RequestContentTypesMap.json].schema" :definitions="definitions" />
        </div>

        <div v-show="currentContentTypeRef === RequestContentTypesMap.xml">
          <JSONSchemaEditor v-model="contentValues[RequestContentTypesMap.xml].schema" :definitions="definitions" />
        </div>

        <div v-show="currentContentTypeRef === RequestContentTypesMap.raw">
          <CodeEditor v-model="contentValues[RequestContentTypesMap.raw].schema.example" />
        </div>

        <div v-show="currentContentTypeRef === RequestContentTypesMap.binary">
          <FileUploaderWrapper class="flex items-center border border-gray-200 border-solid rounded cursor-pointer h-30px" v-slot="{ fileName }" @change="handleChooseFile">
            <label class="flex items-center h-full bg-gray-200 border-r border-gray-200 border-solid px-6px font-500">{{ RequestContentTypesMap.binary }}</label>
            <span class="flex-1 truncate px-6px"> {{ fileName ?? '请选择文件' }}</span>
          </FileUploaderWrapper>
        </div>
      </el-tab-pane>
      <el-tab-pane label="Path">
        <template #label>
          Path
          <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px" v-if="pathsCount">{{ pathsCount }}</span>
        </template>
        <SimpleParameterEditor v-model="paths" />
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

const props = defineProps<{ modelValue: HttpDocument; definitions?: Definition[] }>()

const { headers, cookies, queries, paths, headersCount, cookiesCount, queriesCount, pathsCount } = useParameter(props)
const { RequestContentTypesMap, currentContentTypeRef, contentValues, bodyCount, handleChooseFile } = useContentType(props)
</script>
