<template>
  <div :class="ns.b()">
    <div :class="[ns.e('item'), activeClass(null)]" @click="() => handleClick()">
      <Iconfont icon="ac-diedai" :size="18" />
      <span>所有迭代</span>
    </div>
    <div :class="ns.e('item')" @click="handleCreate">
      <el-icon size="18"><ac-icon-ep-plus /></el-icon>
      <span>创建迭代</span>
    </div>
  </div>

  <p class="text-#101010 font-500 my-10px">关注的项目</p>
  <ul :class="ns.bm('followed')">
    <li v-for="project in projects" :class="[ns.e('item'), activeClass(project.id as number)]" @click="() => handleClick(project)">
      <span class="mr-8px">·</span><span>{{ project.title }}</span>
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
const emits = defineEmits(['add'])
const props = withDefaults(defineProps<Props>(), {
  projects: () => [],
  selectedProjectKey: null,
})

const ns = useNamespace('iteration-tree')
const selectedProjectKeyRef = useVModel(props, 'selectedProjectKey', emits)

const activeClass = (id: number | null = null) => (selectedProjectKeyRef.value === id ? 'active' : '')

const handleCreate = () => emits('add')

const handleClick = (project?: ProjectInfo) => {
  selectedProjectKeyRef.value = project ? (project.id as number) : null
}
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(iteration-tree) {
  @apply mt-24px text-#101010;

  @include e(item) {
    @apply rounded-5px cursor-pointer h-40px text-14px flex flex-y-center px-20px;

    & + & {
      margin-top: 4px;
    }

    .ac-iconfont,
    .el-icon {
      @apply mr-8px;
    }

    span {
      @apply truncate;
    }

    &.active {
      background-color: #e5f0ff;
      @apply text-#101010 font-500;
    }

    &:hover {
      @apply bg-gray-110;
    }
  }

  @include m(followed) {
    @include e(item) {
      list-style-type: disc;
    }
  }
}
</style>
