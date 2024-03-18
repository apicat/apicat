<script setup lang="ts">
import { useNamespace } from '@apicat/hooks'
import CollectionHistoryOperateHeader from './components/CollectionHistoryOperateHeader.vue'
import { useCollectionContext } from '@/hooks/useCollectionContext'
import useApi from '@/hooks/useApi'
import { getCollectionHistoryRecordDetail } from '@/api/project/collectionHistoryRecord'

const props = defineProps<{
  projectID: string
  collectionID: string
  historyID?: string
}>()

const AcEditor = defineAsyncComponent(() => import('@apicat/editor'))

const ns = useNamespace('document')
const collectionIDRef = computed(() => Number(props.collectionID))
const historyIDRef = computed(() => Number(props.historyID))
const { activeUrl, urls, parameters, responses, schemas, acEditorOptions } = useCollectionContext()
const [isLoading, getCollectionHistoryRecordDetailApi] = useApi(getCollectionHistoryRecordDetail)
const detail = ref<HistoryRecord.CollectionDetail>()

watch(historyIDRef, async (historyID) => {
  if (!historyID || Number.isNaN(historyID)) {
    detail.value = undefined
    return
  }

  detail.value = await getCollectionHistoryRecordDetailApi(props.projectID, collectionIDRef.value, historyID)
}, { immediate: true })
</script>

<template>
  <el-empty v-if="!detail" />
  <div v-else>
    <CollectionHistoryOperateHeader :title="detail?.title" :history-id="historyIDRef" :collection-id="collectionIDRef" :project-id="projectID" />
    <div v-loading="isLoading" :class="[ns.b()]">
      <AcEditor
        v-if="!isLoading"
        v-model:active-url="activeUrl"
        readonly
        :content="detail.content"
        :urls="urls"
        :schemas="schemas"
        :responses="responses"
        :parameters="parameters"
        :options="acEditorOptions"
      />
    </div>
  </div>
</template>
