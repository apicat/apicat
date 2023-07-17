<script setup lang="ts">
import SchemaEditorRaw from './SchemaEditorRaw.vue'
import { useNamespace } from '@/hooks/useNamespace'
import { JSONSchema } from './types'
// import SchemaTreeStore from './model/SchemaStore'
import SchemaStore from './schema/SchemaStore'

const props = defineProps<{
  schema: JSONSchema
  definitionSchemas: Array<any>
}>()

const nsEditor = useNamespace('schema-editor')
const nsRow = useNamespace('schema-row')

// const store = new SchemaTreeStore({ schema: props.schema, definitionSchemas: props.definitionSchemas })
const store = new SchemaStore(props.schema, props.definitionSchemas)

const root = ref(store.root.rootNode)

window['root'] = store.root
</script>

<template>
  <div :class="nsEditor.b()">
    <div :class="[nsRow.b(), nsRow.m('header')]">
      <div :class="nsRow.e('content')">
        <div :class="nsRow.e('name')">Name</div>
        <div :class="nsRow.e('type')">Type</div>
        <div :class="nsRow.e('required')">Required</div>
        <div :class="nsRow.e('example')">Example</div>
        <div :class="nsRow.e('description')">Description</div>
        <div :class="nsRow.e('operation')"></div>
      </div>
    </div>

    <SchemaEditorRaw :data="root" />
  </div>
</template>

<style lang="scss">
@use './style.scss';
</style>
