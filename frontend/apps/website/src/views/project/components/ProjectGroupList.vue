<template>
  <div :class="ns.b()">
    <div :class="[ns.e('item'), activeClass('all')]" @click="handleItemClick('all')">
      <Iconfont icon="ac-diedai" :size="18" />
      <span>所有项目</span>
    </div>
    <div :class="[ns.e('item'), activeClass('followed')]" @click="handleItemClick('followed')">
      <Iconfont icon="ac-diedai" :size="18" />
      <span>关注的项目</span>
    </div>
    <div :class="[ns.e('item'), activeClass('my')]" @click="handleItemClick('my')">
      <Iconfont icon="ac-diedai" :size="18" />
      <span>我的项目</span>
    </div>
    <div v-if="!isNormalUser" :class="[ns.e('item'), activeClass('create')]" @click="handleItemClick('create')">
      <el-icon size="18"><ac-icon-ep-plus /></el-icon>
      <span>{{ $t('app.project.createModal.title') }}</span>
    </div>
  </div>

  <div :class="ns.e('segment')">
    <p>项目分组</p>
    <el-icon @click="onCreateProjectGroup" class="cursor-pointer"><ac-icon-ep-plus /></el-icon>
  </div>
  <ul>
    <li v-for="item in groups" :class="[ns.e('item'), activeClass(item.id!)]" @click="handleItemClick(item.id!)">
      <span :class="ns.e('dot')"></span>
      <span>{{ item.name }}</span>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { useUserStore } from '@/store/user'
import { ProjectGroup, ProjectGroupSelectKey } from '@/typings'

const emits = defineEmits<{
  (event: 'create-project'): void
  (event: 'create-group'): void
  (event: 'click', key: ProjectGroupSelectKey): void
}>()

const props = withDefaults(
  defineProps<{
    groups?: ProjectGroup[]
  }>(),
  {
    groups: () => [],
  }
)

const { isNormalUser } = useUserStore()
const ns = useNamespace('group-list')
const selectedRef = ref<ProjectGroupSelectKey>('all')
let selectedHistory: ProjectGroupSelectKey = 'all'

const handleItemClick = (key: ProjectGroupSelectKey) => {
  selectedRef.value = key

  if (key !== 'create') {
    selectedHistory = key
    emits('click', key)
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
