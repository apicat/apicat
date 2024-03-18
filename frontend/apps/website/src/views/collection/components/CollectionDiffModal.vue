<script setup lang="ts">
import { Close } from '@element-plus/icons-vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useModal } from '@apicat/hooks'
import { parseJSONWithDefault } from '@apicat/shared'
import useApi from '@/hooks/useApi'
import { useCollectionsStore } from '@/store/collections'
import { compareCollection } from '@/api/project/collectionHistoryRecord'
import { useCollectionContext } from '@/hooks/useCollectionContext'

const AcEditor = defineAsyncComponent(() => import('@apicat/editor'))

const { currentRoute } = useRouter()
const { historyRecordForOptions } = storeToRefs(useCollectionsStore())
const { dialogVisible, showModel } = useModal()
const { responses, schemas, acEditorOptions } = useCollectionContext()

const [isLoading, compareCollectionApi] = useApi(compareCollection)
const closeBtnRef = ref()
const leftDocTitle = ref('')

const leftDoc: any = ref(null)
const rightDoc: any = ref(null)
const leftContent = ref([])
const rightContent = ref([])

const changeDocSelectRef = ref(0)

async function getDocumentDiff() {
  const projectID = currentRoute.value.params.projectID as string
  const collectionID = Number.parseInt(currentRoute.value.params.collectionID as any, 10)
  const originalID = Number.parseInt(currentRoute.value.params.historyID as any, 10)
  const targetID = unref(changeDocSelectRef)
  try {
    const res = await compareCollectionApi(projectID, collectionID, originalID, targetID)

    leftDoc.value = res?.doc1
    rightDoc.value = res?.doc2
    leftContent.value = parseJSONWithDefault(res?.doc1?.content, [])
    rightContent.value = parseJSONWithDefault(res?.doc2?.content, [])
  }
  catch (error) {
    leftDoc.value = null
    rightDoc.value = null
  }
}

async function show() {
  showModel()
  changeDocSelectRef.value = 0
  const activeDoc = historyRecordForOptions.value.find(
    item => item.id === Number.parseInt(currentRoute.value.params.historyID as any, 10),
  )
  if (activeDoc)
    leftDocTitle.value = activeDoc.title

  await getDocumentDiff()
}

function onCloseBtnClick() {
  closeBtnRef.value?.click()
}

watch(changeDocSelectRef, async () => {
  await getDocumentDiff()
})

defineExpose({
  show,
})
</script>

<template>
  <el-dialog
    v-model="dialogVisible"
    fullscreen
    :show-close="false"
    class="hide-dialog-header no-padding"
    append-to-body
    destroy-on-close
  >
    <template #header="{ close }">
      <button ref="closeBtnRef" @click="close">
        close
      </button>
    </template>

    <el-button class="dialog-close" :icon="Close" @click="onCloseBtnClick" />

    <section v-loading="isLoading" class="ac-diff-doc" element-loading-background="#fff">
      <div class="ac-diff-main scroll-content">
        <header class="font-500 text-gray-title">
          {{ leftDocTitle }}
        </header>
        <div v-if="leftDoc" class="ac-diff-content diff-left-detail">
          <h1 class="ac-document__title">
            {{ leftDoc.title }}
          </h1>
          <AcEditor
            v-if="!isLoading"
            readonly
            :content="leftContent"
            :schemas="schemas"
            :responses="responses"
            :options="acEditorOptions"
          />
        </div>
      </div>

      <div v-show="!isLoading" class="ac-diff-drag-line" />

      <div class="ac-diff-main scroll-content">
        <header>
          <el-select
            v-model="changeDocSelectRef"
            :placeholder="$t('app.project.collection.history.diff.holder')"
            class="font-medium w-[260px]"
          >
            <el-option v-for="item in historyRecordForOptions" :key="item.id" :value="item.id" :label="item.title">
              <div class="truncate max-w-260px">
                {{ item.title }}
              </div>
            </el-option>
          </el-select>
        </header>

        <div v-if="rightDoc" class="ac-diff-content diff-right-detail">
          <h1 class="ac-document__title">
            {{ rightDoc.title }}
          </h1>
          <AcEditor
            v-if="!isLoading"
            readonly
            :content="rightContent"
            :schemas="schemas"
            :responses="responses"
            :options="acEditorOptions"
          />
        </div>
      </div>
    </section>
  </el-dialog>
</template>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(diff-doc) {
  @apply flex max-h-screen overflow-hidden;
}

@include b(diff-main) {
  @apply flex-1 overflow-y-auto pt-[56px];

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
