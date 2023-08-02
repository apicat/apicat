<template>
  <div class="ac-header-operate">
    <p class="flex-1"></p>
    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="onSaveOrEditBtnClick" :loading="isLoading">还原此历史记录</el-button>
      <el-button @click="onShowDocumentDiffModal">对比</el-button>
    </div>
  </div>

  <!-- <DocumentDiffModal ref="documentDiffModal" /> -->
</template>

<script setup lang="tsx">
import { ref, computed, inject } from 'vue'
// import DocumentDiffModal from '@/views/document/components/DocumentDiffModal.vue'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
// import { restoreDocumentByHistoryRecord } from '@/api/document'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage as $Message } from 'element-plus'
import { storeToRefs } from 'pinia'
// import { useDocumentStore } from '@/stores/document'

defineProps({
  title: {
    type: String,
    default: '',
  },
})

// const goBack: any = inject('goBack')
// const documentStore = useDocumentStore()
const route = useRoute()
const { currentRoute } = useRouter()
// const { documentHistoryRecordTree } = storeToRefs(documentStore)
const { project_id_public } = route.params

const documentDiffModal = ref()
const isLoading = ref(false)
const isShowTitle = ref(false)

// const hasDocument = computed(() => !!currentRoute.value.params.id && documentHistoryRecordTree.value.length !== 0)

const onSaveOrEditBtnClick = () => {
  if (!currentRoute.value.params.id) {
    // $Message.error('')
    return
  }

  AsyncMsgBox({
    title: '提示',
    content: <div class="break-all">确定还原此历史记录吗？</div>,
    onOk: () => {
      // return restoreDocumentByHistoryRecord(project_id_public, currentRoute.value.params.id).then((res: any) => {
      //   $Message.success(res.msg || '还原成功！')
      //   goBack()
      // })
    },
  })
}

const onShowDocumentDiffModal = () => {
  documentDiffModal.value?.show()
}
</script>
