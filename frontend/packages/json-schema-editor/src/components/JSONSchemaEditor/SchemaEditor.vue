<script setup lang="ts">
import SchemaEditorRaw from './SchemaEditorRaw.vue'
import { useNamespace } from '@/hooks/useNamespace'
import { JSONSchema } from './types'
import SchemaStore from './schema/SchemaStore'

const props = defineProps<{
  schema: JSONSchema
  definitionSchemas: Array<any>
}>()

const emits = defineEmits(['update:schema'])

const nsEditor = useNamespace('schema-editor')
const nsRow = useNamespace('schema-row')
const { schema, definitionSchemas } = toRefs(props)

const store = ref(new SchemaStore(props.schema, props.definitionSchemas, (schema) => schema && emits('update:schema', schema)))

watch(schema, () => {
  store.value.setSchema(schema.value)
})

watch(definitionSchemas, () => {
  store.value.setDefinitionSchemas(definitionSchemas.value)
})
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

    <SchemaEditorRaw v-if="store.root" :data="store.root" />
  </div>
</template>

<style lang="scss">
@use './style.scss';
</style>
