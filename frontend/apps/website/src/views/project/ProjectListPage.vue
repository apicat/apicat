<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import ProjectGroups from './components/ProjectGroups.vue'
import ProjectList from './components/ProjectList.vue'
import { useProjects } from './logic/useProjects'

import CreateOrUpdateProjectGroup from './components/CreateOrUpdateProjectGroup.vue'
import CreateProjectForm from './components/CreateProjectForm.vue'
import SelectProjectGroup from './components/SelectProjectGroup.vue'

import { useProjectListProvider } from './logic/useProjectListContext'
import { usePageMode } from '@/views/composables/usePageMode'
import LeftRightLayout from '@/layouts/LeftRightLayout.vue'
import { useTeamStore } from '@/store/team'
import useProjectGroupStore from '@/store/projectGroup'

const groupListRef = ref<InstanceType<typeof ProjectGroups>>()
const teamStore = useTeamStore()
const { createOrUpdateProjectGroupRef } = useProjectListProvider()
const { isFormMode, isListMode, switchMode } = usePageMode()
const { projectGroups } = storeToRefs(useProjectGroupStore())

const {
  isLoading,
  projects,
  selectProjectGroupRef,
  handleFollowProject,
  navigateToProjectDetail,
  showProjectGroupModal,
  refreshProjectList,
  loadPrjectListByGroupId,
} = useProjects()

// 创建项目
const createProjectFormRef = ref<InstanceType<typeof CreateProjectForm>>()
const currentRrojectGroupInfo = ref<SwitchProjectGroupInfo>()

const titleRef = computed(() => {
  if (!currentRrojectGroupInfo.value)
    return ''

  const info = currentRrojectGroupInfo.value

  if (typeof info.key === 'string')
    return info.title

  if (typeof info.key === 'number')
    return projectGroups.value.find(group => group.id === info.key)?.name || ''

  return currentRrojectGroupInfo.value?.title || ''
})

function onCreateProjectMenuClick() {
  createProjectFormRef.value?.reset()
  switchMode('form')
}

// 切换项目分组
function handleSwitchProjectGroup(info: SwitchProjectGroupInfo) {
  currentRrojectGroupInfo.value = info
  switchMode('list')
  loadPrjectListByGroupId()
}

function handleCancelCreateProject() {
  switchMode('list')
  groupListRef.value?.goBackSelected()
}
</script>

<template>
  <LeftRightLayout main-width="auto">
    <template #left>
      <ProjectGroups
        ref="groupListRef"
        @switch-group="handleSwitchProjectGroup"
        @create-project="onCreateProjectMenuClick"
        @delete="refreshProjectList"
      />
    </template>

    <ProjectList
      v-show="isListMode"
      :loading="isLoading"
      :title="titleRef"
      :projects="projects"
      @click="navigateToProjectDetail"
      @follow="handleFollowProject"
      @change-group="showProjectGroupModal"
    />

    <CreateProjectForm
      v-show="isFormMode && !teamStore.isMember"
      ref="createProjectFormRef"
      @cancel="handleCancelCreateProject"
    />
  </LeftRightLayout>

  <CreateOrUpdateProjectGroup ref="createOrUpdateProjectGroupRef" />
  <SelectProjectGroup ref="selectProjectGroupRef" @success="refreshProjectList" />
</template>
