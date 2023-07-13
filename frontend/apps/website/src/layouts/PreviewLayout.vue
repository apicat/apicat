<template>
  <template v-if="isExistSecretKey">
    <PreviewHeader />
    <main class="flex flex-col min-h-screen bg-gray-100 p-30px pt78px">
      <div class="flex-1 bg-white">
        <RouterView />
      </div>
    </main>
    <ac-backtop :bottom="100" :right="100" />
  </template>
  <DocumentVerification v-else />
</template>
<script setup lang="ts">
import PreviewHeader from '@/views/share/components/PreviewHeader.vue'
import DocumentVerification from '@/views/share/DocumentVerification.vue'

import { useDefinitionSchemaStore } from '@/store/definition'
import { useDefinitionParametersStore } from '@/store/globalParameters'
import useProjectStore from '@/store/project'
import useDefinitionResponseStore from '@/store/definitionResponse'
import useShareStore from '@/store/share'
import { storeToRefs } from 'pinia'

const projectStore = useProjectStore()
const globalParametersStore = useDefinitionParametersStore()
const definitionResponseStore = useDefinitionResponseStore()
const definitionSchemaStore = useDefinitionSchemaStore()
const shareStore = useShareStore()

const { sharedDocumentInfo, isExistSecretKey } = storeToRefs(shareStore)

onBeforeMount(async () => {
  if (!isExistSecretKey.value) {
    return
  }

  try {
    const { project_id } = sharedDocumentInfo.value!
    await projectStore.getUrlServers(project_id)
    await globalParametersStore.getGlobalParameters(project_id)
    await definitionResponseStore.getDefinitions(project_id)
    await definitionSchemaStore.getDefinitions(project_id)
  } catch (error) {
    //
  }
})
</script>
