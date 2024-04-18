<script setup lang="ts">
import { storeToRefs } from 'pinia'
import ProjectInfoHeader from './components/ProjectInfoHeader.vue'
import CollectionTree from './components/CollectionTree'
import DefinitionResponseTree from './components/DefinitionResponseTree'
import DefinitionSchemaTree from './components/DefinitionSchemaTree'
import { useActiveTree } from './composables/useActiveTree'
import { useProjectLayoutProvider } from './composables/useProjectLayoutContext'
import { useRefresh } from './composables/useRefresh'
import ProjectVerification from '@/views/share/ProjectVerification.vue'
import ProjectShareModal from '@/views/project/components/ProjectShareModal.vue'
import DocumentShareModal from '@/views/collection/components/DocumentShareModal.vue'
import ExportDocumentModal from '@/views/collection/components/ExportDocumentModal.vue'
import { useNamespace } from '@/hooks/useNamespace'
import useProjectStore from '@/store/project'
import { useCollectionsStore } from '@/store/collections'
import { useCollectionContextWithoutMounted } from '@/hooks/useCollectionContext'
import useApi from '@/hooks/useApi'
import AICreateDialog from '@/layouts/ProjectDetailLayout/components/AICreateDialog.vue'
import { providePagesMode } from '@/layouts/ProjectDetailLayout/composables/usePagesMode'
import { useGlobalLoading } from '@/hooks/useGlobalLoading'
import { provideAsyncInitTask } from '@/hooks/useWaitAsyncTask'

const ns = useNamespace('doc-layout')
const { exportDocModalRef, shareDocModalRef, shareProjectModalRef, AIDialogRef } = useProjectLayoutProvider()
const { currentNode, activeCollectionKey, activeResponseKey, activeSchemaKey } = useActiveTree()

const projectStore = useProjectStore()
const collectionStore = useCollectionsStore()
const { isShowProjectSecretLayer, projectID } = storeToRefs(projectStore)
const { hideGlobalLoading } = useGlobalLoading()
const collectionTreeRef = ref<InstanceType<typeof CollectionTree>>()
const schemaTreeRef = ref<InstanceType<typeof DefinitionSchemaTree>>()
const responseTreeRef = ref<InstanceType<typeof DefinitionResponseTree>>()
const { initContextData, context } = useCollectionContextWithoutMounted()

async function initialize() {
  await initContextData(projectID.value!)
  await collectionStore.getCollections(projectID.value!)

  if (currentNode.value.id === undefined) {
    const selectTask = collectionTreeRef.value?.selectFirstNode()
    selectTask && ctx.addTask(selectTask)
    return
  }

  collectionTreeRef.value?.expandOnStartup()
  schemaTreeRef.value?.expandOnStartup()
  responseTreeRef.value?.expandOnStartup()
}

const [isLoading, init] = useApi(initialize, { defaultLoadingStatus: true })

watch(
  isShowProjectSecretLayer,
  async () => {
    if (!isShowProjectSecretLayer.value)
      await init()
  },
)

const ctx = provideAsyncInitTask()

onMounted(async () => {
  !isShowProjectSecretLayer.value && ctx.addTask(init)

  try {
    await ctx.done()
    console.log('layout mounted task done')
  }
  finally {
    hideGlobalLoading()
  }
})

onUnmounted(() => collectionStore.$reset())

// 解引用后同步相关数据
useRefresh(activeCollectionKey, activeSchemaKey, activeResponseKey, context)

// page mode
providePagesMode()
</script>

<template>
  <ProjectVerification v-if="projectStore.isShowProjectSecretLayer" />
  <div v-else>
    <main :class="ns.b()">
      <ProjectInfoHeader />
      <div v-loading="isLoading" :class="ns.e('left')">
        <div class="flex flex-col h-full overflow-y-auto scroll-content">
          <CollectionTree ref="collectionTreeRef" />
          <DefinitionSchemaTree ref="schemaTreeRef" />
          <DefinitionResponseTree ref="responseTreeRef" />
        </div>
      </div>

      <div class="scroll-content" :class="ns.e('right')">
        <RouterView :project_id="projectID" />
      </div>
    </main>
    <ExportDocumentModal ref="exportDocModalRef" />
    <DocumentShareModal ref="shareDocModalRef" />
    <ProjectShareModal ref="shareProjectModalRef" />
    <AICreateDialog ref="AIDialogRef" />
  </div>
</template>
