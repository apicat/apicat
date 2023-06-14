<template>
  <div v-loading="isLoading">
    <el-table ref="tableRef" :data="tableData" :border="true" @selection-change="handleSelectionChange" class="w-full" :empty-text="$t('app.common.emptyDataTip')">
      <el-table-column type="selection" width="55" />
      <el-table-column property="title" :label="$t('app.table.name')" show-overflow-tooltip />
      <el-table-column property="deleted_at" :label="$t('app.table.deleteAt')" />
      <el-table-column :label="$t('app.table.operation')" v-if="!isReader">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleRestore([row])">{{ $t('app.common.restore') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-button v-if="!isReader" v-show="tableData.length" class="mt-20px" type="primary" :disabled="!multipleSelectionRef.length" @click="handleRestore(multipleSelectionRef)">
      {{ $t('app.common.restore') }}
    </el-button>
  </div>
</template>
<script setup lang="ts">
import { restoreDoc, getProjectTranshList } from '@/api/project'
import { useParams } from '@/hooks/useParams'
import uesProjectStore from '@/store/project'
import { TrashModel } from '@/typings/project'
import type { ElTable } from 'element-plus'
import { storeToRefs } from 'pinia'

const { project_id } = useParams()
const [isLoading, getProjectTranshListApi] = getProjectTranshList()
const [isRestoring, restoreDocApi] = restoreDoc()
const directoryTree: any = inject('directoryTree')
const projectStore = uesProjectStore()
const { isReader } = storeToRefs(projectStore)

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
    directoryTree?.reload()
  } finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  tableData.value = await getProjectTranshListApi(project_id)
})
</script>
