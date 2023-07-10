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
import useProjectStore from '@/store/project'
import useDefinitionResponseStore from '@/store/definitionResponse'
import useShareStore from '@/store/share'

const projectStore = useProjectStore()
const globalParametersStore = useDefinitionParametersStore()
const definitionResponseStore = useDefinitionResponseStore()
const definitionSchemaStore = useDefinitionSchemaStore()
const shareStore = useShareStore()

const { project_id } = shareStore.sharedDocumentInfo!

onBeforeMount(async () => {
  try {
    await projectStore.getUrlServers(project_id)
    await globalParametersStore.getGlobalParameters(project_id)
    await definitionResponseStore.getDefinitions(project_id)
    await definitionSchemaStore.getDefinitions(project_id)
  } catch (error) {
    //
  }
})
</script>
