<template>
  <template v-if="!isShowProjectSecretLayer">
    <main :class="ns.b()">
      <ProjectInfoHeader />
      <div :class="ns.e('left')">
        <div class="flex flex-col h-full overflow-y-scroll scroll-content">
          <DirectoryTree ref="directoryTree" />
          <div class="my-10px"></div>
          <SchemaTree ref="schemaTree" />
          <div class="my-10px"></div>
          <DefinitionResponseTree ref="definitionResponseTree" />
        </div>
      </div>
      <div class="scroll-content" :class="ns.e('right')">
        <router-view />
      </div>
    </main>
    <ExportDocumentModal ref="exportDocumentModalRef" />
    <DocumentShareModal ref="documentShareModalRef" />
    <ProjectShareModal ref="projectShareModalRef" />
  </template>
  <ProjectVerification v-else />
</template>

<script setup lang="ts">
import ProjectVerification from '@/views/share/ProjectVerification.vue'
import ExportDocumentModal from '@/views/component/ExportDocumentModal.vue'
import DocumentShareModal from '@/views/document/components/DocumentShareModal.vue'
import ProjectShareModal from '@/views/project/components/ProjectShareModal.vue'
import { useNamespace } from '@/hooks/useNamespace'
import ProjectInfoHeader from './components/ProjectInfoHeader.vue'
import DirectoryTree from './components/DirectoryTree'
import SchemaTree from './components/SchemaTree'
import DefinitionResponseTree from './components/DefinitionResponseTree'
import useProjectStore from '@/store/project'
import uesGlobalParametersStore from '@/store/globalParameters'
import { useParams } from '@/hooks/useParams'
import { ProjectDetailModalsContextKey } from './constants'
import { storeToRefs } from 'pinia'
import useDefinitionStore from '@/store/definitionSchema'
import useDefinitionResponseStore from '@/store/definitionResponse'

const ns = useNamespace('doc-layout')
const projectStore = useProjectStore()
const globalParametersStore = uesGlobalParametersStore()
const definitionStore = useDefinitionStore()
const definitionResponseStore = useDefinitionResponseStore()

const { project_id } = useParams()

const { isShowProjectSecretLayer } = storeToRefs(projectStore)
const directoryTree = ref<InstanceType<typeof DirectoryTree>>()
const schemaTree = ref<InstanceType<typeof SchemaTree>>()
const definitionResponseTree = ref<InstanceType<typeof DefinitionResponseTree>>()
const exportDocumentModalRef = ref<InstanceType<typeof ExportDocumentModal>>()
const documentShareModalRef = ref<InstanceType<typeof DocumentShareModal>>()
const projectShareModalRef = ref<InstanceType<typeof ProjectShareModal>>()

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
  reload: async () => await schemaTree.value?.reload(),
  reactiveNode: () => schemaTree.value?.reactiveNode(),
  redirecToSchemaDetail: (activeId?: any) => schemaTree.value?.redirecToSchemaEdit(activeId),
})

provide('definitionResponseTree', {
  updateTitle: (id: any, title: string) => definitionResponseTree.value?.updateTitle(id, title),
  activeNode: (id: any) => definitionResponseTree.value?.activeNode(id),
  reload: async () => await definitionResponseTree.value?.reload(),
  reactiveNode: () => definitionResponseTree.value?.reactiveNode(),
})

provide('exportModal', {
  exportDocument: (project_id?: string, doc_id?: string) => exportDocumentModalRef.value?.show(project_id, doc_id),
})

provide(ProjectDetailModalsContextKey, {
  exportDocument: (project_id?: string, doc_id?: string | number) => exportDocumentModalRef.value?.show(project_id, doc_id),
  shareDocument: (project_id: string, doc_id: string) => documentShareModalRef.value?.show({ project_id, collection_id: doc_id }),
  shareProject: (project_id: string) => projectShareModalRef.value?.show({ project_id }),
})

globalParametersStore.$onAction(({ name, after }) => {
  if (name === 'deleteGlobalParameter') {
    after(() => globalParametersStore.getGlobalParameters(project_id as string))
  }
})

definitionStore.$onAction(({ name, after }) => {
  if (name === 'deleteDefinition') {
    after(async () => {
     await definitionStore.getDefinitions(project_id as string)
     await definitionResponseStore.getDefinitions(project_id as string)
    })
  }
})

definitionResponseStore.$onAction(({ name, after }) => {
  if (name === 'deleteDefinition') {
    after(() =>  definitionResponseStore.getDefinitions(project_id as string))
  }
})

onMounted(async () => {
  if (isShowProjectSecretLayer.value) {
    return
  }

  await projectStore.getUrlServers(project_id as string)
  await globalParametersStore.getGlobalParameters(project_id as string)
})
</script>
