<template>
  <div class="flex h-full">
    <div class="w-246px b-r b-solid b-gray-110 px-20px">
      <IterationTree ref="iterationTreeRef" @create="onCreateOrUpdateIteration" @click="onClickFollowedProject" />
    </div>
    <div class="flex-1">
      <div class="m-auto w-776px pt-22px">
        <IterationTable
          v-if="isListMode"
          :title="currentSelectedProjectRef ? currentSelectedProjectRef.title : null"
          :project-id="currentSelectedProjectRef ? currentSelectedProjectRef.id : null"
          ref="iterationTableRef"
          @edit="onCreateOrUpdateIteration"
        />
        <IterationForm
          v-if="isFormMode"
          :iteration-id="currentEditableItreationIdRef"
          :projects="projectList"
          @success="onCreateOrUpdateIterationSuccess"
          @cancel="switchMode('list')"
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

const onCreateOrUpdateIteration = (iteration?: Iteration) => {
  currentEditableItreationIdRef.value = iteration ? iteration.id : null
  switchMode('form')
}

const onCreateOrUpdateIterationSuccess = async () => {
  await nextTick()
  iterationTableRef.value?.reload()
}

const onClickFollowedProject = (project: ProjectInfo | null) => {
  switchMode('list')
  currentSelectedProjectRef.value = project
}
</script>
