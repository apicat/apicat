<template>
  <div class="flex h-full">
    <div class="w-246px b-r b-solid b-gray-110 px-20px">
      <IterationTree ref="iterationTreeRef" @create="onCreateIteration" @click="onClickFollowedProject" />
    </div>
    <div class="flex-1">
      <div class="m-auto w-776px pt-22px">
        <IterationTable
          v-show="isListMode"
          :title="currentSelectedProjectRef ? currentSelectedProjectRef.title : null"
          :project-id="currentSelectedProjectRef ? currentSelectedProjectRef.id : null"
          ref="iterationTableRef"
          @edit="onEditIteration"
        />
        <IterationForm
          v-if="isFormMode"
          :iteration-id="currentEditableItreationIdRef"
          :projects="projectList"
          @success="onCreateOrUpdateIterationSuccess"
          @cancel="onCancelCreateOrUpdateIteration"
        />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import IterationTree from './components/IterationTree.vue'
import IterationTable from './components/IterationTable.vue'
import IterationForm from './components/IterationForm.vue'
import { usePageMode } from './logic/usePageMode'
import { useProjectList } from '../project/logic/useProjectList'
import { Iteration, ProjectInfo } from '@/typings'

const iterationTableRef = ref<InstanceType<typeof IterationTable>>()
const iterationTreeRef = ref<InstanceType<typeof IterationTree>>()
const { isFormMode, isListMode, switchMode } = usePageMode()
const { projectList } = useProjectList()

const currentEditableItreationIdRef = ref<number | string | null>()
const currentSelectedProjectRef = ref<ProjectInfo | null>()

const onEditIteration = async (iteration: Iteration) => {
  await onCreateIteration()
  currentEditableItreationIdRef.value = iteration.id
  iterationTreeRef.value?.removeSelected()
}

const onCreateIteration = async () => {
  currentEditableItreationIdRef.value = null
  await nextTick()
  switchMode('form')
}

const onCreateOrUpdateIterationSuccess = async () => {
  await nextTick()
  switchMode('list')
  iterationTableRef.value?.reload()
  iterationTreeRef.value?.goBackSelected()
}

const onClickFollowedProject = (project: ProjectInfo | null) => {
  switchMode('list')
  currentSelectedProjectRef.value = project
}

const onCancelCreateOrUpdateIteration = () => {
  iterationTreeRef.value?.goBackSelected()
  switchMode('list')
}
</script>
