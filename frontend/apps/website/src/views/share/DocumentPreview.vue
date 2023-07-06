<template>
  <div class="p-30px" v-loading="isLoading">
    <h1 class="ac-document__title">{{ httpDoc?.title }}</h1>
    <div class="mt-10px" v-if="httpDoc">
      <RequestMethodRaw class="mb-10px" :doc="httpDoc" :urls="urlServers" />
      <RequestParamRaw class="mb-10px" :doc="httpDoc" :definitions="definitions" />
      <ResponseParamTabsRaw :doc="httpDoc" :definitions="definitions" :project-id="project_id" />
    </div>
  </div>
</template>
<script setup lang="ts">
import { HttpDocument } from '@/typings'
import ResponseParamTabsRaw from '@/components/ResponseParamTabsRaw.vue'
import { storeToRefs } from 'pinia'
import { getCollectionDetail } from '@/api/collection'
import useDefinitionStore from '@/store/definition'
import useProjectStore from '@/store/project'
import useShareStore from '@/store/share'

const projectStore = useProjectStore()
const definitionStore = useDefinitionStore()
const shareStore = useShareStore()

const [isLoading, getCollectionDetailApi] = getCollectionDetail()
const { urlServers } = storeToRefs(projectStore)
const { definitions } = storeToRefs(definitionStore)

const httpDoc: Ref<HttpDocument | null> = ref(null)
const { project_id, collection_id } = shareStore.sharedDocumentInfo!

onMounted(async () => {
  httpDoc.value = await getCollectionDetailApi({ project_id, collection_id })
})
</script>
