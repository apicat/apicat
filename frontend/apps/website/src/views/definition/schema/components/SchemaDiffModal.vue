<template>
  <el-dialog v-model="dialogVisible" fullscreen :show-close="false" class="hide-dialog-header no-padding" append-to-body destroy-on-close>
    <template #header="{ close }">
      <button @click="close" ref="closeBtnRef">close</button>
    </template>

    <el-button class="dialog-close" :icon="Close" @click="onCloseBtnClick" />

    <section class="ac-diff-doc" v-loading="isLoading" element-loading-background="#fff">
      <div class="ac-diff-main">
        <header>{{ leftDocTitle }}</header>
        <div class="ac-diff-content diff-left-detail" v-if="leftSchema">
          <h1 class="ac-document__title" ref="title">{{ leftSchema.name }}</h1>
          <h4>{{ leftSchema.description }}</h4>
          <div class="ac-editor mt-10px">
            <JSONSchemaEditor readonly v-model="leftSchema.schema" :definitions="definitions" />
          </div>
        </div>
      </div>

      <div class="ac-diff-drag-line" v-show="!isLoading"></div>

      <div class="ac-diff-main">
        <header>
          <el-select placeholder="对比模型" v-model="changeDocSelectRef" class="font-medium w-[260px]">
            <el-option v-for="item in historyRecordForOptions" :value="item.id" :label="item.title" :key="item.id">{{ item.title }}</el-option>
          </el-select>
        </header>

        <div class="ac-diff-content diff-right-detail" v-if="rightSchema">
          <h1 class="ac-document__title" ref="title">{{ rightSchema.name }}</h1>
          <h4>{{ rightSchema.description }}</h4>
          <div class="ac-editor mt-10px">
            <JSONSchemaEditor readonly v-model="rightSchema.schema" :definitions="definitions" />
          </div>
        </div>
      </div>
    </section>
  </el-dialog>
</template>

<script setup lang="ts">
import JSONSchemaEditor from '@/components/APIEditor/Editor.vue'
import { Close } from '@element-plus/icons-vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import useApi from '@/hooks/useApi'
import { compareSchema } from '@/api/definitionSchema'
import { useParams } from '@/hooks/useParams'
import { useDefinitionSchemaStore } from '@/store/definitionSchema'
import { useModal } from '@/hooks/useModel'

const { currentRoute } = useRouter()
const { project_id, schema_id } = useParams()
const definitionSchemaStore = useDefinitionSchemaStore()
const { definitions, historyRecordForOptions } = storeToRefs(definitionSchemaStore)
const { dialogVisible, showModel } = useModal()

const [isLoading, fetchDiffApi] = useApi(compareSchema)
const closeBtnRef = ref()
const leftDocTitle = ref('')

const leftSchema: any = ref(null)
const rightSchema: any = ref(null)

const changeDocSelectRef = ref(0)

const getDocumentDiff = async () => {
  const history_id1 = parseInt(currentRoute.value.params.history_id as any, 10)
  const selectedId = unref(changeDocSelectRef)
  const data = await fetchDiffApi({ project_id, def_id: schema_id, history_id1, history_id2: selectedId })
  leftSchema.value = data.schema1 || {}
  rightSchema.value = data.schema2 || {}
}

const show = async () => {
  showModel()
  changeDocSelectRef.value = 0
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
