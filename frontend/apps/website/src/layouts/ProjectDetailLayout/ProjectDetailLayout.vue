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
      <template v-if="createMode">
        <DocumentCreatePage v-if="createMode === CreateModeEnum.document" />
        <SchemaCreatePage v-if="createMode === CreateModeEnum.schema" />
      </template>
      <router-view v-else />
    </div>
  </main>
</template>

<script setup lang="ts">
import { useNamespace } from '@/hooks/useNamespace'
import ProjectInfoHeader from './components/ProjectInfoHeader.vue'
import DirectoryTree from './components/DirectoryTree'
import SchemaTree from './components/SchemaTree'
import SchemaCreatePage from '@/views/document/SchemaCreatePage.vue'
import DocumentCreatePage from '@/views/document/DocumentCreatePage.vue'
import { uesAppStore, CreateModeEnum } from '@/store/app'
import { storeToRefs } from 'pinia'
import uesProjectStore from '@/store/project'
import { useParams } from '@/hooks/useParams'

const ns = useNamespace('doc-layout')
const appsStore = uesAppStore()
const projectStore = uesProjectStore()
const { project_id } = useParams()
const { createMode } = storeToRefs(appsStore)

const directoryTree = ref<InstanceType<typeof DirectoryTree>>()
const schemaTree = ref<InstanceType<typeof SchemaTree>>()

provide('directoryTree', {
  updateTitle: (id: any, title: string) => directoryTree.value?.updateTitle(id, title),
  createNodeByData: (data: any) => directoryTree.value?.createNodeByData(data),
})

provide('schemaTree', {
  updateTitle: (id: any, title: string) => schemaTree.value?.updateTitle(id, title),
  activeNode: (id: any) => schemaTree.value?.activeNode(id),
})

onMounted(async () => {
  await projectStore.getUrlServers(project_id as string)
})
</script>
