<script setup lang="ts">
import { storeToRefs } from 'pinia'
import HistoryLayout from './HistoryLayout.vue'
import HistoryTree from './HistoryTree.vue'
import useDefinitionSchemaStore from '@/store/definitionSchema'
import type Node from '@/components/AcTree/model/node'
import { getSchemaHistoryPath } from '@/router/history'
import { ITERATION_SCHEMA_PATH_NAME, PROJECT_SCHEMA_PATH_NAME } from '@/router'

const props = defineProps<{
  projectID: string
  schemaID: string
}>()

const router = useRouter()
const schemaStore = useDefinitionSchemaStore()
const { historyRecord } = storeToRefs(schemaStore)
const loading = ref(true)
const historyTreeRef = ref<InstanceType<typeof HistoryTree>>()

const activeKey = computed(() => {
  const historyID = Number.parseInt(router.currentRoute.value.params.historyID as string)
  return Number.isNaN(historyID) ? undefined : historyID
})

watch(activeKey, () => {
  activeKey.value && historyTreeRef.value?.setCurrentKey(activeKey.value)
})

async function onNodeClick(node: Node) {
  await router.push({
    path: getSchemaHistoryPath(props.projectID, Number.parseInt(props.schemaID), node.key),
    query: router.currentRoute.value.query,
  })
}

// 返回
function goBack() {
  const iterationID = router.currentRoute.value.query.iterationID as string
  const { projectID: project_id, schemaID } = props
  // from itreation
  if (iterationID) {
    router.push({ name: ITERATION_SCHEMA_PATH_NAME, params: { iterationID, schemaID } })
    return
  }

  // 返回项目列表中的schema信息
  router.push({ name: PROJECT_SCHEMA_PATH_NAME, params: { project_id, schemaID } })
}

onMounted(async () => {
  try {
    const [histories] = await Promise.all([
      schemaStore.getHistories(props.projectID as string, Number.parseInt(props.schemaID as string)),
      schemaStore.getSchemas(props.projectID),
    ])

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
