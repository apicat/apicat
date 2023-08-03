<template>
  <div class="ac-header-operate" v-if="hasDocument && definition">
    <div class="ac-header-operate__main">
      <p class="ac-header-operate__title">{{ definition.name }}</p>
    </div>

    <div class="ac-header-operate__btns" v-if="isManager || isWriter">
      <el-button type="primary" @click="() => goSchemaEditPage()">{{ $t('app.common.edit') }}</el-button>
      <el-tooltip effect="dark" content="历史记录" placement="bottom">
        <Iconfont class="cursor-pointer ac-history" :size="24" @click="goSchemaHistoryRecord" />
      </el-tooltip>
    </div>
  </div>

  <Result v-show="!hasDocument && !isLoading">
    <template #icon>
      <img class="h-auto w-260px mb-26px" src="@/assets/images/icon-empty.png" alt="" />
    </template>
  </Result>

  <div :class="[ns.b(), { 'h-20vh': !definition && hasDocument }]" v-loading="isLoading">
    <template v-if="definition">
      <h4>{{ definition.description }}</h4>
      <div class="text-right">
        <el-button @click="onShowCodeGenerate" :icon="AcIconoirCode">{{ $t('app.common.generateCode') }}</el-button>
      </div>
      <div class="ac-editor mt-10px"></div>
      <JSONSchemaEditor readonly v-model="definition.schema" :definitions="definitions" />
    </template>
  </div>

  <GenerateCodeModal ref="generateCodeModalRef" />
</template>

<script setup lang="ts">
import AcIconoirCode from '~icons/pepicons-pop/code'
import JSONSchemaEditor from '@/components/APIEditor/Editor.vue'
import { getDefinitionSchemaDetail } from '@/api/definitionSchema'
import { DefinitionSchema } from '@/components/APIEditor/types'
import { useNamespace } from '@/hooks'
import { useGoPage } from '@/hooks/useGoPage'
import { useParams } from '@/hooks/useParams'
import useDefinitionStore from '@/store/definition'
import useProjectStore from '@/store/project'
import { storeToRefs } from 'pinia'
import { getSchemaHistoryPath } from '@/router'

const GenerateCodeModal = defineAsyncComponent(() => import('@/components/GenerateCode/GenerateCodeModal.vue'))

const ns = useNamespace('document')
const route = useRoute()
const router = useRouter()
const definitionStore = useDefinitionStore()
const projectStore = useProjectStore()
const { project_id } = useParams()
const { goSchemaEditPage } = useGoPage()

const { definitions } = storeToRefs(definitionStore)
const { isManager, isWriter } = storeToRefs(projectStore)
const [isLoading, getDefinitionDetailApi] = getDefinitionSchemaDetail()

const definition = ref<DefinitionSchema | null>(null)
const hasDocument = ref(true)
const generateCodeModalRef = ref<InstanceType<typeof GenerateCodeModal>>()

const getDetail = async () => {
  const def_id = parseInt(route.params.schema_id as string, 10)

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
const goSchemaHistoryRecord = () => router.push(getSchemaHistoryPath(project_id, route.params.schema_id as string))

const onShowCodeGenerate = () => {
  generateCodeModalRef.value?.show(toRaw(definition.value)!)
}

definitionStore.$onAction(({ name, after, args }) => {
  // 删除全局模型
  if (name === 'deleteDefinition' && args[1] !== parseInt(route.params.schema_id as string, 10)) {
    after(() => getDetail())
  }
})

watch(
  () => route.params.schema_id,
  async () => await getDetail(),
  { immediate: true }
)
</script>
