<template>
  <div class="ac-header-operate">
    <div class="ac-header-operate__main">
      <p>
        <i class="ac-iconfont"></i>
        {{ isSaving ? '保存中...' : '已保存在云端' }}
      </p>
    </div>
    <div class="ac-header-operate__btns">
      <el-button type="primary" :loading="isLoadingForSaveBtn" @click="handleSave">预览</el-button>
    </div>
  </div>
  <div v-loading="isLoading" :class="{ 'h-50vh': !httpDoc }">
    <HttpDocumentEditor v-if="httpDoc" v-model="httpDoc" />
  </div>
</template>

<script setup lang="ts">
import { HttpDocument } from '@/typings'
import { getCollectionDetail, updateCollection } from '@/api/collection'
import HttpDocumentEditor from './components/HttpDocumentEditor.vue'
import { useParams } from '@/hooks/useParams'
import { useGoPage } from '@/hooks/useGoPage'
import { debounce, isEmpty } from 'lodash-es'
import useApi from '@/hooks/useApi'
import { ElMessage } from 'element-plus'
import uesGlobalParametersStore from '@/store/globalParameters'
import useDefinitionStore from '@/store/definition'

const { project_id } = useParams()
const route = useRoute()
const globalParametersStore = uesGlobalParametersStore()
const definitionStore = useDefinitionStore()

const [isLoading, getCollectionDetailApi] = getCollectionDetail()
const [isLoadingForSaveBtn, updateCollectionApiWithLoading] = useApi(updateCollection)()

const { goDocumentDetailPage } = useGoPage()

const isSaving = ref(false)
const httpDoc: Ref<HttpDocument | null> = ref(null)

const directoryTree: any = inject('directoryTree')

const stringifyHttpDoc = (doc: any) => {
  const data: any = { ...unref(doc) }
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
  if (isInvalidId()) {
    // ElMessage.error('文档id无效')
    return
  }

  try {
    httpDoc.value = await getCollectionDetailApi({ project_id, collection_id: route.params.doc_id })
  } catch (error) {
    console.error(error)
  }
}

globalParametersStore.$onAction(({ name, after }) => {
  // 删除全局参数
  if (name === 'deleteGlobalParameter') {
    after(() => getDetail())
  }
})

definitionStore.$onAction(({ name, after }) => {
  // 删除全局模型
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
      ElMessage.error('请输入文档标题')
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
