<template>
  <DocumentHistoryOperateHeader :title="httpDoc?.title" />
  <div :class="[ns.b(), { 'h-20vh': !httpDoc && hasDocument }]" v-loading="isLoading">
    <div class="ac-editor mt-10px" v-if="httpDoc">
      <RequestMethodRaw class="mb-10px" :doc="httpDoc" :urls="urlServers" />
      <RequestParamRaw class="mb-10px" :doc="httpDoc" :definitions="definitions" />
      <ResponseParamTabsRaw :doc="httpDoc" :definitions="definitions" :project-id="project_id" />
    </div>
  </div>
</template>
<script setup lang="ts">
import DocumentHistoryOperateHeader from '@/views/document/components/DocumentHistoryOperateHeader.vue'
import { getDocumentHistoryRecordDetail } from '@/api/collection'
import { useNamespace } from '@/hooks/useNamespace'
import { useParams } from '@/hooks/useParams'
import { useDefinitionSchemaStore } from '@/store/definitionSchema'
import useProjectStore from '@/store/project'
import { HttpDocument } from '@/typings'
import { storeToRefs } from 'pinia'
import useApi from '@/hooks/useApi'

const route = useRoute()

const projectStore = useProjectStore()
const definitionSchemaStore = useDefinitionSchemaStore()

const ns = useNamespace('document')

const { project_id, computedRouteParams } = useParams()

const hasDocument = ref(false)

const httpDoc: Ref<HttpDocument | null> = ref(null)

const [isLoading, getDocumentHistoryRecordDetailApi] = useApi(getDocumentHistoryRecordDetail)

const { urlServers } = storeToRefs(projectStore)
const { definitions } = storeToRefs(definitionSchemaStore)

const getDetail = async (hid: string) => {
  const history_id = parseInt(hid, 10)

  isLoading.value = true

  if (isNaN(history_id)) {
    hasDocument.value = false
    httpDoc.value = null
    isLoading.value = false
    return
  }

  try {
    hasDocument.value = true
    const { doc_id: collection_id } = unref(computedRouteParams)
    httpDoc.value = await getDocumentHistoryRecordDetailApi({ project_id, collection_id, history_id })
  } catch (error) {
    //
  }
}

watch(
  () => route.params.history_id,
  async () => await getDetail(route.params.history_id as string),
  {
    immediate: true,
  }
)
</script>
