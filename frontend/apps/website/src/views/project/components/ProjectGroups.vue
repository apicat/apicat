<script setup lang="tsx">
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { VueDraggableNext as Draggable } from 'vue-draggable-next'
import { Icon } from '@iconify/vue'
import { useNamespace } from '@apicat/hooks'
import { useProjectListContext } from '../logic/useProjectListContext'
import AcIconAdd from '~icons/fe/plus'
import AcIconMyProject from '~icons/mdi/folder-outline'
import AcIconMyFollowedProject from '~icons/clarity/star-line'
import { usePopover } from '@/hooks/usePopover'
import useProjectGroupStore from '@/store/projectGroup'
import { useTeamStore } from '@/store/team'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'

const emits = defineEmits<{
  (event: 'createProject' | 'delete'): void
  (event: 'switchGroup', info: SwitchProjectGroupInfo): void
}>()

let currentModifyProjectGroup: ProjectAPI.ResponseGroup | null = null

const { t } = useI18n()
const ns = useNamespace('group-list')
const teamStore = useTeamStore()
const groupStore = useProjectGroupStore()
const { selectedGroupRef, projectGroups } = storeToRefs(groupStore)
const { createOrUpdateProjectGroupRef } = useProjectListContext()
// 仅用于UI显示,处理回退时，创建项目选项被选中问题
const selectedGroupKeyForUI = ref<ProjectGroupSelectKey>(selectedGroupRef.value)

const {
  isShow: isShowPopoverMenu,
  popoverRefEl,
  showPopover,
  hidePopover,
} = usePopover({
  onHide: () => {
    currentModifyProjectGroup = null
  },
})

const activeClass = (key: ProjectGroupSelectKey) => (selectedGroupKeyForUI.value === key ? 'active' : '')

let selectedHistory: ProjectGroupSelectKey = selectedGroupRef.value

// 导航菜单
const menus = computed(() => {
  const allMenus: Array<{
    name: string
    id: ProjectGroupSelectKey
    iconfont?: string
    elIcon?: any
  }> = [
    {
      name: t('app.project.projects.title'),
      id: 'all',
      iconfont: 'ac-xiangmu',
    },
    {
      name: t('app.project.stars.title'),
      id: 'followed',
      elIcon: markRaw(AcIconMyFollowedProject),
    },
    {
      name: t('app.project.mypro.title'),
      id: 'my',
      elIcon: markRaw(AcIconMyProject),
    },
  ]

  if (!teamStore.isMember) {
    allMenus.push({
      name: t('app.project.createModal.title'),
      id: 'create',
      elIcon: markRaw(AcIconAdd),
    })
  }

  return allMenus
})

// popover菜单
const popoverMenus = [
  {
    content: (
      <span style="display: flex">
        <Icon icon={'mdi:rename-outline'} width={18} />
        <span class="ml-1">{t('app.project.groups.rename')}</span>
      </span>
    ),
    onClick: () => currentModifyProjectGroup && onClickRenameProjectGroup(currentModifyProjectGroup),
  },
  {
    content: (
      <span style="display: flex; color: #ff5353">
        <Icon icon={'material-symbols:delete-outline'} width={18} />
        <span class="ml-1">{t('app.project.groups.delete')}</span>
      </span>
    ),
    onClick: () => currentModifyProjectGroup && onClickDeleteProjectGroup(currentModifyProjectGroup),
  },
]

// 是否需要重新选中分组

function needReSelectGroup(group: ProjectAPI.ResponseGroup) {
  if (selectedGroupRef.value === group?.id) {
    selectedGroupRef.value = 'all'
    selectedGroupKeyForUI.value = selectedGroupRef.value
    return
  }

  emits('delete')
}

// 重命名项目分组
function onClickRenameProjectGroup(group: ProjectAPI.ResponseGroup) {
  createOrUpdateProjectGroupRef?.value.show(group)
  hidePopover()
}

// 删除项目分组
function onClickDeleteProjectGroup(group: ProjectAPI.ResponseGroup) {
  AsyncMsgBox({
    title: t('app.project.groups.pop.title'),
    confirmButtonText: t('app.common.delete'),
    content: t('app.project.groups.pop.content'),
    onOk: async () => {
      await groupStore.deleteProjectGroup(group)
      hidePopover()
      needReSelectGroup(group)
    },
  })
}

// 点击项目分组（切换分组）
function handleItemClick(key: ProjectGroupSelectKey) {
  selectedGroupKeyForUI.value = key
}

