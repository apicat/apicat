<template>
  <HistoryLayout>
    <template #back>
      <div :class="historyInfo.b()">
        <div :class="historyInfo.e('img')">
          <a href="javascript:void(0)" @click="handleGoBack">
            <el-icon :class="historyInfo.e('back')"><ac-icon-ep-arrow-left-bold /></el-icon>
          </a>
        </div>
        <div :class="historyInfo.e('title')" title="projectDetailInfo?.title">历史记录</div>
      </div>
    </template>

    <template #left>
      <DocumentHistoryRecordTree />
    </template>
  </HistoryLayout>
</template>
<script setup lang="ts">
import HistoryLayout from './HistoryLayout.vue'
import DocumentHistoryRecordTree from '@/views/document/components/DocumentHistoryRecordTree.vue'
import DocumentHistoryOperateHeader from '@/views/document/components/DocumentHistoryOperateHeader.vue'
import { useNamespace } from '@/hooks/useNamespace'
import { useGoPage } from '@/hooks/useGoPage'

const historyInfo = useNamespace('history-info')
const { goDocumentDetailPage } = useGoPage()

const handleGoBack = () => {
  goDocumentDetailPage()
}

provide('goBack', handleGoBack)
provide('setTitle', handleGoBack)
</script>
