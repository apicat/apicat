<template>
  <el-dialog v-model="dialogVisible" append-to-body :close-on-click-modal="false" :close-on-press-escape="false" class="fullscree hide-header" destroy-on-close center width="70%">
    <ModalLayout>
      <template #nav>
        <p class="text-16px text-gray-950 font-500">{{ $t('app.project.setting.title') }}</p>
        <ul class="mt-20px">
          <li
            v-for="(menu, type) in menus"
            :text="menu.text"
            @click="onMenuTabClick(menu, type)"
            class="cursor-pointer py-10px"
            :class="{ 'text-blue-primary': activeTab.type === type }"
          >
            <Iconfont :icon="menu.icon" />
            {{ menu.text }}
          </li>
        </ul>
      </template>
      <template v-if="activeTab" #title>{{ activeTab.menu.text }}</template>
      <component v-if="activeTab && activeTab.menu.component" :is="activeTab.menu.component" />
    </ModalLayout>
  </el-dialog>
</template>
<script setup lang="ts">
import ModalLayout from '@/layouts/ModalLayout.vue'
import { useModal } from '@/hooks'
import { Menu } from '@/components/typings'
import { ProjectNavigateObject } from '@/typings/project'
import { getProjectNavigateList, ProjectNavigateListEnum } from '@/commons/constant'
import BaseInfoSetting from './ProjectSettingPages/BaseInfoSetting.vue'
import ProjectMemberList from './ProjectSettingPages/ProjectMemberList.vue'
import ServerUrlSetting from './ProjectSettingPages/ServerUrlSetting.vue'
import GlobalParametersSetting from './ProjectSettingPages/GlobalParametersSetting.vue'
import ProjectExportPage from './ProjectSettingPages/ProjectExportPage.vue'
import ProjectTrashPage from './ProjectSettingPages/ProjectTrashPage.vue'
import uesProjectStore from '@/store/project'
import { storeToRefs } from 'pinia'

const projectStore = uesProjectStore()
const { isManager } = storeToRefs(projectStore)

const menus: ComputedRef<ProjectNavigateObject> = computed(() => {
  const menus = getProjectNavigateList({
    [ProjectNavigateListEnum.BaseInfoSetting]: { component: BaseInfoSetting },
    [ProjectNavigateListEnum.ProjectMemberList]: { component: ProjectMemberList },
    [ProjectNavigateListEnum.ServerUrlSetting]: { component: ServerUrlSetting },
    [ProjectNavigateListEnum.GlobalParamsSetting]: { component: GlobalParametersSetting },
    [ProjectNavigateListEnum.ProjectExport]: { component: ProjectExportPage },
    [ProjectNavigateListEnum.ProjectTrash]: { component: ProjectTrashPage },
  })

  if (!isManager.value) {
    menus[ProjectNavigateListEnum.QuitProject] = {
      text: '退出项目',
      type: ProjectNavigateListEnum.QuitProject,
      action: handlerQuitProject,
    }
  }

  return menus
})

const handlerQuitProject = async () => {
  console.log('退出项目')
}

const activeTab = shallowRef<{ menu: any; type: ProjectNavigateListEnum }>({
  menu: menus.value[ProjectNavigateListEnum.BaseInfoSetting],
  type: ProjectNavigateListEnum.BaseInfoSetting,
})

const { dialogVisible, showModel } = useModal()

const onMenuTabClick = async (menu: Menu, type: ProjectNavigateListEnum) => {
  if (menu.action) {
    await menu.action()
    return
  }

  activeTab.value = {
    type,
    menu,
  }
}

const show = (type: ProjectNavigateListEnum) => {
  onMenuTabClick(menus.value[type], type)
  showModel()
}

defineExpose({
  show,
})
</script>
