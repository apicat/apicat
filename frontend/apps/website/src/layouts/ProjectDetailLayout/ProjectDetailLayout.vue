<template>
  <main :class="ns.b()">
    <ProjectInfoHeader />

    <div :class="ns.e('left')">
      <div class="flex flex-col h-full overflow-y-scroll scroll-content">
        <DirectoryTree ref="directoryTree" />
        <div class="my-10px"></div>
        <SchemaTree ref="schemaTree" />
      </div>
    </div>
    <div class="scroll-content" :class="ns.e('right')">
      <router-view />
    </div>
  </main>
</template>

<script setup lang="ts">
import { useNamespace } from '@/hooks/useNamespace'
import ProjectInfoHeader from './components/ProjectInfoHeader.vue'
import DirectoryTree from './components/DirectoryTree'
import SchemaTree from './components/SchemaTree'
import uesProjectStore from '@/store/project'
import uesGlobalParametersStore from '@/store/globalParameters'
import { useParams } from '@/hooks/useParams'
import useCommonResponseStore from '@/store/commonResponse'

const ns = useNamespace('doc-layout')
const projectStore = uesProjectStore()
const globalParametersStore = uesGlobalParametersStore()
const commonResponseStore = useCommonResponseStore()
const { project_id } = useParams()

const directoryTree = ref<InstanceType<typeof DirectoryTree>>()
const schemaTree = ref<InstanceType<typeof SchemaTree>>()

provide('directoryTree', {
  updateTitle: (id: any, title: string) => directoryTree.value?.updateTitle(id, title),
  createNodeByData: (data: any) => directoryTree.value?.createNodeByData(data),
  reload: () => directoryTree.value?.reload(),
  reactiveNode: () => directoryTree.value?.reactiveNode(),
  redirecToDocumentDetail: (activeId?: any) => directoryTree.value?.redirecToDocumentDetailPage(activeId),
})

provide('schemaTree', {
  updateTitle: (id: any, title: string) => schemaTree.value?.updateTitle(id, title),
  activeNode: (id: any) => schemaTree.value?.activeNode(id),
  reactiveNode: () => schemaTree.value?.reactiveNode(),
  redirecToSchemaDetail: (activeId?: any) => schemaTree.value?.redirecToSchemaEdit(activeId),
})

onMounted(async () => {
  await projectStore.getUrlServers(project_id as string)
  await globalParametersStore.getGlobalParameters(project_id as string)
  await commonResponseStore.getCommonResponseList(project_id as string)
})
</script>
