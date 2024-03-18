<script setup lang="tsx">
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import CollectionDiffModal from './CollectionDiffModal.vue'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { useHistoryLayoutContext } from '@/layouts/useHistoryLayoutContext'
import { restoreCollectionHistoryRecord } from '@/api/project/collectionHistoryRecord'

const props = defineProps<{
  title: string
  projectId: string
  collectionId: number
  historyId: number
}>()

const { t } = useI18n()
const { goBack } = useHistoryLayoutContext()
const collectionDiffModal = ref<InstanceType<typeof CollectionDiffModal>>()

function onSaveOrEditBtnClick() {
  const { historyId, collectionId, projectId } = props
  if (!historyId)
    return

  AsyncMsgBox({
    title: t('app.project.collection.history.restore.poptitle'),
    content: <p>{t('app.project.collection.history.restore.popcontent')}</p>,
    onOk: async () => {
      await restoreCollectionHistoryRecord(projectId, collectionId, historyId)
      ElMessage.success(t('app.project.collection.history.restore.success'))
      goBack()
    },
  })
}

function onShowDocumentDiffModal() {
  collectionDiffModal.value?.show()
}
</script>

<template>
  <div class="ac-header-operate">
    <div class="ac-header-operate__main">
      <p class="ac-header-operate__title">
        {{ title }}
      </p>
    </div>
    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="onSaveOrEditBtnClick">
        {{ $t('app.project.collection.history.restore.title') }}
      </el-button>
      <el-button @click="onShowDocumentDiffModal">
        {{ $t('app.project.collection.history.diff.title') }}
      </el-button>
    </div>
  </div>

  <CollectionDiffModal ref="collectionDiffModal" />
</template>
