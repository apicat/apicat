<template>
  <div :class="ns.b()">
    <div :class="[ns.e('item'), activeClass()]" @click="() => handleClick()">
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
const emits = defineEmits(['add'])
const props = withDefaults(defineProps<Props>(), {
  projects: () => [],
  selectedProjectKey: null,
})

const ns = useNamespace('iteration-tree')
const selectedProjectKeyRef = useVModel(props, 'selectedProjectKey', emits)

const activeClass = (project?: ProjectInfo) => (project?.id === selectedProjectKeyRef.value ? 'active' : '')

const handleCreate = () => emits('add')

const handleClick = (project?: ProjectInfo) => {
  selectedProjectKeyRef.value = project ? (project.id as number) : null
}
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(iteration-tree) {
  @apply mt-24px;

  @include e(item) {
    @apply rounded-5px cursor-pointer h-40px text-14px flex flex-y-center px-20px;
    .ac-iconfont,
    .el-icon {
      @apply mr-8px;
    }

    span {
      @apply truncate;
    }

    &.active {
      background-color: rgba(204, 225, 255, 50);
      @apply text-#101010 font-500;
    }

    &:hover {
      @apply bg-gray-110;
    }

    & + & {
      margin-top: 10px;
    }
  }

  @include m(followed) {
    @include e(item) {
      list-style-type: disc;
    }
  }
}
</style>
