<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import { useIterationList } from '../logic/useIterationList'
import { useNamespace } from '@/hooks/useNamespace'
import { Authority } from '@/commons'
import { useTeamStore } from '@/store/team'

const props = withDefaults(defineProps<IterationTableProps>(), {
  projectId: null,
  title: null,
})

const emits = defineEmits<{
  (e: 'remove' | 'edit', i: IterationAPI.ResponseIteration): void
}>()

export interface IterationTableProps {
  title: string | null
  projectId?: string | null
}

const { t } = useI18n()
const ns = useNamespace('iteration-list')
const { isLoading, data, currentPage, pageSize, total, handleRemoveIteration, handleRowClick, fetchIterationList }
  = useIterationList(props)
const titleRef = computed(() => props.title || t('app.iter.table.title'))

function onDeleteBtnClick(i: IterationAPI.ResponseIteration) {
  handleRemoveIteration(i.id)
  emits('remove', i)
}

const onEditBtnClick = (i: IterationAPI.ResponseIteration) => emits('edit', i)

defineExpose({
  reload: () => {
    currentPage.value = 1
    fetchIterationList()
  },
  refresh: fetchIterationList,
})
</script>

<template>
  <div v-loading="isLoading">
    <h3 class="truncate text-18px text-gray-title mb-30px">
      {{ titleRef }}
    </h3>

    <div v-if="data.length" :class="ns.b()">
      <div v-for="item in data" :key="item.id" :class="ns.e('item')" @click="handleRowClick(item)">
        <div class="flex-1 overflow-hidden">
          <p :class="ns.e('title')" class="m-0 overflow-hidden truncate" :title="item.title">
            {{ item.title }}
          </p>
          <div v-if="item.description" class="mt-1 truncate text-gray-helper" :title="item.description">
            {{ item.description }}
          </div>
          <div class="flex mt-1 text-gray-helper">
            <p class="w-35% truncate" :title="item.title">
              {{ $t('app.iter.table.project') + item.project?.title }}
            </p>
            <p class="w-20% mx-10px">
              {{ $t('app.iter.table.apis') + item.apisCount }}
            </p>
            <p>
              {{ $t('app.iter.table.created_at') + dayjs(item.createdAt * 1000).format('LLL LT') }}
            </p>
          </div>
        </div>
        <div :class="ns.e('operation')">
          <div v-if="item.project && [Authority.Manage, Authority.Write].includes(item.project.selfMember.permission)">
            <el-icon size="20" class="mr-10px" @click.stop="onEditBtnClick(item)">
              <ac-icon-ep-edit />
            </el-icon>
            <el-icon size="20" @click.stop="onDeleteBtnClick(item)">
              <ac-icon-ep-delete />
            </el-icon>
          </div>
        </div>
      </div>
    </div>

    <el-empty v-else :description="$t('app.iteration.list.emptyDataTip')" />

    <div v-if="total" class="flex justify-end mt-4">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        layout="prev, pager, next"
        :total="total"
      />
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(iteration-list) {
  border-top: 1px solid;
  @apply b-gray-110;

  @include e(item) {
    border-bottom: 1px solid;
    height: 100px;
    align-items: center;
    @apply flex cursor-pointer b-gray-110 py-10px px-24px rounded;

    @include e(title) {
      @apply text-gray-title mb-6px;
    }

    @include e(operation) {
      @apply flex-y-center;
      visibility: hidden;
      min-width: 100px;
      justify-content: flex-end;
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
