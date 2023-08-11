<template>
  <div class="flex h-full">
    <div class="w-246px b-r b-solid b-gray-110 px-20px">
      <IterationTree v-model:selected-key="selectedProjectKeyRef" :projects="followedProjects" @create="switchMode('form')" @click-item="switchMode('list')" />
    </div>
    <div class="flex-1">
      <div class="m-auto w-776px pt-22px">
        <IterationTable v-if="isListMode" :data="iterations" v-model:page="currentPage" :total="total" v-loading="isLoading" @remove="handleRemoveIteration" />
        <IterationForm v-if="isFormMode" :projects="projectList" :id="null" @success="fetchIterationList" @cancel="switchMode('list')" />
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

const { isFormMode, isListMode, switchMode } = usePageMode()
const { projectList } = useProjectList()
const { followedProjects } = useFollowedProjectList()
const { isLoading, selectedProjectKeyRef, data: iterations, currentPage, total, fetchIterationList, handleRemoveIteration } = useIterationList()
</script>
