<template>
  <div class="flex flex-col justify-center mx-auto px-36px">
    <p class="border-b border-solid border-gray-lighter pb-30px text-18px text-gray-title">
      {{ titleRef }}
    </p>

    <ul :class="ns.b()">
      <li :class="ns.e('item')" v-for="project in projects" @click="handleClick(project)" :key="project.id">
        <div :class="ns.e('cover')" :style="{ backgroundColor: (project.cover as ProjectCover).coverBgColor }">
          <Iconfont class="text-white" :icon="(project.cover as ProjectCover).coverIcon" :size="55" />
        </div>
        <div :class="ns.e('title')">
          <p class="flex-1 truncate">{{ project.title }}</p>
          <div :class="ns.e('icons')">
            <el-tooltip :content="project.is_followed ? '取消关注' : '关注项目'" placement="bottom">
              <el-icon size="18" v-if="!project.is_followed" @click.stop="handleFollowProject(project)"><ac-icon-mdi:star-outline /></el-icon>
              <el-icon size="18" v-else color="#FF9966" @click.stop="handleFollowProject(project)"><ac-icon-mdi:star /></el-icon>
            </el-tooltip>

            <el-tooltip content="项目分组" placement="bottom">
              <Iconfont :size="18" class="ml-10px" icon="ac-fenlei" @click.stop="handleProjectGroup(project)" />
            </el-tooltip>
          </div>
        </div>
      </li>
    </ul>
    <el-empty v-if="!projects.length" :image-size="200" :description="$t('app.project.tips.noData')" />
  </div>
</template>

<script setup lang="ts">
import { ProjectCover, ProjectInfo } from '@/typings'
import { useNamespace } from '@/hooks'

const emits = defineEmits<{
  (e: 'click', project: ProjectInfo): void
  (e: 'follow', project: ProjectInfo): void
  (e: 'group', project: ProjectInfo): void
}>()

const props = withDefaults(
  defineProps<{
    title: string
    projects: ProjectInfo[]
  }>(),
  {
    projects: () => [],
    title: '所有项目',
  }
)

const ns = useNamespace('project-list')

const titleRef = computed(() => props.title || '所有项目')

const handleClick = (project: ProjectInfo) => emits('click', project)
const handleFollowProject = (project: ProjectInfo) => emits('follow', project)
const handleProjectGroup = (project: ProjectInfo) => emits('group', project)
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
      @include e(icons) {
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

  @include e(icons) {
    visibility: hidden;
    @apply flex pt-1px;
  }
}
</style>
