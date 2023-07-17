<template>
  <div :class="nsRow.b()">
    <div :class="nsRow.e('content')">
      <div :class="[nsRow.e('item'), nsRow.e('name')]">
        <div :style="intentRowStyle"></div>
        <input v-model="data.schemaName" />
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('type')]">
        {{ data.isRefSchema ? data.name : data.schema.type }}
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('required')]">
        <input type="checkbox" v-model="data.isRequired" />
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('example')]">Example</div>
      <div :class="[nsRow.e('item'), nsRow.e('description')]">Description</div>
      <div :class="[nsRow.e('item'), nsRow.e('operation')]"></div>
    </div>

    <div :class="nsRow.e('children')">
      <div :style="intentLineStyle" :class="nsRow.e('line')"></div>
      <template v-for="item in data.childNodes" :key="item.id">
        <SchemaEditorRaw :data="item" />
        <!-- <SchemaEditorRaw v-if="item.isConstantNode && item.childNodes.length" :data="item.childNodes[0]" /> -->
        <!-- <SchemaEditorRaw v-if="item.isRef && item.childNodes.length" :data="item.childNodes[0]" /> -->
      </template>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks/useNamespace'

const props = defineProps<{ data: any }>()

const nsRow = useNamespace('schema-row')
const intent = props.data.level * 16 + 12
const intentLineStyle = { left: intent + 'px' }
const intentRowStyle = { width: intent + 'px' }

onMounted(() => {
  //   console.log(props.data)
})
</script>
