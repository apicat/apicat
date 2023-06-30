<template>
  <el-dialog v-model="dialogVisible" append-to-body :close-on-click-modal="false" class="fullscree hide-header" align-center destroy-on-close width="1000px">
    <ModalLayout v-loading="isLoading">
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
import { getProjectNavigateList, ProjectNavigateListEnum } from '@/commons/constant'
import BaseInfoSetting from './ProjectSettingPages/BaseInfoSetting.vue'
import ProjectMemberList from './ProjectSettingPages/ProjectMemberList.vue'
import ServerUrlSetting from './ProjectSettingPages/ServerUrlSetting.vue'
import GlobalParametersSetting from './ProjectSettingPages/GlobalParametersSetting.vue'
import ProjectExportPage from './ProjectSettingPages/ProjectExportPage.vue'
import ProjectTrashPage from './ProjectSettingPages/ProjectTrashPage.vue'
import uesProjectStore from '@/store/project'
import { useParams } from '@/hooks/useParams'
import useApi from '@/hooks/useApi'

const menus = getProjectNavigateList({
  [ProjectNavigateListEnum.BaseInfoSetting]: { component: BaseInfoSetting },
  [ProjectNavigateListEnum.ProjectMemberList]: { component: ProjectMemberList },
  [ProjectNavigateListEnum.ServerUrlSetting]: { component: ServerUrlSetting },
  [ProjectNavigateListEnum.GlobalParamsSetting]: { component: GlobalParametersSetting },
  [ProjectNavigateListEnum.ProjectExport]: { component: ProjectExportPage },
  [ProjectNavigateListEnum.ProjectTrash]: { component: ProjectTrashPage },
})

const activeTab = shallowRef<{ menu: any; type: ProjectNavigateListEnum }>({
  menu: menus[ProjectNavigateListEnum.BaseInfoSetting],
  type: ProjectNavigateListEnum.BaseInfoSetting,
})

const { dialogVisible, showModel } = useModal()
const projectStore = uesProjectStore()
const [isLoading, getProjectDetailInfo] = useApi(projectStore.getProjectDetailInfo)
const { project_id } = useParams()
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

const show = async (type: ProjectNavigateListEnum) => {
  onMenuTabClick(menus[type], type)
  showModel()
  await getProjectDetailInfo(project_id as string)
}

defineExpose({
  show,
})
</script>
