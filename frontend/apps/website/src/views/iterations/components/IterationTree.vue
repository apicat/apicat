<template>
  <div :class="ns.b()">
    <div :class="[ns.e('item'), activeClass('all')]" @click="handleItemClick('all')">
      <Iconfont icon="ac-diedai" :size="18" />
      <span>所有迭代</span>
    </div>
    <div :class="[ns.e('item'), activeClass('create')]" @click="handleItemClick('create')">
      <el-icon size="18"><ac-icon-ep-plus /></el-icon>
      <span>创建迭代</span>
    </div>
  </div>

  <p class="text-#101010 font-500 my-10px">关注的项目</p>
  <ul :class="ns.bm('followed')">
    <li v-for="project in projects" :class="[ns.e('item'), activeClass(project.id as number)]" @click="handleItemClick(project.id as number)">
      <span class="mr-8px">·</span><span>{{ project.title }}</span>
    </li>
  </ul>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { ProjectInfo } from '@/typings'

type SelectedKey = number | string | 'all' | 'create'

interface Props {
  projects: ProjectInfo[]
  selectedKey?: number | string | null
}

interface Events {
  (event: 'update:selectedKey', id: number | string | null): void
  (event: 'create'): void
  (event: 'click-item', id: number | string | null): void
}

const emits = defineEmits<Events>()

const props = withDefaults(defineProps<Props>(), {
  projects: () => [],
  selectedKey: null,
})

const ns = useNamespace('iteration-tree')
const selectedRef = ref<SelectedKey>(props.selectedKey ?? 'all')
const activeClass = (key: SelectedKey) => (selectedRef.value === key ? 'active' : '')

const handleItemClick = (key: SelectedKey) => {
  selectedRef.value = key

  if (key === 'create') {
    emits('create')
    return
  }
  emits('update:selectedKey', key === 'all' ? null : key)
  emits('click-item', key === 'all' ? null : key)
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
