<template>
  <HistoryLayout>
    <template #back>
      <div :class="historyInfo.b()">
        <div :class="historyInfo.e('img')">
          <a href="javascript:void(0)" @click="handleGoBack">
            <el-icon :class="historyInfo.e('back')"><ac-icon-ep-arrow-left-bold /></el-icon>
          </a>
        </div>
        <div :class="historyInfo.e('title')">历史记录</div>
      </div>
    </template>

    <template #left>
      <SchemaHistoryRecordTree />
    </template>
  </HistoryLayout>
</template>
<script setup lang="ts">
import HistoryLayout from './HistoryLayout.vue'
import SchemaHistoryRecordTree from '@/views/definition/schema/components/SchemaHistoryRecordTree.vue'
import { useNamespace } from '@/hooks/useNamespace'
import { useGoPage } from '@/hooks/useGoPage'

const historyInfo = useNamespace('history-info')
const { goSchemaDetailPage } = useGoPage()

const handleGoBack = () => {
  goSchemaDetailPage()
}

provide('goBack', handleGoBack)
</script>
