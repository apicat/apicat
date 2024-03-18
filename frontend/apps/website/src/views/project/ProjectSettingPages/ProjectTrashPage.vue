<script setup lang="ts">
import { ElTable } from 'element-plus'
import { storeToRefs } from 'pinia'
import dayjs from 'dayjs'

import useProjectStore from '@/store/project'
import useApi from '@/hooks/useApi'
import { apiGetTrashList, apiRestoreTrash } from '@/api/project/setting/trash'
import { useCollectionsStore } from '@/store/collections'

const collectionStore = useCollectionsStore()
const [isLoading, getTrash] = useApi(apiGetTrashList)
const [, restoreTrash] = useApi(apiRestoreTrash)
const projectStore = useProjectStore()
const { isReader } = storeToRefs(projectStore)

const tableRef = ref<InstanceType<typeof ElTable>>()
const tableData = shallowRef<ProjectAPI.Trash[]>([])
const multipleSelectionRef = ref<ProjectAPI.Trash[]>([])

function handleSelectionChange(val: ProjectAPI.Trash[]) {
  multipleSelectionRef.value = val
}

async function handleRestore(multipleSelection: ProjectAPI.Trash[]) {
  try {
    isLoading.value = true
    await restoreTrash(projectStore.project!.id, {
      collectionIDs: multipleSelection.map(item => item.collectionID),
    })
    tableData.value = (await getTrash(projectStore.project!.id))!
    tableRef.value!.clearSelection()
    collectionStore.getCollections(projectStore.projectID!)
  }
  finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  tableData.value = (await getTrash(projectStore.project!.id))!
})
</script>

<template>
  <div v-loading="isLoading">
    <ElTable
      ref="tableRef"
      :data="tableData"
      :border="true"
      class="w-full"
      :empty-text="$t('app.common.emptyDataTip')"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column
        property="collectionTitle"
        :label="$t('app.project.setting.trash.table.title')"
        show-overflow-tooltip
      />
      <el-table-column property="deletedBy" :label="$t('app.project.setting.trash.table.deletedName')" />
      <el-table-column property="deletedAt" :label="$t('app.project.setting.trash.table.deletedAt')">
        <template #default="{ row }">
          {{ dayjs(new Date(row.deletedAt).getTime()).format('LLL LT') }}
        </template>
      </el-table-column>
      <el-table-column v-if="!isReader">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleRestore([row])">
            {{ $t('app.project.setting.trash.table.restore') }}
          </el-button>
        </template>
      </el-table-column>
    </ElTable>

    <el-button
      v-if="!isReader"
      v-show="tableData.length"
      class="mt-20px"
      type="primary"
      :disabled="!multipleSelectionRef.length"
      @click="handleRestore(multipleSelectionRef)"
    >
      {{ $t('app.common.restore') }}
    </el-button>
  </div>
</template>
