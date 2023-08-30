<template>
  <LeftRightLayout main-width="auto">
    <template #left>
      <ProjectGroups
        ref="groupListRef"
        v-model:selected="selectedGroupRef"
        :groups="projectGroups"
        @change-group="onSwitchProjectGroup"
        @create-project="onCreateProject"
        @create-group="onCreateProjectGroup"
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
      @group="changeProjectGroup"
    />

    <div v-if="isFormMode">
      <el-button type="primary" @click="onCancel">取消</el-button>
    </div>
  </LeftRightLayout>
</template>

<script lang="ts" setup>
import LeftRightLayout from '@/layouts/LeftRightLayout.vue'
import ProjectGroups from './components/ProjectGroups.vue'
import ProjectList from './components/ProjectList.vue'
import { usePageMode } from '@/views/composables/usePageMode'

import { SwitchProjectGroupInfo } from '@/typings'
import { useProjects } from './logic/useProjects'
import { useProjectGroups } from './logic/useProjectGroups'

const titleRef = ref('')
const groupListRef = ref<InstanceType<typeof ProjectGroups>>()
const { isFormMode, isListMode, switchMode } = usePageMode()
const { selectedGroupRef, projectGroups, handleDeleteProjectGroup, handleRenameProjectGroup, handleSortProjectGroup } = useProjectGroups()
const { isLoading, projects, handleFollowProject, goProjectDetail, refreshProjectList } = useProjects(selectedGroupRef)

// 创建项目
const onCreateProject = () => switchMode('form')

// 切换项目分组
const onSwitchProjectGroup = ({ title }: SwitchProjectGroupInfo) => {
  titleRef.value = title
  switchMode('list')
}

// 创建项目分组
const onCreateProjectGroup = () => {
  console.log('创建项目分组')
}

// 调整项目分组
const changeProjectGroup = () => {
  console.log('调整项目分组', selectedGroupRef.value)
}

const onCancel = () => {
  switchMode('list')
  groupListRef.value?.goBackSelected()
}
</script>
