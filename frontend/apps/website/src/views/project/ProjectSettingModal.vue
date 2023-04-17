<template>
  <el-dialog v-model="dialogVisible" fullscreen append-to-body :close-on-click-modal="false" :close-on-press-escape="false" class="fullscree hide-header" destroy-on-close center>
    <ProjectSettingLayout>
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
      <div v-if="activeTab">
        <p class="text-16px text-gray-950 font-500 mb-20px">{{ activeTab.menu.text }}</p>
        <component v-if="activeTab.menu.component" :is="activeTab.menu.component" />
      </div>
    </ProjectSettingLayout>
  </el-dialog>
</template>
<script setup lang="ts">
import ProjectSettingLayout from '@/layouts/ProjectSettingLayout.vue'
import { useModal } from '@/hooks'
import { Menu } from '@/components/typings'
import { ProjectNavigateObject } from '@/typings/project'
import { getProjectNavigateList, ProjectNavigateListEnum } from '@/commons/constant'
import BaseInfoSetting from './ProjectSettingPages/BaseInfoSetting.vue'
import ServerUrlSetting from './ProjectSettingPages/ServerUrlSetting.vue'
import CommonRequestSetting from './ProjectSettingPages/CommonRequestSetting.vue'
import CommonResponseSetting from './ProjectSettingPages/CommonResponseSetting.vue'
import ProjectExportPage from './ProjectSettingPages/ProjectExportPage.vue'
import ProjectTrashPage from './ProjectSettingPages/ProjectTrashPage.vue'

const menus: ProjectNavigateObject = getProjectNavigateList({
  [ProjectNavigateListEnum.BaseInfoSetting]: { component: BaseInfoSetting },
  [ProjectNavigateListEnum.ServerUrlSetting]: { component: ServerUrlSetting },
  [ProjectNavigateListEnum.RequestParamsSetting]: { component: CommonRequestSetting },
  [ProjectNavigateListEnum.ResponseParamsSetting]: { component: CommonResponseSetting },
  [ProjectNavigateListEnum.ProjectExport]: { component: ProjectExportPage },
  [ProjectNavigateListEnum.ProjectTrash]: { component: ProjectTrashPage },
})

const activeTab = shallowRef<{ menu: any; type: ProjectNavigateListEnum }>({ menu: menus[ProjectNavigateListEnum.BaseInfoSetting], type: ProjectNavigateListEnum.BaseInfoSetting })

const { dialogVisible, showModel } = useModal()

const onMenuTabClick = (menu: Menu, type: ProjectNavigateListEnum) => {
  activeTab.value = {
    type,
    menu,
  }
}

const show = (type: ProjectNavigateListEnum) => {
  onMenuTabClick(menus[type], type)
  showModel()
}

defineExpose({
  show,
})
</script>
