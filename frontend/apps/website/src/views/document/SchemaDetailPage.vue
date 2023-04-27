<template>
  <div class="ac-header-operate" v-if="hasDocument">
    <div class="ac-header-operate__main">
      <p class="ac-header-operate__title" v-if="definition">{{ definition.name }}</p>
    </div>

    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="() => goSchemaEditPage()">编辑</el-button>
    </div>
  </div>

  <SchmaEditor v-if="hasDocument" v-loading="isLoading" :readonly="true" v-model="definition" :definitions="definitions" />

  <Result v-if="!hasDocument">
    <template #icon>
      <img class="h-auto w-260px mb-26px" src="@/assets/images/icon-empty.png" alt="" />
    </template>
    <template #title>
      <div class="m-auto">您当前尚未创建模型。</div>
    </template>
  </Result>
</template>
<script setup lang="ts">
import SchmaEditor from './components/SchemaEditor.vue'
import { getDefinitionDetail } from '@/api/definition'
import { Definition } from '@/components/APIEditor/types'
import { useGoPage } from '@/hooks/useGoPage'
import { useParams } from '@/hooks/useParams'
import useDefinitionStore from '@/store/definition'
import { storeToRefs } from 'pinia'
import createDefaultDefinition from './components/createDefaultDefinition'

const route = useRoute()
const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)
const [isLoading, getDefinitionDetailApi] = getDefinitionDetail()
const { project_id } = useParams()
const { goSchemaEditPage } = useGoPage()
const definition = ref<Definition>(createDefaultDefinition())
const hasDocument = ref(true)

const getDetail = async (defId: string) => {
  const def_id = parseInt(defId, 10)

  if (isNaN(def_id)) {
    hasDocument.value = false
    return
  }
  hasDocument.value = true
  const data = await getDefinitionDetailApi({ project_id, def_id })

  Object.defineProperty(data.schema, '_id', {
    value: data.id,
    enumerable: false,
    configurable: false,
    writable: false,
  })
  definition.value = data
}

watch(
  () => route.params.shcema_id,
  async () => await getDetail(route.params.shcema_id as string),
  {
    immediate: true,
  }
)
</script>
