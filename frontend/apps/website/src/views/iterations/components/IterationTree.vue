<template>
  <div :class="ns.b()">
    <div :class="[ns.e('item'), activeClass('all')]" @click="handleItemClick('all')">
      <Iconfont icon="ac-diedai" :size="18" />
      <span>所有迭代</span>
    </div>
    <div :class="[ns.e('item'), activeClass('create')]" @click="handleItemClick('create')">
      <el-icon :size="18"><ac-icon-fe:plus /></el-icon>
      <span>创建迭代</span>
    </div>
  </div>

  <p v-if="followedProjects.length" :class="ns.e('segment')">关注的项目</p>
  <ul>
    <li v-for="project in followedProjects" :class="[ns.e('item'), activeClass(project)]" @click="handleItemClick(project)">
      <span :class="ns.e('dot')"></span>
      <span>{{ project.title }}</span>
    </li>
  </ul>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { useFollowedProjectList } from '../logic/useFollowedProjectList'
import { ProjectInfo, SelectedKey } from '@/typings'
interface Events {
  (event: 'create'): void
  (event: 'click', project: ProjectInfo | null): void
}

const emits = defineEmits<Events>()
const ns = useNamespace('group-list')

const { followedProjects, activeClass, selectedRef, setSelectedHistory, goBackSelected, goSelectedAll, removeSelected } = useFollowedProjectList()

const handleItemClick = (project: SelectedKey) => {
  selectedRef.value = project

  if (project === 'create') {
    emits('create')
  } else {
    setSelectedHistory(project === 'all' ? 0 : (project as ProjectInfo))
    emits('click', project === 'all' ? null : (project as ProjectInfo))
  }
}

defineExpose({
  goSelectedAll,
  goBackSelected,
  removeSelected,
})
</script>
