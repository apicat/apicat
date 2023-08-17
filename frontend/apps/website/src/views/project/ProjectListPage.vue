<template>
  <div class="container flex flex-col justify-center mx-auto px-20px" v-loading="isLoading">
    <p class="border-b border-solid border-gray-lighter py-20px mt-20px text-20px">
      {{ $t('app.project.list.tabTitle') }}
    </p>

    <ul :class="ns.b()">
      <li
        v-if="!isNormalUser"
        class="flex flex-col justify-between rounded cursor-pointer w-250px h-156px hover:shadow-lg bg-gray-110 px-20px py-16px"
        @click="handleShowModelClick"
      >
        <ac-icon-ep-plus class="text-18px" />
        <p>{{ $t('app.project.createModal.title') }}</p>
      </li>

      <li :class="ns.e('item')" v-for="project in projectList" @click="$router.push(getProjectDetailPath(project.id))">
        <div :class="ns.e('cover')" :style="{ backgroundColor: (project.cover as ProjectCover).coverBgColor }">
          <Iconfont class="text-white" :icon="(project.cover as ProjectCover).coverIcon" :size="55" />
        </div>
        <div :class="ns.e('title')">
          <p class="flex-1 truncate">{{ project.title }}</p>
          <el-tooltip :content="project.is_followed ? '取消关注' : '关注项目'" placement="bottom">
            <div :class="ns.e('follow')">
              <el-icon size="16" v-if="!project.is_followed" @click.stop="handleFollowProject(project)"><ac-icon-ph:star-light /></el-icon>
              <el-icon size="16" v-else color="#FF9966" @click.stop="handleFollowProject(project)"><ac-icon-ph:star-fill /></el-icon>
            </div>
          </el-tooltip>
        </div>
      </li>
    </ul>
    <el-empty v-if="isNormalUser && !projectList.length" :image-size="200" :description="$t('app.project.tips.noData')" />
  </div>

  <CreateProjectModal ref="createProjectModal" />
</template>

<script lang="ts" setup>
import { getProjectDetailPath } from '@/router'
import CreateProjectModal from './CreateProjectModal.vue'
import { ProjectCover } from '@/typings'
import { useUserStore } from '@/store/user'
import { useProjectList } from './logic/useProjectList'
import { useNamespace } from '@/hooks'

const ns = useNamespace('project-list')
const createProjectModal = ref<InstanceType<typeof CreateProjectModal>>()
const { isNormalUser } = useUserStore()
const { isLoading, projectList, handleFollowProject } = useProjectList()

const handleShowModelClick = () => {
  createProjectModal.value!.show()
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
