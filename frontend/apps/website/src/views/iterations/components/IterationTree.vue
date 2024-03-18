<script setup lang="ts">
import { useNamespace } from '@apicat/hooks'
import { useFollowedProjectList } from '../logic/useFollowedProjectList'

export interface IterationTreeEmits {
  (event: 'create'): void
  (event: 'click', project: ProjectAPI.ResponseProject | null): void
}

const emits = defineEmits<IterationTreeEmits>()
const ns = useNamespace('group-list')
const {
  followedProjects,
  activeClass,
  selectedRef,
  setSelectedHistory,
  goBackSelected,
  goSelectedAll,
  removeSelected,
} = useFollowedProjectList(emits)

function handleItemClick(projectID: IterationSelectedKey) {
  selectedRef.value = projectID
  setSelectedHistory(projectID)
  if (projectID === 'create') {
    emits('create')
    return
  }
  emits('click', projectID === 'all' ? null : followedProjects.value.find(val => val.id === projectID)!)
}

defineExpose({
  goSelectedAll,
  goBackSelected,
  removeSelected,
})
</script>

<template>
  <div>
    <div :class="[ns.b(), ns.m('header')]">
      <div :class="[ns.e('item'), activeClass('all')]" @click="handleItemClick('all')">
        <Iconfont icon="ac-diedai" :size="18" />
        <span> {{ $t('app.iter.table.title') }}</span>
      </div>
      <div :class="[ns.e('item'), activeClass('create')]" @click="handleItemClick('create')">
        <el-icon :size="18">
          <ac-icon-fe:plus />
        </el-icon>
        <span> {{ $t('app.iter.create.title') }}</span>
      </div>
    </div>

    <p v-if="followedProjects.length" :class="ns.e('segment')">
      {{ $t('app.iter.star.title') }}
    </p>
    <ul :class="[ns.b()]">
      <li
        v-for="project in followedProjects"
        :key="project.id"
        :class="[ns.e('item'), activeClass(project.id)]"
        :title="project.title"
        @click="handleItemClick(project.id)"
      >
        <span :class="ns.e('dot')" />
        <span :class="ns.e('title')">{{ project.title }}</span>
      </li>
    </ul>
  </div>
</template>
