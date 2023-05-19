<template>
  <el-dialog
    v-model="dialogVisible"
    append-to-body
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    destroy-on-close
    :title="$t('app.definitionResponse.title')"
    width="40%"
  >
    <el-input v-model="search" :placeholder="$t('app.response.fullname')" />
    <el-table ref="multipleTableRef" :data="filterResponse" class="w-full" @selection-change="handleSelectionChange" @row-click="handelRowClick">
      <el-table-column type="selection" width="55" />
      <el-table-column property="name" :label="$t('app.response.table.name')" />
      <el-table-column property="code" :label="$t('app.response.table.code')" width="120" />
      <el-table-column property="description" :label="$t('app.response.table.desc')" show-overflow-tooltip />
    </el-table>
    <el-button type="primary" class="mt-20px" @click="handelConfrim">{{ $t('app.common.confirm') }}</el-button>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import useCommonResponseStore from '@/store/definitionResponse'
import { storeToRefs } from 'pinia'

const emits = defineEmits(['ok'])

const { dialogVisible, showModel, hideModel } = useModal()
const multipleSelection: any = ref([])
const multipleTableRef = ref()

const commonResponseStore = useCommonResponseStore()
const { responses } = storeToRefs(commonResponseStore)

const search = ref('')
const filterResponse = computed(() => responses.value.filter((data: any) => !search.value || data.name.toLowerCase().includes(search.value.toLowerCase())))

const show = async (selectedNameList: number[]) => {
  showModel()
  await nextTick()
  multipleTableRef.value.clearSelection()

  responses.value
    .filter((row) => selectedNameList.includes(row.id as any))
    .forEach((row) => {
      multipleTableRef.value.toggleRowSelection(row, true)
    })
}

const handleSelectionChange = (val: any) => {
  multipleSelection.value = val
}

const handelRowClick = (row: any) => {
  const hasItem = multipleSelection.value.find((item: any) => item.id === row.id)
  multipleTableRef.value!.toggleRowSelection(row, !hasItem)
}

const handelConfrim = () => {
  emits(
    'ok',
    multipleSelection.value.map((item: any) => item.id)
  )
  hideModel()
}

defineExpose({
  show,
})
</script>
