<template>
  <HistoryLayout>
    <template #left>
      <SchemaHistoryRecordTree />
    </template>
  </HistoryLayout>
</template>
<script setup lang="ts">
import HistoryLayout from './HistoryLayout.vue'
import SchemaHistoryRecordTree from '@/views/definition/schema/components/SchemaHistoryRecordTree.vue'
</script>
