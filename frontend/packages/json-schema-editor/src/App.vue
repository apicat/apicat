<script setup lang="ts">
import JSONSchemaEditor from './components/JSONSchemaEditor'
import { exampleSchema, definitionSchemas } from '@/mock/schemas'
import SchemaStore from './components/JSONSchemaEditor/schema/SchemaStore'

new SchemaStore(exampleSchema, definitionSchemas as any)
</script>

<template>
  <div class="m-auto w-80vw mt-10vh">
    <JSONSchemaEditor :schema="exampleSchema" :definition-schemas="definitionSchemas" />
  </div>
</template>
