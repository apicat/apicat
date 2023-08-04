<template>
  <div class="ac-header-operate" v-if="hasDocument">
    <div class="ac-header-operate__main">
      <p class="ac-header-operate__title">{{ title }}</p>
    </div>
    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="onSaveOrEditBtnClick">还原此历史记录</el-button>
      <el-button @click="() => schemaDiffModalRef?.show()">对比</el-button>
    </div>
  </div>
  <SchemaDiffModal ref="schemaDiffModalRef" />
</template>

<script setup lang="tsx">
import SchemaDiffModal from './SchemaDiffModal.vue'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { restoreSchemaByHistoryRecord } from '@/api/definitionSchema'
import { ElMessage as $Message } from 'element-plus'
import { storeToRefs } from 'pinia'
import { useParams } from '@/hooks/useParams'
import { useDefinitionSchemaStore } from '@/store/definition'

defineProps({
  title: {
    type: String,
    default: '',
  },
})

const goBack = inject('goBack') as () => void
const definitionSchemaStore = useDefinitionSchemaStore()
const { historyRecordTree } = storeToRefs(definitionSchemaStore)
const { currentRoute } = useRouter()
const { project_id, schema_id } = useParams()

const schemaDiffModalRef = ref<InstanceType<typeof SchemaDiffModal>>()

const hasDocument = computed(() => !!currentRoute.value.params.history_id && historyRecordTree.value.length !== 0)

const onSaveOrEditBtnClick = () => {
  const { history_id } = currentRoute.value.params
  if (!history_id) {
    return
  }

  AsyncMsgBox({
    title: '提示',
    content: <div class="break-all">确定还原此历史记录吗？</div>,
    onOk: () =>
      restoreSchemaByHistoryRecord({ project_id, def_id: schema_id, history_id }).then((res: any) => {
        $Message.success(res.msg || '还原成功')
        goBack()
      }),
  })
}
</script>
