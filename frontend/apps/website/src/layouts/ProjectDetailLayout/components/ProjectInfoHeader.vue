<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { useProjectLayoutContext } from '../composables/useProjectLayoutContext'
import { useProjectStore } from '@/store/project'
import { useNamespace } from '@/hooks/useNamespace'
import ProjectSettingModal from '@/views/project/ProjectSettingModal.vue'
import type { Menu } from '@/components/typings'
import { ProjectNavigateListEnum, getProjectNavigateList } from '@/layouts/ProjectDetailLayout/constants'
import { useIterationStore } from '@/store/iteration'

const ns = useNamespace('project-info')
const projectSettingModalRef = ref<InstanceType<typeof ProjectSettingModal>>()
const projectStore = useProjectStore()
const iterationStore = useIterationStore()
const { isIterationRoute, iterationInfo } = storeToRefs(iterationStore)
const { project, isPrivate, isGuest } = storeToRefs(projectStore)

const handleShareProject = useProjectLayoutContext('handleShareProject')

const { t } = useI18n()
const allMenus = computed(() =>
  getProjectNavigateList(false, project.value!.selfMember.permission, {
    [ProjectNavigateListEnum.ProjectShare]: {
      action: () => handleShareProject!(project.value!.id),
    },
    [ProjectNavigateListEnum.General]: {
      text: t('app.project.setting.basic.alias'),
    },
  }),
)

async function onMenuItemClick(menu: Menu) {
  if (menu.action) return await menu.action()
  unref(projectSettingModalRef)!.show(menu.key)
}
</script>

<template>
  <div :class="[ns.b(), { [ns.m('hover')]: !isGuest }]">
    <div v-if="isGuest" :class="ns.e('img')">
      <router-link to="/">
        <img src="@/assets/images/logo-square.svg" :alt="project?.title" />
      </router-link>
    </div>

    <template v-else>
      <div v-if="isIterationRoute" :class="ns.e('img')">
        <img src="@/assets/images/logo-square.svg" :alt="project?.title" />
        <router-link to="/iterations">
          <el-icon :class="ns.e('back')">
            <ac-icon-ep-arrow-left-bold />
          </el-icon>
        </router-link>
      </div>

      <el-popover v-else transition="fade" placement="bottom" width="250px" :show-arrow="false">
        <template #reference>
          <div :class="ns.e('img')">
            <img src="@/assets/images/logo-square.svg" :alt="project?.title" />
            <router-link to="/main">
              <el-icon :class="ns.e('back')">
                <ac-icon-ep-arrow-left-bold />
              </el-icon>
            </router-link>
          </div>
        </template>

        <PopperMenu :menus="allMenus" class="normal-popover-space" @menu-click="onMenuItemClick" />
      </el-popover>
    </template>

    <template v-if="isIterationRoute">
      <div class="pr-2 overflow-hidden">
        <div :title="project?.title" class="flex-y-center">
          <p class="text-base truncate">
            {{ project?.title }}
          </p>
          <el-tooltip
            v-if="isPrivate"
            effect="dark"
            placement="bottom"
            :content="$t('app.project.infoHeader.private')"
            :show-arrow="false">
            <el-icon class="ml-4px">
              <ac-icon-ep-lock />
            </el-icon>
          </el-tooltip>
        </div>
        <p class="text-sm truncate" :title="iterationInfo?.title">
          {{ iterationInfo?.title }}
        </p>
      </div>
    </template>

    <template v-else>
      <div :class="ns.e('title')" :title="project?.title">
        {{ project?.title }}
        <el-tooltip
          v-if="isPrivate"
          effect="dark"
          placement="bottom"
          :content="$t('app.project.infoHeader.private')"
          :show-arrow="false">
          <el-icon :class="ns.e('icon')">
            <ac-icon-ep-lock />
          </el-icon>
        </el-tooltip>
      </div>
    </template>
  </div>

  <ProjectSettingModal ref="projectSettingModalRef" />
</template>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;
@use '@/styles/variable' as *;

// 项目信息
@include b(project-info) {
  height: $doc-header-height;
  width: $doc-layout-left-width;
  padding: 0 $doc-layout-padding;
  @apply flex items-center fixed left-0 top-0 bg-gray-100;

  @include e(img) {
    @apply flex-none w-32px h-32px mr-10px cursor-pointer;
  }

  @include e(back) {
    @apply w-32px h-32px rounded-4px hidden bg-#dddddd text-12px;
  }

  @include e(title) {
    @apply truncate text-16px relative pr-20px;
  }

  @include e(icon) {
    @apply absolute right-0 top-50% -mt-8px;
  }

  @include m(hover) {
    @include e(img) {
      &:hover {
        img {
          @apply hidden;
        }
        @include e(back) {
          @apply inline-flex;
        }
      }
    }
  }
}
</style>
