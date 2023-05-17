<template>
  <div v-if="responseRef.header && responseRef.header.length">
    <p class="mb-5px text-16px mt-10px">Header</p>
    <SimpleParameterEditor :readonly="true" v-model="responseRef.header" />
  </div>

  <div v-for="(_, contentTypeKey) in responseRef.content" :key="contentDefaultType">
    <p class="mb-10px text-16px mt-14px">Content({{ contentTypeKey }})</p>
    <JSONSchemaEditor :readonly="true" :definitions="definitionSchemas" v-model="responseRef.content[contentTypeKey].schema" v-if="isJsonSchema" />
    <p class="my-10px" v-if="responseRef.content[contentTypeKey].schema.example">
      {{ $t('app.response.tips.responseExample') }}
    </p>
    <CodeEditor
      v-if="responseRef.content[contentTypeKey].schema.example"
      class="mt-14px"
      :readonly="true"
      :lang="contentTypes[contentTypeKey]"
      v-model="responseRef.content[contentTypeKey].schema.example"
    />
  </div>
</template>

<script setup lang="ts">
import type { DefinitionSchema } from '../APIEditor/types'
import SimpleParameterEditor from '../APIEditor/SimpleEditor.vue'
import JSONSchemaEditor from '../APIEditor/Editor.vue'
import CodeEditor from '../APIEditor/CodeEditor.vue'
import { DefinitionResponse } from '@/typings'
import { useDefinitionResponse, contentTypes } from './useDefinitionResponse'

const props = defineProps<{
  response: DefinitionResponse
  definitionSchemas?: DefinitionSchema[]
}>()

const { responseRef, contentDefaultType, isJsonSchema } = useDefinitionResponse(props)
</script>
