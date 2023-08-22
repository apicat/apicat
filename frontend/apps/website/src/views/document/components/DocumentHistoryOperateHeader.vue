<template>
  <div class="ac-header-operate" v-if="hasDocument">
    <div class="ac-header-operate__main">
      <p class="ac-header-operate__title">{{ title }}</p>
    </div>
    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="onSaveOrEditBtnClick">还原此历史记录</el-button>
      <el-button @click="onShowDocumentDiffModal">对比</el-button>
    </div>
  </div>

  <DocumentDiffModal ref="documentDiffModal" />
</template>

<script setup lang="tsx">
import { ref, computed, inject } from 'vue'
import DocumentDiffModal from '@/views/document/components/DocumentDiffModal.vue'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { restoreDocumentByHistoryRecord } from '@/api/collection'
import { useRouter } from 'vue-router'
import { ElMessage as $Message } from 'element-plus'
import { storeToRefs } from 'pinia'
import { useParams } from '@/hooks/useParams'
import { useDocumentStore } from '@/store/document'

defineProps({
  title: {
    type: String,
    default: '',
  },
})

const goBack = inject('goBack') as () => void
const documentStore = useDocumentStore()
const { documentHistoryRecordTree } = storeToRefs(documentStore)
const { currentRoute } = useRouter()
const { project_id, computedRouteParams } = useParams()

const documentDiffModal = ref<InstanceType<typeof DocumentDiffModal>>()

const hasDocument = computed(() => !!currentRoute.value.params.history_id && documentHistoryRecordTree.value.length !== 0)

const onSaveOrEditBtnClick = () => {
  const { history_id } = currentRoute.value.params
  if (!history_id) {
    return
  }

  const { doc_id: collection_id } = unref(computedRouteParams)

  AsyncMsgBox({
    title: '提示',
    content: <div class="break-all">确定还原此历史记录吗？</div>,
    onOk: () =>
      restoreDocumentByHistoryRecord({ project_id, collection_id, history_id }).then((res: any) => {
        $Message.success(res.msg || '还原成功')
        goBack()
      }),
  })
}

const onShowDocumentDiffModal = () => {
  documentDiffModal.value?.show()
}
</script>
