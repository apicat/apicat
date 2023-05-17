<template>
  <div class="ac-header-operate">
    <div class="ac-header-operate__main">
      <p class="flex-y-center">
        <el-icon :size="18" class="mt-1px mr-4px"><ac-icon-ic-sharp-cloud-queue /></el-icon>
        {{ isSaving ? $t('app.common.saving') : $t('app.common.savedCloud') }}
      </p>
    </div>
    <div class="ac-header-operate__btns">
      <el-button type="primary" :loading="isLoadingForSaveBtn" @click="handleSave">{{ $t('app.common.preview') }}</el-button>
    </div>
  </div>
  <div v-loading="isLoading" :class="{ 'h-50vh': !httpDoc }">
    <HttpDocumentEditor v-if="httpDoc" v-model="httpDoc" />
  </div>
</template>

<script setup lang="ts">
import { APICatCommonResponse, HttpDocument } from '@/typings'
import { getCollectionDetail, updateCollection } from '@/api/collection'
import HttpDocumentEditor from './components/HttpDocumentEditor.vue'
import { useParams } from '@/hooks/useParams'
import { useGoPage } from '@/hooks/useGoPage'
import { cloneDeep, debounce, isEmpty } from 'lodash-es'
import useApi from '@/hooks/useApi'
import { ElMessage } from 'element-plus'
import uesGlobalParametersStore from '@/store/globalParameters'
import useDefinitionStore from '@/store/definition'
import { useI18n } from 'vue-i18n'
import { DOCUMENT_EDIT_NAME } from '@/router'
import { HTTP_RESPONSE_NODE_KEY } from './components/createHttpDocument'
import useDefinitionResponseStore from '@/store/definitionResponse'

const { t } = useI18n()
const { project_id } = useParams()
const route = useRoute()
const router = useRouter()
const globalParametersStore = uesGlobalParametersStore()
const definitionStore = useDefinitionStore()
const definitionResponseStore = useDefinitionResponseStore()

const [isLoading, getCollectionDetailApi] = getCollectionDetail()
const [isLoadingForSaveBtn, updateCollectionApiWithLoading] = useApi(updateCollection)

const { goDocumentDetailPage } = useGoPage()

const isSaving = ref(false)
const httpDoc: Ref<HttpDocument | null> = ref(null)

const directoryTree: any = inject('directoryTree')

const validResponseName = (responses: APICatCommonResponse[]) => {
  let len = responses.length

  for (let i = 0; i < len; i++) {
    const item = responses[i]
    if (!item.$ref && isEmpty(item.name)) {
      ElMessage.error(t('app.response.rules.name'))
      return false
    }
  }
  return true
}

const stringifyHttpDoc = (doc: any) => {
  const data: any = cloneDeep(unref(doc))
  const responseNode = data.content.find((node: any) => node.type === HTTP_RESPONSE_NODE_KEY)
  responseNode.attrs.list = responseNode.attrs.list.map((item: any) => {
    if (item.$ref) {
      delete item.name
    }
    return item
  })
  data.content = JSON.stringify(data.content)
  const { id: collection_id, ...rest } = data
  return { project_id, collection_id, ...rest }
}

const isInvalidId = () => isNaN(parseInt(route.params.doc_id as string, 10))

const handleSave = async () => {
  await updateCollectionApiWithLoading(stringifyHttpDoc(httpDoc))
  goDocumentDetailPage()
}

const getDetail = async () => {
  // id 无效
  if (isInvalidId() || router.currentRoute.value.name !== DOCUMENT_EDIT_NAME) {
    return
  }

  httpDoc.value = null
  try {
    httpDoc.value = await getCollectionDetailApi({ project_id, collection_id: route.params.doc_id })
  } catch (error) {
    console.error(error)
  }
}

globalParametersStore.$onAction(({ name, after }) => {
  if (name === 'deleteGlobalParameter') {
    after(() => getDetail())
  }
})

definitionStore.$onAction(({ name, after }) => {
  if (name === 'deleteDefinition') {
    after(() => getDetail())
  }
})

definitionResponseStore.$onAction(({ name, after }) => {
  if (name === 'deleteDefinition') {
    after(() => getDetail())
  }
})

watch(
  httpDoc,
  debounce(async (newVal, oldVal) => {
    if (!oldVal || !oldVal.id || isInvalidId()) {
      // id 不存在
      return
    }

    if (isEmpty(newVal.title)) {
      ElMessage.error(t('app.interface.form.title'))
      return
    }

    const responsesNode = unref(httpDoc)?.content.find((node) => node.type === HTTP_RESPONSE_NODE_KEY)

    if (!validResponseName(responsesNode.attrs.list || [])) {
      return
    }

    isSaving.value = true
    try {
      await updateCollection(stringifyHttpDoc(httpDoc))
      directoryTree.updateTitle(newVal.id, newVal.title)
    } finally {
      isSaving.value = false
    }
  }, 300),
  {
    deep: true,
  }
)

watch(
  () => route.params.doc_id,
  async () => await getDetail(),
  { immediate: true }
)
</script>
