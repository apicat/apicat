<template>
  <el-dialog v-model="dialogVisible" append-to-body :close-on-click-modal="false" :close-on-press-escape="false" destroy-on-close title="Tips" width="30%">
    <el-table ref="multipleTableRef" :data="tableData" style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column label="Date" width="120">
        <template #default="scope">{{ scope.row.date }}</template>
      </el-table-column>
      <el-table-column property="name" label="Name" width="120" />
      <el-table-column property="address" label="Address" show-overflow-tooltip />
    </el-table>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'

const { dialogVisible, showModel } = useModal()
const multipleSelection = ref([])

const handleSelectionChange = (val: any) => {
  multipleSelection.value = val
}

const tableData: any = [
  {
    date: '2016-05-03',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-02',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-04',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-01',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-08',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-06',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-07',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
]

const show = (selected?: any) => {
  console.log(selected)
  showModel()
}

defineExpose({
  show,
})
</script>
