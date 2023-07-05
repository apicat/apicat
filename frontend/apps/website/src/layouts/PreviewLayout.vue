<template>
  <PreviewHeader />
  <main class="flex flex-col min-h-screen bg-gray-100 p-30px pt78px">
    <div class="flex-1 bg-white">
      <RouterView />
    </div>
  </main>
  <ac-backtop :bottom="100" :right="100" />
</template>
<script setup lang="ts">
import PreviewHeader from '@/views/share/components/PreviewHeader.vue'
import { useDefinitionSchemaStore } from '@/store/definition'
import { useDefinitionParametersStore } from '@/store/globalParameters'
import uesProjectStore from '@/store/project'
import useDefinitionResponseStore from '@/store/definitionResponse'

const projectStore = uesProjectStore()
const globalParametersStore = useDefinitionParametersStore()
const definitionResponseStore = useDefinitionResponseStore()
const definitionSchemaStore = useDefinitionSchemaStore()
</script>
