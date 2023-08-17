<template>
  <div v-loading="isLoading">
    <h3 class="text-18px text-gray-title mb-30px">{{ titleRef }}</h3>

    <div :class="ns.b()" v-if="data.length">
      <div :class="ns.e('item')" v-for="item in data" :key="item.id" @click="handleRowClick(item)">
        <div class="flex-1 overflow-hidden">
          <p :class="ns.e('title')">{{ item.title }}</p>
          <div class="flex text-gray-helper">
            <p class="w-20% truncate" :title="item.project_title">项目: {{ item.project_title }}</p>
            <p class="w-20% mx-10px">API数量: {{ item.api_num }}</p>
            <p>创建时间: {{ item.created_at }}</p>
          </div>
        </div>
        <div :class="ns.e('operation')">
          <el-icon size="20" class="mr-10px" @click.stop="onEditBtnClick(item)"><ac-icon-ep-edit /></el-icon>
          <el-icon size="20" @click.stop="onDeleteBtnClick(item)"><ac-icon-ep-delete /></el-icon>
        </div>
      </div>
    </div>

    <el-empty v-else />

    <div class="flex justify-end mt-4" v-if="total">
      <el-pagination :page-size="pageSize" layout="prev, pager, next" :total="total" v-model:current-page="currentPage" />
    </div>
  </div>
</template>
<script lang="ts" setup>
import { useNamespace } from '@/hooks/useNamespace'
import { Iteration } from '@/typings'
import { useIterationList } from '../logic/useIterationList'

const ns = useNamespace('iteration-list')

interface Props {
  title: string | null
  projectId?: string | number | null
}

const props = withDefaults(defineProps<Props>(), {
  projectId: null,
  title: null,
})

const emits = defineEmits<{
  (e: 'remove', i: Iteration): void
  (e: 'edit', i: Iteration): void
}>()

const { projectId: projectIdRef } = toRefs(props)
const { isLoading, data, currentPage, pageSize, total, handleRemoveIteration, handleRowClick, fetchIterationList } = useIterationList(projectIdRef)

const titleRef = computed(() => props.title || '所有迭代')

const onDeleteBtnClick = (i: Iteration) => {
  handleRemoveIteration(i)
  emits('remove', i)
}

const onEditBtnClick = (i: Iteration) => emits('edit', i)

defineExpose({
  reload: () => {
    currentPage.value = 1
    fetchIterationList()
  },
  refresh: fetchIterationList,
})
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(iteration-list) {
  border-top: 1px solid;
  @apply b-gray-110;

  @include e(item) {
    border-bottom: 1px solid;
    @apply flex cursor-pointer b-gray-110 py-10px px-24px rounded;

    @include e(title) {
      @apply text-gray-title mb-6px;
    }

    @include e(operation) {
      @apply flex-y-center;
      visibility: hidden;
    }

    &:hover {
      @apply bg-gray-110;

      @include e(operation) {
        visibility: visible;
      }
    }
  }

  @include e(operation) {
  }
}
</style>
