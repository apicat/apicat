<template>
  <HistoryLayout>
    <template #header>
      <DocumentHistoryOperateHeader />
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
</script>
