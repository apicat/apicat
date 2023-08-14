<template>
  <div>
    <h3 class="text-18px text-gray-title mb-30px">{{ titleRef }}</h3>

    <div :class="ns.b()" v-if="data.length">
      <div :class="ns.e('item')" v-for="item in data" :key="item.id">
        <div class="flex-1">
          <p :class="ns.e('title')">{{ item.title }}</p>
          <div class="text-gray-helper">
            <span>项目:{{ item.project_title }}</span>
            <span class="mx-10px">API数量:{{ item.api_num }}</span>
            <span>创建时间:{{ item.created_at }}</span>
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
      <el-pagination :page-size="pageSize" layout="prev, pager, next" :total="total" v-model:current-page="pageVModel" />
    </div>
  </div>
</template>
<script lang="ts" setup>
import { useNamespace } from '@/hooks/useNamespace'
import { Iteration } from '@/typings'

const ns = useNamespace('iteration-list')

interface Props {
  title?: string
  data: Iteration[]
  page: number
  total: number
  pageSize?: number
}

const props = withDefaults(defineProps<Props>(), {
  data: () => [],
  page: 1,
  total: 0,
  pageSize: 15,
})

const emits = defineEmits<{
  (e: 'remove', i: Iteration): void
  (e: 'edit', i: Iteration): void
  (e: 'update:page', page: number): void
}>()

const pageVModel = useVModel(props, 'page', emits)
const titleRef = computed(() => props.title || '所有迭代')

const onDeleteBtnClick = (i: Iteration) => emits('remove', i)
const onEditBtnClick = (i: Iteration) => emits('edit', i)
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
      display: none;
    }

    &:hover {
      @apply bg-gray-110;

      @include e(operation) {
        display: flex;
      }
    }
  }

  @include e(operation) {
  }
}
</style>
