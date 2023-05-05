<template>
  <el-dialog v-model="dialogVisible" center append-to-body :close-on-click-modal="false" :close-on-press-escape="false" destroy-on-close width="40%">
    <template #header>
      <div class="flex-y-center">
        <el-icon class="mr-5px"><ac-icon-bi-robot /></el-icon>{{ $t('app.interface.common.aiGenerateInterface') }}
      </div>
    </template>

    <div v-loading="isLoading">
      <el-table :data="collectList" class="w-full" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column property="method" :label="$t('app.interface.table.method')" width="100">
          <template #default="{ row }">
            <div class="flex-y-center">
              <el-icon v-if="row.isLoading || row.isFinish" :class="{ 'animate-spin': row.isLoading }">
                <ac-icon-ep-loading v-if="!row.isFinish && row.isLoading" />
                <ac-icon-ep:circle-check-filled class="text-[var(--el-color-success)]" v-if="row.isFinish && row.isSuccess" />
                <ac-icon-ep:circle-close-filled class="text-[var(--el-color-danger)]" v-if="row.isFinish && !row.isSuccess" />
              </el-icon>
              {{ row.method }}
            </div>
          </template>
        </el-table-column>
        <el-table-column property="path" :label="$t('app.interface.table.path')" show-overflow-tooltip />
        <el-table-column property="description" :label="$t('app.interface.table.desc')" show-overflow-tooltip />
      </el-table>
      <el-button class="mt-20px" :disabled="!collectList.length" :loading="isStartCreate" type="primary" @click="handleCreate(multipleSelection)">{{
        $t('app.common.create')
      }}</el-button>
    </div>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import { createCollection, createCollectionByAI, createCollectionWithSchemaByAI } from '@/api/collection'
import useApi from '@/hooks/useApi'
import { useParams } from '@/hooks/useParams'
import { ElMessage } from 'element-plus'
import { uuid } from '@apicat/shared'
import { DocumentTypeEnum } from '@/commons'
import { useI18n } from 'vue-i18n'

const emits = defineEmits(['ok'])
const { t } = useI18n()

const { dialogVisible, showModel, hideModel } = useModal()
const [isLoading, createCollectionWithSchemaByAIApi] = useApi(createCollectionWithSchemaByAI)()
const { project_id } = useParams()

const multipleSelection = ref([])
const collectList = shallowRef([])
const isStartCreate = ref(false)

const show = async (schema: any) => {
  showModel()
  try {
    collectList.value = []
    const data = await createCollectionWithSchemaByAIApi({ project_id, schema_id: schema.id })
    collectList.value = (data || []).map((item: any, idx: number) => {
      return {
        ...item,
        _id: uuid(),
        sort: idx,
        schema,
        isLoading: false,
        isSuccess: false,
        isFinish: false,
      }
    })
  } catch (error) {
    //
  }
}

const handleSelectionChange = (val: any) => {
  multipleSelection.value = val
}

// 创建文档分类
const createCategory = async (schema: any) => {
  try {
    const data: any = await createCollection({ project_id, title: schema.name, type: DocumentTypeEnum.DIR })
    return data.id
  } catch (error) {
    return 0
  }
}

const handleCreate = async (selectedRows: Array<any>) => {
  if (!selectedRows.length) {
    ElMessage.error(t('app.interface.tips.unselectedInterface'))
    return
  }

  isStartCreate.value = true

  const parent_id = await createCategory(selectedRows[0].schema)
  selectedRows.sort((a, b) => a.sort - b.sort)
  const len = selectedRows.length

  for (let i = 0; i < len; i++) {
    // 隐藏modal
    if (!dialogVisible.value) {
      break
    }

    const item = selectedRows[i]
    item.isLoading = true
    item.abortController = new AbortController()

    try {
      const data: any = await createCollectionByAI({ project_id, parent_id, schema_id: item.schema.id, title: item.description }, { signal: item.abortController.signal })
      item.isSuccess = true
      item.isFinish = true
      item.isLoading = false

      emits('ok', data.id)
    } catch (error) {
      item.isSuccess = false
      item.isFinish = true
      item.isLoading = false
    }
  }

  isStartCreate.value = false

  if (selectedRows.every((item) => !item.isSuccess)) {
    ElMessage.error(t('app.interface.tips.allInterfaceCreateFailure'))
  }

  hideModel()
}

watch(dialogVisible, () => {
  if (!dialogVisible.value) {
    multipleSelection.value.forEach((item: any) => item.abortController?.abort())
  }
})

defineExpose({
  show,
})
</script>
