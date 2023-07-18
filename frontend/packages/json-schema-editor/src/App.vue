<script setup lang="ts">
import JSONSchemaEditor from './components/JSONSchemaEditor'
import { exampleSchema, definitionSchemas } from '@/mock/schemas'

// const schema = ref({})
const ss = ref({})
const ds: Ref<Array<any>> = ref([])

const schema = computed({
  get: () => ss.value,
  set: (newVal) => {
    console.log(toRaw(newVal))
  },
})

setTimeout(() => {
  // schema.value = exampleSchema
  ss.value = exampleSchema
}, 500)

setTimeout(() => {
  ds.value = definitionSchemas
}, 1000)
</script>

<template>
  <div class="m-auto w-80vw mt-10vh">
    <JSONSchemaEditor v-model:schema="schema" :definition-schemas="ds" />
  </div>
</template>
