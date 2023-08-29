<template>
  <LeftRightLayout main-width="auto">
    <template #left>
      <ProjectGroupList ref="groupListRef" @click="onSwitchProjectGroup" @create-project="onCreateProject" @create-group="onCreateProjectGroup" />
    </template>
    <ProjectList v-show="isListMode" title="呃呃沙发" ref="iterationTableRef" :projects="[]" />
  </LeftRightLayout>
</template>

<script lang="ts" setup>
import LeftRightLayout from '@/layouts/LeftRightLayout.vue'
import ProjectGroupList from './components/ProjectGroupList.vue'
import ProjectList from './components/ProjectList.vue'
import { usePageMode } from '@/views/composables/usePageMode'

import { getProjectDetailPath } from '@/router'
import { ProjectCover, ProjectGroupSelectKey } from '@/typings'
import { useUserStore } from '@/store/user'
import { useProjects } from './logic/useProjects'
import { useNamespace } from '@/hooks'

const ns = useNamespace('project-list')
const { isNormalUser } = useUserStore()
const { isLoading, projects, handleFollowProject } = useProjects()
const groupListRef = ref<InstanceType<typeof ProjectGroupList>>()
const { isFormMode, isListMode, switchMode } = usePageMode()

// 创建项目
const onCreateProject = () => {
  console.log('创建项目')
}

// 切换项目分组
const onSwitchProjectGroup = (key: ProjectGroupSelectKey) => {
  console.log('切换项目分组', key)
}

// 创建项目分组
const onCreateProjectGroup = () => {
  console.log('创建项目分组')
}
</script>
<style scoped lang="scss">
@use '@/styles/mixins/mixins' as *;

@include b(project-list) {
  display: grid;
  justify-content: space-between;
  grid-template-columns: repeat(auto-fill, 250px);
  grid-gap: 20px;
  @apply my-20px py-10px;

  @include e(item) {
    @apply flex flex-col overflow-hidden rounded shadow-md cursor-pointer hover:shadow-lg w-250px h-156px;

    &:hover {
      @include e(follow) {
        visibility: visible;
      }
    }
  }

  @include e(cover) {
    @apply flex items-center justify-center h-112px;
  }

  @include e(title) {
    @apply flex items-center flex-1 px-16px;
  }

  @include e(follow) {
    visibility: hidden;
    @apply flex pt-1px;
  }
}
</style>
