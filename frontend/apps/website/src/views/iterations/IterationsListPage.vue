<template>
  <LeftRightLayout>
    <template #left>
      <IterationTree ref="iterationTreeRef" @create="onCreateIteration" @click="onClickFollowedProject" />
    </template>

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
  </LeftRightLayout>
</template>
<script setup lang="ts">
import LeftRightLayout from '@/layouts/LeftRightLayout.vue'
import IterationTree from './components/IterationTree.vue'
import IterationTable from './components/IterationTable.vue'
import IterationForm from './components/IterationForm.vue'
import { usePageMode } from '@/views/composables/usePageMode'
import { useProjectList } from '@/views/composables/useProjectList'
import { Iteration, ProjectInfo } from '@/typings'
import { MemberAuthorityInProject } from '@/typings/member'

const iterationTableRef = ref<InstanceType<typeof IterationTable>>()
const iterationTreeRef = ref<InstanceType<typeof IterationTree>>()
const { isFormMode, isListMode, switchMode } = usePageMode()
const { projectList } = useProjectList({ auth: [MemberAuthorityInProject.MANAGER, MemberAuthorityInProject.WRITE] })

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

  if (!unref(currentEditableItreationIdRef)) {
    currentSelectedProjectRef.value = null
    iterationTableRef.value?.reload()
    iterationTreeRef.value?.goSelectedAll()
  } else {
    iterationTableRef.value?.refresh()
    iterationTreeRef.value?.goBackSelected()
  }
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
