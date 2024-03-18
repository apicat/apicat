<script setup lang="ts">
import { storeToRefs } from 'pinia'
import HistoryLayout from './HistoryLayout.vue'
import HistoryTree from './HistoryTree.vue'
import { useCollectionsStore } from '@/store/collections'
import { useCollectionProvider } from '@/hooks/useCollectionContext'
import { getCollectionHistoryPath } from '@/router/history'
import type Node from '@/components/AcTree/model/node'
import { ITERATION_COLLECTION_PATH_NAME, PROJECT_COLLECTION_PATH_NAME } from '@/router'

const props = defineProps<{
  projectID: string
  collectionID: string
}>()

useCollectionProvider({ projectID: props.projectID })

const router = useRouter()
const collectionsStore = useCollectionsStore()
const { historyRecord } = storeToRefs(collectionsStore)
const loading = ref(true)
const historyTreeRef = ref<InstanceType<typeof HistoryTree>>()

const activeKey = computed(() => {
  const historyID = Number.parseInt(router.currentRoute.value.params.historyID as string)
  return Number.isNaN(historyID) ? undefined : historyID
})

async function onNodeClick(node: Node) {
  await router.push({
    path: getCollectionHistoryPath(props.projectID, Number.parseInt(props.collectionID), node.key),
    query: router.currentRoute.value.query,
  })
}

// 返回
function goBack() {
  const iterationID = router.currentRoute.value.query.iterationID as string
  const { projectID: project_id, collectionID } = props
  // from itreation
  if (iterationID) {
    router.push({ name: ITERATION_COLLECTION_PATH_NAME, params: { iterationID, collectionID } })
    return
  }

  // 返回项目列表中的schema信息
  router.push({ name: PROJECT_COLLECTION_PATH_NAME, params: { project_id, collectionID } })
}
watch(activeKey, () => {
  activeKey.value && historyTreeRef.value?.setCurrentKey(activeKey.value)
})

onMounted(async () => {
  try {
    const histories = await collectionsStore.getHistories(
      props.projectID,
      Number.parseInt(props.collectionID as string),
    )
    loading.value = false
    await nextTick()
    // set default active
    if (histories.length && activeKey.value === undefined) await onNodeClick({ key: histories[0].id } as Node)
    // set current for expand
    activeKey.value && historyTreeRef.value?.setCurrentKey(activeKey.value)
  } catch (error) {
    loading.value = false
  }
})
</script>

<template>
  <HistoryLayout :go-back="goBack">
    <template #left>
      <HistoryTree v-if="!loading" ref="historyTreeRef" :history-record="historyRecord" @on-node-click="onNodeClick" />
      <div v-else class="h-full" v-loading="loading"></div>
    </template>
  </HistoryLayout>
</template>

<style lang="scss" scoped>
:deep(.el-loading-mask) {
  background: transparent;
}
</style>
