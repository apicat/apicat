<template>
  <div class="mt-10px">
    <ToggleHeading title="Header" :expand="responseRef.header && !!responseRef.header.length">
      <SimpleEditor v-model="responseRef.header" />
    </ToggleHeading>

    <ToggleHeading title="Content">
      <template #extra>
        <el-select :modelValue="contentDefaultType" @change="changeContentType">
          <template #prefix>Content-Type</template>
          <el-option v-for="(_, cti) in contentTypes" :key="cti" :label="cti" :value="cti" />
        </el-select>
      </template>

      <div v-for="(_, ct) in responseRef.content" :key="ct" class="mt-10px">
        <Editor :definitions="definitionSchemas" v-model="responseRef.content[ct].schema" v-if="isJsonSchema" />
        <p class="my-10px">
          <span class="mr-4px">{{ $t('app.response.tips.responseExample') }}</span>
          <el-tag disable-transitions effect="plain">format:{{ contentTypes[ct] }}</el-tag>
        </p>
        <CodeEditor v-model="responseRef.content[ct].schema.example" :lang="contentTypes[ct]" />
      </div>
    </ToggleHeading>
  </div>
</template>

<script setup lang="ts">
import type { DefinitionSchema } from '../APIEditor/types'
import SimpleEditor from '../APIEditor/SimpleEditor.vue'
import Editor from '../APIEditor/Editor.vue'
import CodeEditor from '../APIEditor/CodeEditor.vue'
import { DefinitionResponse } from '@/typings'
import { useDefinitionResponse, contentTypes } from './useDefinitionResponse'

const props = defineProps<{
  response: DefinitionResponse
  definitionSchemas?: DefinitionSchema[]
}>()

const { responseRef, contentDefaultType, isJsonSchema, changeContentType } = useDefinitionResponse(props)
</script>
