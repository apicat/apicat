<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useNamespace } from '@apicat/hooks'

const props = withDefaults(
  defineProps<{
    title: string
    projects: ProjectAPI.ResponseProject[]
  }>(),
  {
    projects: () => [],
  },
)

const emits = defineEmits<{
  (e: 'click' | 'follow' | 'change-group', project: ProjectAPI.ResponseProject): void
}>()

const { t } = useI18n()

const ns = useNamespace('project-list')

const titleRef = computed(() => props.title || t('app.project.projects.title'))

function handleClick(project: ProjectAPI.ResponseProject) {
  return emits('click', project)
}
function handleFollowProject(project: ProjectAPI.ResponseProject) {
  return emits('follow', project)
}
function handleProjectGroup(project: ProjectAPI.ResponseProject) {
  return emits('change-group', project)
}
</script>

<template>
  <div class="flex flex-col justify-center mx-auto px-36px">
    <p class="border-b border-solid border-gray-lighter pb-30px text-18px text-gray-title">
      {{ titleRef }}
    </p>

    <div :class="ns.b()">
      <div v-for="project in projects" :key="project.id" :class="ns.e('item')" @click="handleClick(project)">
        <div class="content">
          <div
            :class="ns.e('cover')"
            :style="{
              backgroundColor: (project.cover as ProjectAPI.ProjectCover).coverBgColor,
            }"
          >
            <Iconfont class="text-white" :icon="(project.cover as ProjectAPI.ProjectCover).coverIcon" :size="55" />
          </div>
          <div :class="ns.e('title')">
            <p class="flex-1 truncate">
              {{ project.title }}
            </p>
            <div :class="ns.e('icons')">
              <el-tooltip
                :content="project.selfMember.isFollowed ? $t('app.project.stars.unstar') : $t('app.project.stars.star')"
                :show-arrow="false"
                placement="bottom"
              >
                <el-icon v-if="!project.selfMember.isFollowed" size="18" @click.stop="handleFollowProject(project)">
                  <ac-icon-mdi:star-outline />
                </el-icon>
                <el-icon v-else size="18" color="#FF9966" @click.stop="handleFollowProject(project)">
                  <ac-icon-mdi:star />
                </el-icon>
              </el-tooltip>

              <el-tooltip :content="$t('app.project.groups.grouping')" placement="bottom" :show-arrow="false">
                <Iconfont :size="18" class="ml-10px" icon="ac-fenlei" @click.stop="handleProjectGroup(project)" />
              </el-tooltip>
            </div>
          </div>
        </div>
      </div>
    </div>
    <el-empty v-if="!projects.length" :image-size="200" :description="$t('app.project.list.emptyDataTip')" />
  </div>
</template>

<style scoped lang="scss">
@use '@/styles/mixins/mixins' as *;

@include b(project-list) {
  // justify-content: space-between;
  // grid-template-columns: repeat(auto-fill, 250px);
  display: grid;
  grid-gap: 20px;
  @apply my-20px py-10px;
  // grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));

  @include e(item) {
    transition: 0.2s all;
    @apply flex justify-center items-center;
    .content {
      @apply flex  flex-col overflow-hidden rounded shadow-md cursor-pointer hover:shadow-lg w-250px h-156px;

      &:hover {
        @include e(icons) {
          visibility: visible;
        }
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