// 显示创建项目分组弹窗
function showCreateProjectGroupModal() {
  createOrUpdateProjectGroupRef?.value.show()
}

// 处理拖拽
let currentGroupList: number[] = []
function onDragStart() {
  currentGroupList = projectGroups.value.map(group => group.id as number)
}
function handleDragEnd() {
  let different = false
  for (let i = 0; i < projectGroups.value.length; i++) {
    if (projectGroups.value[i].id !== currentGroupList[i]) {
      different = true
      break
    }
  }
  if (different)
    groupStore.sortGroup()
}

// 点击更多菜单
function onMoreMenuClick(e: MouseEvent, group: ProjectAPI.ResponseGroup) {
  currentModifyProjectGroup = group
  showPopover(e.target as HTMLElement)
}

// 选中项回退
function goBackSelected() {
  selectedGroupKeyForUI.value = selectedHistory
}

// 移除选中导航（所有导航不选中）
function removeSelected() {
  selectedGroupKeyForUI.value = null
}

// 导航选中（默认选中所有）
function goSelected(key: ProjectGroupSelectKey = 'all') {
  selectedHistory = key
  goBackSelected()
}

// 触发切换分组事件
function triggerSwitchProjectGroupEvent() {
  const key = selectedGroupRef.value
  const title: string | undefined = (typeof key === 'string' ? menus.value : projectGroups.value).find(
    (menu: any) => menu.id === key,
  )?.name
  emits('switchGroup', {
    key,
    title: title ?? t('app.project.projects.title'),
  })
}

watch(selectedGroupKeyForUI, (key: ProjectGroupSelectKey) => {
  if (key !== 'create') {
    selectedHistory = key
    selectedGroupRef.value = key
    triggerSwitchProjectGroupEvent()
    groupStore.saveGroupKeyToStorage(key)
    return
  }

  emits('createProject')
})

// 获取分组，手动触发切换分组事件
onBeforeMount(async () => {
  const groups = await groupStore.getProjectGroups()
  // 检测默认选中分组是否存在，不存在时切换分组
  if (typeof selectedGroupRef.value === 'number' && !groups.find(item => item.id === selectedGroupRef.value))
    selectedGroupRef.value = 'all'
  else triggerSwitchProjectGroupEvent()
})

defineExpose({
  goSelected,
  removeSelected,
  goBackSelected,
})
</script>

<template>
  <div class="flex flex-col h-full">
    <div :class="[ns.b(), ns.m('header')]">
      <div
        v-for="menu in menus"
        :key="menu.name"
        :class="[ns.e('item'), activeClass(menu.id)]"
        @click="handleItemClick(menu.id)"
      >
        <Iconfont v-if="menu.iconfont" :icon="menu.iconfont" :size="18" />
        <el-icon v-if="menu.elIcon" :size="18">
          <component :is="menu.elIcon" />
        </el-icon>
        <span>{{ menu.name }}</span>
      </div>
    </div>
    <div>
      <div :class="ns.e('segment')">
        <p>{{ $t('app.project.groups.title') }}</p>
        <el-icon class="cursor-pointer" @click="showCreateProjectGroupModal">
          <ac-icon-ep-plus />
        </el-icon>
      </div>
      <div class="flex-1 overflow-hidden mb-10px">
        <Draggable tag="ul" :class="ns.b()" :list="projectGroups" @start="onDragStart" @end="handleDragEnd">
          <li
            v-for="item in projectGroups"
            :key="item.id"
            :class="[
              ns.e('item'),
              ns.em('item', 'more'),
              activeClass(item.id as number),
            ]"
            :title="item.name"
            @click="handleItemClick(item.id as number)"
          >
            <div class="w-full flex-y-center">
              <span :class="ns.e('dot')" />
              <span :class="ns.e('title')">{{ item.name }}</span>
            </div>
            <el-icon @click.stop="onMoreMenuClick($event, item)">
              <ac-icon-ep-more-filled />
            </el-icon>
          </li>
        </Draggable>
      </div>
    </div>
  </div>

  <el-popover
    trigger="click"
    width="auto"
    virtual-triggering
    :virtual-ref="popoverRefEl"
    :visible="isShowPopoverMenu"
    :show-arrow="false"
  >
    <PopperMenu :menus="popoverMenus" size="small" class="clear-popover-space" />
  </el-popover>
</template>
