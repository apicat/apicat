<template>
  <div :class="[ns.b(), ns.m('hover')]">
    <el-popover placement="bottom" width="250px">
      <template #reference>
        <div :class="ns.e('img')">
          <img src="@/assets/images/logo-square.svg" :alt="projectDetailInfo?.title" />
          <router-link to="/main">
            <el-icon :class="ns.e('back')"><ac-icon-ep-arrow-left-bold /></el-icon>
          </router-link>
        </div>
      </template>

      <PopperMenu :menus="allMenus" class="clear-popover-space" @menu-click="onMenuItemClick" />
    </el-popover>

    <div :class="ns.e('title')" :title="projectDetailInfo?.title">{{ projectDetailInfo?.title }}</div>
  </div>

  <ProjectSettingModal ref="projectSettingModalRef" />
</template>
<script setup lang="ts">
import { uesProjectStore } from '@/store/project'
import { useNamespace } from '@/hooks/useNamespace'
import ProjectSettingModal from '@/views/project/ProjectSettingModal.vue'
import { Menu } from '@/components/typings'
import { getProjectNavigateList, ProjectNavigateListEnum } from '@/commons/constant'
import { storeToRefs } from 'pinia'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { useI18n } from 'vue-i18n'
import { quitProject } from '@/api/project'
import { useUserStore } from '@/store/user'
import AcIconLogout from '~icons/mdi/logout'

const ns = useNamespace('project-info')
const projectSettingModalRef = ref<InstanceType<typeof ProjectSettingModal>>()
const projectStore = uesProjectStore()
const { projectDetailInfo, isManager } = storeToRefs(projectStore)
const { t } = useI18n()

const allMenus = computed(() => {
  const menus = getProjectNavigateList()
  if (!isManager.value) {
    menus[ProjectNavigateListEnum.QuitProject] = {
      text: t('app.project.setting.quitProject'),
      elIcon: markRaw(AcIconLogout),
      action: handlerQuitProject,
    }
  }
  return menus
})

const handlerQuitProject = async () => {
  AsyncMsgBox({
    title: t('app.project.tips.quitProjectTitle'),
    content: t('app.project.tips.quitProject'),
    onOk: async () => {
      try {
        await quitProject(projectDetailInfo.value?.id! as string)
        projectStore.clearCurrentProjectInfo()
        useUserStore().goHome()
      } catch (error) {}
    },
  })
}

const onMenuItemClick = async (menu: Menu, key: ProjectNavigateListEnum) => {
  if (menu.action) {
    return await menu.action()
  }
  unref(projectSettingModalRef)!.show(key)
}
</script>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

$doc-layout-padding: 25px;
$doc-header-height: 60px;
$doc-layout-left-width: 315px;
$doc-layout-left-padding-right: 30px;
$doc-content-padding-left: $doc-layout-left-width + $doc-layout-left-padding-right;
$document-padding: 40px;

// 项目信息
@include b(project-info) {
  height: $doc-header-height;
  width: $doc-layout-left-width;
  padding: 0 $doc-layout-padding;
  @apply flex items-center fixed left-0 top-0 z-50;

  @include e(img) {
    @apply flex-none w-32px h-32px mr-10px cursor-pointer;
  }

  @include e(back) {
    @apply w-32px h-32px rounded-4px hidden bg-#dddddd text-12px;
  }

  @include e(title) {
    @apply truncate text-16px relative pr-20px;
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
