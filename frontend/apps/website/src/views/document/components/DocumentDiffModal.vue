<template>
  <el-dialog v-model="dialogVisible" fullscreen :show-close="false" class="hide-dialog-header no-padding" append-to-body destroy-on-close>
    <template #header="{ close }">
      <button @click="close" ref="closeBtnRef">close</button>
    </template>

    <el-button class="dialog-close" :icon="Close" @click="onCloseBtnClick" />

    <section class="ac-diff-doc" v-loading="isLoading" element-loading-background="#fff">
      <div class="ac-diff-main">
        <header>{{ leftDocTitle }}</header>
        <div class="ac-diff-content diff-left-detail">
          <h1 class="ac-document__title" ref="title">{{ leftDoc.title }}</h1>
          <div class="ac-editor mt-10px" v-if="leftDoc">
            <RequestMethodRaw class="mb-10px" :doc="leftDoc" :urls="urlServers" />
            <RequestParamRaw class="mb-10px" :doc="leftDoc" :definitions="definitions" />
            <ResponseParamTabsRaw :doc="leftDoc" :definitions="definitions" :project-id="project_id" />
          </div>
        </div>
      </div>

      <div class="ac-diff-drag-line" v-show="!isLoading"></div>

      <div class="ac-diff-main">
        <header>
          <el-select placeholder="对比文档" v-model="changeDocSelectRef" class="font-medium w-[260px]">
            <el-option v-for="item in historyRecordForOptions" :value="item.id" :label="item.title" :key="item.id">{{ item.title }}</el-option>
          </el-select>
        </header>

        <div class="ac-diff-content diff-right-detail">
          <h1 class="ac-document__title" ref="title">{{ rightDoc.title }}</h1>
          <div class="ac-editor mt-10px" v-if="rightDoc">
            <RequestMethodRaw class="mb-10px" :doc="rightDoc" :urls="urlServers" />
            <RequestParamRaw class="mb-10px" :doc="rightDoc" :definitions="definitions" />
            <ResponseParamTabsRaw :doc="rightDoc" :definitions="definitions" :project-id="project_id" />
          </div>
        </div>
      </div>
    </section>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, unref, watch } from 'vue'
import { Close } from '@element-plus/icons-vue'
import { useDocumentStore } from '@/store/document'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import useApi from '@/hooks/useApi'
import { compareDocument } from '@/api/collection'
import { useParams } from '@/hooks/useParams'
import useProjectStore from '@/store/project'
import { useDefinitionSchemaStore } from '@/store/definition'
import { useModal } from '@/hooks/useModel'

const { currentRoute } = useRouter()
const { project_id, doc_id } = useParams()
const documentStore = useDocumentStore()
const projectStore = useProjectStore()
const definitionSchemaStore = useDefinitionSchemaStore()
const { historyRecordForOptions } = storeToRefs(documentStore)
const { urlServers } = storeToRefs(projectStore)
const { definitions } = storeToRefs(definitionSchemaStore)
const { dialogVisible, showModel } = useModal()

const [isLoading, fetchDiffApi] = useApi(compareDocument)
const closeBtnRef = ref()
const leftDocTitle = ref('')

const leftDoc: any = ref({})
const rightDoc: any = ref({})

const changeDocSelectRef = ref(0)

const getDocumentDiff = async () => {
  const history_id1 = parseInt(currentRoute.value.params.history_id as any, 10)
  const selectedId = unref(changeDocSelectRef)
  const data = await fetchDiffApi({ project_id, collection_id: doc_id, history_id1, history_id2: selectedId })
  leftDoc.value = data.doc1 || {}
  rightDoc.value = data.doc2 || {}
}

const show = async () => {
  showModel()
  const activeDoc = historyRecordForOptions.value.find((item) => item.id === parseInt(currentRoute.value.params.history_id as any, 10))
  if (activeDoc) {
    leftDocTitle.value = activeDoc.title
  }
  await getDocumentDiff()
}

const onCloseBtnClick = () => {
  closeBtnRef.value?.click()
}

watch(changeDocSelectRef, async () => {
  await getDocumentDiff()
})

defineExpose({
  show,
})
</script>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(diff-doc) {
  @apply flex max-h-screen overflow-hidden;
}

@include b(diff-main) {
  @apply flex-1 overflow-y-scroll pt-[56px];

  header {
    position: fixed;
    width: 100%;
    background: #ffffff;
    z-index: 99;
    top: 0;
    @apply border-b border-gray-200 h-[56px] flex items-center px-[20px] text-[16px] font-medium;
  }
}

@include b(diff-content) {
  @apply px-[30px] py-[10px];
}

@include b(diff-drag-line) {
  @apply border-l border-gray-200 fixed top-0 bottom-0 left-1/2;
  z-index: 100;
}
</style>
