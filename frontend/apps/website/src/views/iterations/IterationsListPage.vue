<template>
  <div class="flex h-full">
    <div class="w-246px b-r b-solid b-gray-110 px-20px">
      <IterationTree v-model:selected-key="selectedProjectKeyRef" :projects="followedProjects" @create="onCreateBtnClick" @click-item="switchMode('list')" />
    </div>
    <div class="flex-1">
      <div class="m-auto w-776px pt-22px">
        <IterationTable
          v-if="isListMode"
          :data="iterations"
          v-model:page="currentPage"
          :total="total"
          v-loading="isLoading"
          @remove="handleRemoveIteration"
          @edit="onEditBtnClick"
        />
        <IterationForm v-if="isFormMode" :projects="projectList" :id="editableItreationIdRef" @success="fetchIterationList" @cancel="switchMode('list')" />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import IterationTree from './components/IterationTree.vue'
import IterationTable from './components/IterationTable.vue'
import IterationForm from './components/IterationForm.vue'
import { usePageMode } from './logic/usePageMode'
import { useFollowedProjectList } from './logic/useFollowedProjectList'
import { useIterationList } from './logic/useIterationList'
import { useProjectList } from '../project/logic/useProjectList'
import { Iteration } from '@/typings'

const { isFormMode, isListMode, switchMode } = usePageMode()
const { projectList } = useProjectList()
const { followedProjects, selectedProjectKeyRef } = useFollowedProjectList()
const { isLoading, data: iterations, currentPage, total, editableItreationIdRef, fetchIterationList, handleRemoveIteration } = useIterationList(selectedProjectKeyRef)

const onCreateBtnClick = () => {
  editableItreationIdRef.value = null
  switchMode('form')
}

const onEditBtnClick = (i: Iteration) => {
  editableItreationIdRef.value = i.id
  switchMode('form')
}
</script>
