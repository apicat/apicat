<script setup lang="ts">
import JSONSchemaEditor from './components/JSONSchemaEditor'
import { exampleSchema, definitionSchemas } from '@/mock/schemas'

const schema = ref(exampleSchema)
</script>

<template>
  <div class="m-auto w-80vw mt-10vh">
    <JSONSchemaEditor :schema="schema" :definition-schemas="definitionSchemas" />
  </div>
</template>
