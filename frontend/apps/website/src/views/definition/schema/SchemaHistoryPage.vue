<template>
  <SchemaHistoryOperateHeader :title="definition?.name" />
  <div :class="[ns.b(), { 'h-20vh': !definition && hasDocument }]" v-loading="isLoading">
    <template v-if="definition">
      <h4>{{ definition.description }}</h4>
      <div class="ac-editor mt-10px"></div>
      <JSONSchemaEditor readonly v-model="definition.schema" :definitions="definitions" />
    </template>
  </div>
</template>
<script setup lang="ts">
import SchemaHistoryOperateHeader from './components/SchemaHistoryOperateHeader.vue'
import JSONSchemaEditor from '@/components/APIEditor/Editor.vue'
import { getSchemaHistoryRecordDetail } from '@/api/definitionSchema'
import { DefinitionSchema } from '@/components/APIEditor/types'
import { useNamespace } from '@/hooks'
import useApi from '@/hooks/useApi'
import { useParams } from '@/hooks/useParams'
import { useDefinitionSchemaStore } from '@/store/definitionSchema'
import { storeToRefs } from 'pinia'

const ns = useNamespace('document')
const route = useRoute()

const definitionSchemaStore = useDefinitionSchemaStore()
const { definitions } = storeToRefs(definitionSchemaStore)
const [isLoading, getDocumentHistoryRecordDetailApi] = useApi(getSchemaHistoryRecordDetail)
const { project_id, schema_id } = useParams()

const definition = ref<DefinitionSchema | null>(null)
const hasDocument = ref(true)

const getDetail = async (hid: string) => {
  const history_id = parseInt(hid, 10)

  if (isNaN(history_id)) {
    hasDocument.value = false
    return
  }
  hasDocument.value = true

  try {
    const data = await getDocumentHistoryRecordDetailApi({ project_id, def_id: schema_id, history_id })
    definition.value = data
  } catch (error) {
    //
  }
}

watch(
  () => route.params.history_id,
  async () => await getDetail(route.params.history_id as string),
  {
    immediate: true,
  }
)
</script>
