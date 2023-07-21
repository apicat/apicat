<template>
  <div :class="nsRow.b()">
    <div :class="nsRow.e('content')">
      <div :class="[nsRow.e('item'), nsRow.e('name')]">
        <div :style="intentRowStyle"></div>
        <input v-model="data.name" :disabled="data.isInRefSchemaNode" />
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('type')]">
        {{ data.isRefSchemaNode ? (data as RefSchemaNode).definitionSchema?.name : data.schema.type }}
        <select v-model="type" class="hidden">
          <option value="string">string</option>
          <option value="number">number</option>
          <option value="boolean">boolean</option>
          <option value="object">object</option>
          <option value="array">array</option>
          <option value="null">null</option>
          <option value="any">any</option>
        </select>
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('required')]">
        <input type="checkbox" :disabled="data.isInRefSchemaNode || data.isConstantSchemaNode" v-model="data.isRequired" />
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('example')]">
        <input v-model="data.example" :disabled="data.isInRefSchemaNode" />
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('description')]">
        <input v-model="data.description" :disabled="data.isInRefSchemaNode" />
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('operation')]"></div>
    </div>

    <div :class="nsRow.e('children')">
      <div :style="intentLineStyle" :class="nsRow.e('line')"></div>
      <SchemaEditorRaw v-for="item in data.childNodes" :key="item.id" :data="item" />
    </div>
  </div>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks/useNamespace'
import SchemaNode from './schema/SchemaNode'
import RefSchemaNode from './schema/compose/RefSchemaNode'

const props = defineProps<{ data: SchemaNode | RefSchemaNode }>()

const nsRow = useNamespace('schema-row')
const intent = props.data.level * 16 + 12
const intentLineStyle = { left: intent + 'px' }
const intentRowStyle = { width: intent + 'px' }

const type = ref('string')
</script>
