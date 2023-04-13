<template>
  <div v-loading="isLoading">
    <el-table ref="tableRef" :data="tableData" :border="true" @selection-change="handleSelectionChange" class="w-full" :empty-text="$t('app.common.emptyDataTip')">
      <el-table-column type="selection" width="55" />
      <el-table-column property="title" label="名称" show-overflow-tooltip />
      <el-table-column property="deleted_at" label="删除时间" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleRestore([row])">{{ $t('app.common.restore') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-button v-show="tableData.length" class="mt-20px" type="primary" :disabled="!multipleSelectionRef.length" @click="handleRestore(multipleSelectionRef)">{{
      $t('app.common.restore')
    }}</el-button>
  </div>
</template>
<script setup lang="ts">
import { restoreDoc, getProjectTranshList } from '@/api/project'
import { useProjectId } from '@/hooks/useProjectId'
import { TrashModel } from '@/typings/project'
import type { ElTable } from 'element-plus'

const project_id = useProjectId()
const [isLoading, getProjectTranshListApi] = getProjectTranshList()
const [isRestoring, restoreDocApi] = restoreDoc()

const tableRef = ref<InstanceType<typeof ElTable>>()
const tableData = shallowRef<TrashModel[]>([])
const multipleSelectionRef = ref<TrashModel[]>([])

const handleSelectionChange = (val: TrashModel[]) => {
  multipleSelectionRef.value = val
}

const handleRestore = async (multipleSelection: TrashModel[]) => {
  try {
    isLoading.value = true
    await restoreDocApi({ project_id, ids: multipleSelection.map((item) => item.id) })
    tableData.value = await getProjectTranshListApi(project_id)
    tableRef.value!.clearSelection()
  } finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  tableData.value = await getProjectTranshListApi(project_id)
})
</script>