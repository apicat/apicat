<template>
  <div class="ac-header-operate" v-if="hasDocument">
    <div class="ac-header-operate__main">
      <p class="ac-header-operate__title">{{ httpDoc.title }}</p>
    </div>

    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="goDocumentEditPage()">编辑</el-button>
    </div>
  </div>

  <div :class="ns.b()" v-loading="isLoading" v-if="hasDocument">
    <div class="ac-editor mt-10px">
      <RequestMethodRaw class="mb-10px" :doc="httpDoc" :urls="urlServers" />

      <RequestParamRaw class="mb-10px" :doc="httpDoc" :definitions="definitions" />

      <ResponseParamTabsRaw :doc="httpDoc" :definitions="definitions" />
    </div>
  </div>

  <Result v-if="!hasDocument">
    <template #icon>
      <img class="h-auto w-260px mb-26px" src="@/assets/images/icon-empty.png" alt="" />
    </template>
    <template #title>
      <div class="m-auto">您当前尚未创建接口，请从左侧目录栏点击添加 API 接口。</div>
    </template>
  </Result>
</template>
<script setup lang="ts">
import { HttpDocument } from '@/typings'
import { useNamespace } from '@/hooks/useNamespace'
import { createHttpDocument } from '@/views/document/components/createHttpDocument'
import ResponseParamTabsRaw from '@/components/ResponseParamTabsRaw.vue'
import { useGoPage } from '@/hooks/useGoPage'
import uesProjectStore from '@/store/project'
import { storeToRefs } from 'pinia'
import { getCollectionDetail } from '@/api/collection'
import { useParams } from '@/hooks/useParams'
import useDefinitionStore from '@/store/definition'
import uesGlobalParametersStore from '@/store/globalParameters'
import useCommonResponseStore from '@/store/commonResponse'

const projectStore = uesProjectStore()
const definitionStore = useDefinitionStore()
const globalParametersStore = uesGlobalParametersStore()
const commonResponseStore = useCommonResponseStore()
const route = useRoute()
const { project_id } = useParams()
const { goDocumentEditPage } = useGoPage()

const [isLoading, getCollectionDetailApi] = getCollectionDetail()
const { urlServers } = storeToRefs(projectStore)
const { definitions } = storeToRefs(definitionStore)

const hasDocument = ref(true)
const ns = useNamespace('document')
const httpDoc: Ref<HttpDocument> = ref(createHttpDocument())

const getDetail = async (docId: string) => {
  const doc_id = parseInt(docId, 10)

  if (isNaN(doc_id)) {
    hasDocument.value = false
    return
  }

  try {
    hasDocument.value = true
    httpDoc.value = await getCollectionDetailApi({ project_id, collection_id: doc_id })
  } catch (error) {
    console.error(error)
  }
}

globalParametersStore.$onAction(({ name, after }) => {
  // 删除全局参数
  if (name === 'deleteGlobalParameter') {
    after(() => getDetail(route.params.doc_id as string))
  }
})

commonResponseStore.$onAction(({ name, after }) => {
  // 删除全局响应
  if (name === 'updateResponseParam') {
    after(() => getDetail(route.params.doc_id as string))
  }
})

definitionStore.$onAction(({ name, after }) => {
  // 删除全局模型
  if (name === 'deleteDefinition') {
    after(() => getDetail(route.params.doc_id as string))
  }
})

watch(
  () => route.params.doc_id,
  async () => await getDetail(route.params.doc_id as string),
  {
    immediate: true,
  }
)
</script>
