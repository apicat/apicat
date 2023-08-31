<template>
  <LeftRightLayout main-width="auto">
    <template #left>
      <ProjectGroups
        ref="groupListRef"
        v-model:selected="selectedGroupRef"
        :groups="projectGroups"
        @switch-group="onSwitchProjectGroup"
        @create-project="onCreateProject"
        @create-group="handleRenameProjectGroup"
        @delete-group="handleDeleteProjectGroup"
        @rename-group="handleRenameProjectGroup"
        @sort-group="handleSortProjectGroup"
      />
    </template>

    <ProjectList
      ref="iterationTableRef"
      v-show="isListMode"
      v-loading="isLoading"
      :title="titleRef"
      :projects="projects"
      @click="goProjectDetail"
      @follow="handleFollowProject"
      @change-group="changeProjectGroup"
    />

    <CreateProjectForm v-if="isFormMode" :groups="groupsForOptions" @cancel="onCancel" :group_id="selectedGroupKeyForCreateForm" />
  </LeftRightLayout>

  <CreateOrUpdateProjectGroup ref="createOrUpdateProjectGroupRef" @success="refreshProjectGroups" />

  <SelectProjectGroup ref="selectProjectGroupRef" @success="refreshProjectList" />
</template>

<script lang="ts" setup>
import LeftRightLayout from '@/layouts/LeftRightLayout.vue'
import ProjectGroups from './components/ProjectGroups.vue'
import ProjectList from './components/ProjectList.vue'
import { usePageMode } from '@/views/composables/usePageMode'
import { SwitchProjectGroupInfo } from '@/typings'
import { useProjects } from './logic/useProjects'
import { useProjectGroups } from './logic/useProjectGroups'
import CreateOrUpdateProjectGroup from './components/CreateOrUpdateProjectGroup.vue'
import CreateProjectForm from './components/CreateProjectForm.vue'
import SelectProjectGroup from './components/SelectProjectGroup.vue'

const titleRef = ref('')
const groupListRef = ref<InstanceType<typeof ProjectGroups>>()
const { isFormMode, isListMode, switchMode } = usePageMode()

const {
  selectedGroupRef,
  createOrUpdateProjectGroupRef,
  projectGroups,
  groupsForOptions,
  handleDeleteProjectGroup,
  handleRenameProjectGroup,
  handleSortProjectGroup,
  refreshProjectGroups,
} = useProjectGroups()

const { isLoading, projects, selectProjectGroupRef, handleFollowProject, goProjectDetail, changeProjectGroup, refreshProjectList } = useProjects(selectedGroupRef)

const selectedGroupKeyForCreateForm = computed<number>(() => (typeof selectedGroupRef.value !== 'number' ? 0 : selectedGroupRef.value))

// 创建项目
const onCreateProject = () => switchMode('form')

// 切换项目分组
const onSwitchProjectGroup = ({ title }: SwitchProjectGroupInfo) => {
  titleRef.value = title
  switchMode('list')
}

const onCancel = () => {
  switchMode('list')
  groupListRef.value?.goBackSelected()
}
</script>
