<script setup lang="ts">
import { JSONSchemaTable } from '@apicat/components'
import { Close } from '@element-plus/icons-vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useModal } from '@apicat/hooks'
import useApi from '@/hooks/useApi'
import { apiDiffSchemaHistory } from '@/api/project/definition/schema'
import useDefinitionSchemaStore from '@/store/definitionSchema'

const { currentRoute } = useRouter()
const schemaStore = useDefinitionSchemaStore()
const { historyRecordForOptions } = storeToRefs(schemaStore)
const { dialogVisible, showModel } = useModal()

const [isLoading, diffSchemaHistory] = useApi(apiDiffSchemaHistory)
const closeBtnRef = ref()
const leftSchemaTitle = ref('')

const leftSchema: any = ref(null)
const rightSchema: any = ref(null)

const changeSchemaSelectRef = ref(0)

async function getSchemaDiff() {
  const projectID = currentRoute.value.params.projectID as string
  const schemaID = Number.parseInt(currentRoute.value.params.schemaID as any, 10)
  const originalID = Number.parseInt(currentRoute.value.params.historyID as any, 10)
  const targetID = unref(changeSchemaSelectRef)
  const res = await diffSchemaHistory(projectID, schemaID, originalID, targetID)
  leftSchema.value = res?.schema1
  rightSchema.value = res?.schema2
}

async function show() {
  showModel()
  changeSchemaSelectRef.value = 0
  const activeSchema = historyRecordForOptions.value.find(
    item => item.id === Number.parseInt(currentRoute.value.params.historyID as any, 10),
  )
  if (activeSchema)
    leftSchemaTitle.value = activeSchema.title

  await getSchemaDiff()
}

function onCloseBtnClick() {
  closeBtnRef.value?.click()
}

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
          {{ leftSchemaTitle }}
        </header>
        <div v-if="leftSchema" class="ac-diff-content diff-left-detail">
          <h1 class="ac-document__title">
            {{ leftSchema.title }}
          </h1>
          <JSONSchemaTable :schema="leftSchema.schema" readonly :definition-schemas="schemaStore.schemas" />
        </div>
      </div>

      <div v-show="!isLoading" class="ac-diff-drag-line" />

      <div class="ac-diff-main scroll-content">
        <header>
          <el-select
            v-model="changeSchemaSelectRef"
            :placeholder="$t('app.schema.history.diff.holder')"
            class="font-medium w-[260px]"
            @change="getSchemaDiff"
          >
            <el-option v-for="item in historyRecordForOptions" :key="item.id" :value="item.id" :label="item.title">
              {{ item.title }}
            </el-option>
          </el-select>
        </header>

        <div v-if="rightSchema" class="ac-diff-content diff-right-detail">
          <h1 class="ac-document__title">
            {{ rightSchema.title }}
          </h1>
          <JSONSchemaTable :schema="rightSchema.schema" readonly :definition-schemas="schemaStore.schemas" />
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
