<script setup lang="ts">
import { ProjectNavigateListEnum, getProjectNavigateList } from '@/layouts/ProjectDetailLayout/constants'
import type { Menu } from '@/components/typings'
import { useModal } from '@/hooks'
import ModalLayout from '@/layouts/ModalLayout.vue'
import useProjectStore from '@/store/project'

const projectStore = useProjectStore()
const menus = computed(() => getProjectNavigateList(true, projectStore.project!.selfMember.permission))

const activeTab = shallowRef<{ menu: any; type: string }>({
  menu: menus.value[ProjectNavigateListEnum.General],
  type: ProjectNavigateListEnum.General,
})

const { dialogVisible, showModel } = useModal()
const isLoading = ref(false)
async function onMenuTabClick(menu: Menu, type: string) {
  if (!menu) throw new Error('ProjectSettingModal active menu is null')

  if (menu.action) {
    await menu.action()
    return
  }

  activeTab.value = {
    type,
    menu,
  }
}

async function show(type: ProjectNavigateListEnum) {
  onMenuTabClick(menus.value[type], type)
  showModel()
}

defineExpose({
  show,
})
</script>

<template>
  <el-dialog
    v-model="dialogVisible"
    :close-on-click-modal="false"
    class="fullscree hide-header"
    append-to-body
    align-center
    destroy-on-close
    width="1200px">
    <ModalLayout v-loading="isLoading">
      <template #nav>
        <p class="text-16px text-gray-950 font-500">
          {{ $t('app.project.setting.title') }}
        </p>
        <div class="grid-container mt-20px">
          <template v-for="(menu, type) in menus" :key="menu.icon">
            <div class="box">
              <Iconfont :icon="menu.icon" />
            </div>
            <div
              class="cursor-pointer box pl-10px line-height-40px"
              :class="{ 'text-blue-primary': activeTab.type === type }"
              @click="onMenuTabClick(menu, type as string)">
              {{ menu.text }}
            </div>
          </template>
        </div>
      </template>
      <template v-if="activeTab" #title>
        {{ activeTab.menu.detailTitle ?? activeTab.menu.text }}
      </template>
      <component :is="activeTab.menu.component" v-if="activeTab && activeTab.menu.component" />
    </ModalLayout>
  </el-dialog>
</template>

<style scoped>
.grid-container {
  display: grid;
  /* grid-template-columns: repeat(2, 1fr); */
  grid-template-columns: max-content 1fr;
  /* grid-template-rows: repeat(2, 1fr); */
  align-items: center;
  justify-content: center;
  /* gap: 10px; */
}

.box {
  display: flex;
  user-select: none;
}
</style>
