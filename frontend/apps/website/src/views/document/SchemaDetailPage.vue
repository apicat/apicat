<template>
  <div class="ac-header-operate" v-if="hasDocument && definition">
    <div class="ac-header-operate__main">
      <p class="ac-header-operate__title">{{ definition.name }}</p>
    </div>

    <div class="ac-header-operate__btns" v-if="!isReader">
      <el-button type="primary" @click="() => goSchemaEditPage()">{{ $t('app.common.edit') }}</el-button>
    </div>
  </div>

  <Result v-show="!hasDocument && !isLoading">
    <template #icon>
      <img class="h-auto w-260px mb-26px" src="@/assets/images/icon-empty.png" alt="" />
    </template>
  </Result>

  <div :class="[ns.b(), { 'h-20vh': !definition && hasDocument }]" v-loading="isLoading">
    <div class="text-right"><el-button @click="onShowCodeGenerate" :icon="AcIconoirCode">Code Generate</el-button></div>
    <SchmaEditor v-if="definition" :readonly="true" v-model="definition" :definitions="definitions" />
    <GenerateCodeModal ref="generateCodeModalRef" />
  </div>
</template>
<script setup lang="ts">
import AcIconoirCode from '~icons/pepicons-pop/code'
import SchmaEditor from './components/SchemaEditor.vue'
import { getDefinitionSchemaDetail } from '@/api/definitionSchema'
import { DefinitionSchema } from '@/components/APIEditor/types'
import { useNamespace } from '@/hooks'
import { useGoPage } from '@/hooks/useGoPage'
import { useParams } from '@/hooks/useParams'
import useDefinitionStore from '@/store/definition'
import uesProjectStore from '@/store/project'
import { storeToRefs } from 'pinia'

const GenerateCodeModal = defineAsyncComponent(() => import('@/components/GenerateCode/GenerateCodeModal.vue'))

const ns = useNamespace('document')
const route = useRoute()
const definitionStore = useDefinitionStore()
const projectStore = uesProjectStore()
const { project_id } = useParams()
const { goSchemaEditPage } = useGoPage()

const { definitions } = storeToRefs(definitionStore)
const { isReader } = storeToRefs(projectStore)
const [isLoading, getDefinitionDetailApi] = getDefinitionSchemaDetail()

const definition = ref<DefinitionSchema | null>(null)
const hasDocument = ref(true)
const generateCodeModalRef = ref<InstanceType<typeof GenerateCodeModal>>()

const getDetail = async () => {
  const def_id = parseInt(route.params.shcema_id as string, 10)

  if (isNaN(def_id)) {
    hasDocument.value = false
    return
  }
  hasDocument.value = true

  try {
    const data = await getDefinitionDetailApi({ project_id, def_id })

    Object.defineProperty(data.schema, '_id', {
      value: data.id,
      enumerable: false,
      configurable: false,
      writable: false,
    })
    definition.value = data
  } catch (error) {
    //
  }
}

const onShowCodeGenerate = () => {
  generateCodeModalRef.value?.show()
}

definitionStore.$onAction(({ name, after, args }) => {
  // 删除全局模型
  if (name === 'deleteDefinition' && args[1] !== parseInt(route.params.shcema_id as string, 10)) {
    after(() => getDetail())
  }
})

watch(
  () => route.params.shcema_id,
  async () => await getDetail(),
  { immediate: true }
)
</script>
