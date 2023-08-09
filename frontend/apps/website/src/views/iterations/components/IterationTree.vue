<template>
  <div :class="ns.b()">
    <div :class="[ns.e('item'), activeClass()]" @click="() => handleClick()">
      <Iconfont icon="" />
      <span>所有迭代</span>
    </div>
    <div :class="ns.e('item')" @click="handleCreate">
      <Iconfont icon="" />
      <span>创建迭代</span>
    </div>
  </div>

  <p>关注的项目</p>
  <ul :class="ns.bm('followed')">
    <li v-for="project in projects" :class="[ns.e('item'), activeClass(project)]" @click="() => handleClick(project)">
      <span>ApiCat</span>
    </li>
  </ul>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { ProjectInfo } from '@/typings'

interface Props {
  projects: ProjectInfo[]
  selectedProjectKey: number | null
}

const props = withDefaults(defineProps<Props>(), {
  projects: () => [],
  selectedProjectKey: null,
})

const emits = defineEmits(['update:selectedProjectKey', 'add'])

const { selectedProjectKey: selectedProjectKeyRef } = toRefs(props)

const ns = useNamespace('iteration-tree')

const activeClass = (project?: ProjectInfo) => (project?.id === selectedProjectKeyRef.value ? 'active' : '')

const handleCreate = () => emits('add')

const handleClick = (project?: ProjectInfo) => {
  emits('update:selectedProjectKey', project ? project.id : null)
}
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(iteration-tree) {
  @include e(item) {
  }

  @include m(followed) {
    @include e(item) {
      list-style-type: disc;
    }
  }
}
</style>
