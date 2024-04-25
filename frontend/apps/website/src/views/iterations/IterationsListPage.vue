<script setup lang="ts">
import IterationTree from './components/IterationTree.vue'
import IterationTable from './components/IterationTable.vue'
import IterationForm from './components/IterationForm.vue'
import LeftRightLayout from '@/layouts/LeftRightLayout.vue'
import Page30pLayout from '@/layouts/Page30pLayout.vue'
import { usePageMode } from '@/views/composables/usePageMode'
import { useProjectList } from '@/views/composables/useProjectList'
import { Authority } from '@/commons/constant'

const iterationTableRef = ref<InstanceType<typeof IterationTable>>()
const iterationTreeRef = ref<InstanceType<typeof IterationTree>>()

const { isFormMode, isListMode, switchMode } = usePageMode()
const { projectList } = useProjectList({
  permissions: [Authority.Manage, Authority.Write],
})

const currentEditableItreationIdRef = ref<string | null>(null)
const currentSelectedProjectRef = ref<ProjectAPI.ResponseProject | null>()

async function onEditIteration(iteration: IterationAPI.ResponseIteration) {
  await onCreateIteration()
  currentEditableItreationIdRef.value = iteration.id
}

async function onCreateIteration() {
  currentEditableItreationIdRef.value = null
  await nextTick()
  switchMode('form')
}

async function onCreateOrUpdateIterationSuccess() {
  await nextTick()
  switchMode('list')

  if (!unref(currentEditableItreationIdRef)) {
    currentSelectedProjectRef.value = null
    iterationTableRef.value?.reload()
    iterationTreeRef.value?.goSelectedAll()
  }
  else {
    iterationTableRef.value?.refresh()
    iterationTreeRef.value?.goBackSelected()
  }
}

function onClickFollowedProject(project: ProjectAPI.ResponseProject | null) {
  switchMode('list')
  currentSelectedProjectRef.value = project
}

function onCancelCreateOrUpdateIteration() {
  iterationTreeRef.value?.goBackSelected()
  switchMode('list')
}
</script>

<template>
  <LeftRightLayout main-width="auto">
    <template #left>
      <IterationTree ref="iterationTreeRef" @create="onCreateIteration" @click="onClickFollowedProject" />
    </template>

    <Page30pLayout>
      <IterationTable
        v-show="isListMode"
        ref="iterationTableRef"
        :title="currentSelectedProjectRef?.title ?? null"
        :project-id="currentSelectedProjectRef?.id ?? null"
        @edit="onEditIteration"
      />

      <IterationForm
        v-if="isFormMode"
        :iteration-i-d="currentEditableItreationIdRef"
        :projects="projectList"
        @success="onCreateOrUpdateIterationSuccess"
        @cancel="onCancelCreateOrUpdateIteration"
      />
    </Page30pLayout>
  </LeftRightLayout>
</template>
