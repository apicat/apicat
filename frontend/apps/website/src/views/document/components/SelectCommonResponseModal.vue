<template>
  <el-dialog v-model="dialogVisible" append-to-body :close-on-click-modal="false" :close-on-press-escape="false" destroy-on-close title="公共响应" width="40%">
    <el-input v-model="search" placeholder="响应名称" />
    <el-table ref="multipleTableRef" :data="filterResponse" class="w-full" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column property="name" label="名称" />
      <el-table-column property="code" label="状态码" width="120" />
      <el-table-column property="description" label="描述" show-overflow-tooltip />
    </el-table>
    <el-button type="primary" class="mt-20px" @click="handelConfrim">确定</el-button>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import useCommonResponseStore from '@/store/commonResponse'
import { storeToRefs } from 'pinia'

const emits = defineEmits(['ok'])

const { dialogVisible, showModel, hideModel } = useModal()
const multipleSelection = ref([])
const multipleTableRef = ref()

const commonResponseStore = useCommonResponseStore()
const { response } = storeToRefs(commonResponseStore)

const search = ref('')
const filterResponse = computed(() => response.value.filter((data: any) => !search.value || data.name.toLowerCase().includes(search.value.toLowerCase())))

const show = async (selectedNameList: string[]) => {
  showModel()
  await nextTick()
  multipleTableRef.value.clearSelection()

  response.value
    .filter((row) => selectedNameList.includes(row.name))
    .forEach((row) => {
      multipleTableRef.value.toggleRowSelection(row, true)
    })
}

const handleSelectionChange = (val: any) => {
  multipleSelection.value = val
}

const handelConfrim = () => {
  console.log(multipleSelection.value)
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
