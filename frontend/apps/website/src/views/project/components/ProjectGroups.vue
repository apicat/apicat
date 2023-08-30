<template>
  <div :class="ns.b()">
    <div v-for="menu in menus" :key="menu.title" :class="[ns.e('item'), activeClass(menu.key)]" @click="handleItemClick(menu.key, menu.title)">
      <Iconfont v-if="menu.iconfont" :icon="menu.iconfont" :size="18" />
      <el-icon v-if="menu.elIcon" :size="18"><component :is="menu.elIcon" /></el-icon>
      <span>{{ menu.title }}</span>
    </div>
  </div>

  <div :class="ns.e('segment')">
    <p>项目分组</p>
    <el-icon v-if="!isNormalUser" @click="onCreateProjectGroup" class="cursor-pointer"><ac-icon-ep-plus /></el-icon>
  </div>

  <ul>
    <li v-for="item in groups" :class="[ns.e('item'), activeClass(item.id!)]" @click="handleItemClick(item.id!, item.name)">
      <span :class="ns.e('dot')"></span>
      <span>{{ item.name }}</span>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { useUserStore } from '@/store/user'
import { ProjectGroup, ProjectGroupSelectKey, SwitchProjectGroupInfo } from '@/typings'
import AcIconAdd from '~icons/fe/plus'
import AcIconMyProject from '~icons/mdi/folder-outline'
import AcIconMyFollowedProject from '~icons/clarity/star-line'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'

const emits = defineEmits<{
  (event: 'create-project'): void
  (event: 'create-group'): void
  (event: 'change-group', info: SwitchProjectGroupInfo): void
  (event: 'update:selected', key: ProjectGroupSelectKey): void
}>()

const props = withDefaults(
  defineProps<{
    selected: ProjectGroupSelectKey
    groups?: ProjectGroup[]
  }>(),
  {
    selected: 'all',
    groups: () => [],
  }
)

const userStore = useUserStore()
const { isNormalUser } = storeToRefs(userStore)
const ns = useNamespace('group-list')
const { t } = useI18n()

const selectedRef = ref<ProjectGroupSelectKey>(props.selected)
let selectedHistory: ProjectGroupSelectKey = props.selected

// 导航菜单
const menus = computed(() => {
  let allMenus: Array<{ title: string; key: ProjectGroupSelectKey; iconfont?: string; elIcon?: any }> = [
    { title: '所有项目', key: 'all', iconfont: 'ac-xiangmu' },
    { title: '关注的项目', key: 'followed', elIcon: markRaw(AcIconMyFollowedProject) },
    { title: '我的项目', key: 'my', elIcon: markRaw(AcIconMyProject) },
    { title: t('app.project.createModal.title'), key: 'create', elIcon: markRaw(AcIconAdd) },
  ]

  // 普通成员移除创建入口
  if (isNormalUser.value) {
    allMenus = allMenus.filter((menu) => menu.key !== 'create')
  }

  return allMenus
})

const handleItemClick = (key: ProjectGroupSelectKey, title: string) => {
  selectedRef.value = key

  if (key !== 'create') {
    selectedHistory = key

    emits('change-group', { key, title })
    emits('update:selected', key)
    return
  }

  emits('create-project')
}

const onCreateProjectGroup = () => emits('create-group')

const activeClass = (key: ProjectGroupSelectKey) => (selectedRef.value === key ? 'active' : '')

const goBackSelected = () => {
  selectedRef.value = selectedHistory
}

const removeSelected = () => {
  selectedRef.value = null
}

const goSelected = (key: ProjectGroupSelectKey = 'all') => {
  selectedHistory = key
  goBackSelected()
}

defineExpose({
  goSelected,
  removeSelected,
  goBackSelected,
})
</script>
