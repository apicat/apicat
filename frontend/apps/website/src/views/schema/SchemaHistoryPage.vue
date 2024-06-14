<script setup lang="ts">
import { useNamespace } from '@apicat/hooks'
import { JSONSchemaTable } from '@apicat/components'
import { useI18n } from 'vue-i18n'
import SchemaDiffDialog from './SchemaDiffDialog.vue'
import { apiGetSchemaHistoryInfo, apiRestoreSchemaHistory } from '@/api/project/definition/schema'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import useApi from '@/hooks/useApi'
import { useHistoryLayoutContext } from '@/layouts/useHistoryLayoutContext'
import useDefinitionSchemaStore from '@/store/definitionSchema'

const props = defineProps<{
  projectID: string
  schemaID: string
  historyID?: string
}>()

const { t } = useI18n()
const ns = useNamespace('document')
const { goBack } = useHistoryLayoutContext()
const schemaStore = useDefinitionSchemaStore()
const schemaHistoryInfo = ref<HistoryRecord.SchemaHistoryInfo>()
const [loading, getInfo] = useApi(apiGetSchemaHistoryInfo)
const schemaDiffRef = ref<InstanceType<typeof SchemaDiffDialog>>()

function restoreHisotry() {
  AsyncMsgBox({
    title: t('app.schema.history.restore.poptitle'),
    content: t('app.schema.history.restore.popcontent'),
    onOk: async () => {
      await apiRestoreSchemaHistory(
        props.projectID,
        Number.parseInt(props.schemaID),
        Number.parseInt(props.historyID!),
      )
      ElMessage.success(t('app.schema.history.restore.success'))
      goBack()
    },
  })
}

watch(
  () => props.historyID,
  async () => {
    if (props.historyID) {
      schemaHistoryInfo.value = await getInfo(
        props.projectID,
        Number.parseInt(props.schemaID),
        Number.parseInt(props.historyID),
      )
    }
  },
  { immediate: true },
)
</script>

<template>
  <ElEmpty v-if="!schemaHistoryInfo" v-loading="loading" />
  <div v-else>
    <div class="ac-header-operate">
      <div class="ac-header-operate__main">
        <p class="ac-header-operate__title">
          {{ schemaHistoryInfo.name }}
        </p>
      </div>
      <div class="ac-header-operate__btns">
        <el-button type="primary" @click="restoreHisotry">
          {{ $t('app.schema.history.restore.title') }}
        </el-button>
        <el-button @click="schemaDiffRef?.show()">
          {{ $t('app.schema.history.diff.title') }}
        </el-button>
      </div>
    </div>
    <div v-loading="loading" :class="[ns.b()]">
      <JSONSchemaTable :schema="schemaHistoryInfo.schema" readonly :definition-schemas="schemaStore.schemas" />
    </div>
    <SchemaDiffDialog ref="schemaDiffRef" />
  </div>
</template>
