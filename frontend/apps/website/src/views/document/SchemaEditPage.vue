<template>
  <div class="ac-header-operate">
    <div class="ac-header-operate__main">
      <p class="flex-y-center">
        <el-icon :size="18" class="mt-1px mr-4px"><ac-icon-ic-sharp-cloud-queue /></el-icon>
        {{ isSaving ? $t('app.common.saving') : $t('app.common.savedCloud') }}
      </p>
    </div>
    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="() => goSchemaDetailPage()">{{ $t('app.common.preview') }}</el-button>
    </div>
  </div>

  <div :class="[ns.b(), { 'h-50vh': !definition }]" v-loading="isLoading">
    <SchmaEditor v-if="definition" v-model="definition" :definitions="definitions" />
  </div>
</template>
<script setup lang="ts">
import { getDefinitionSchemaDetail } from '@/api/definitionSchema'
import SchmaEditor from './components/SchemaEditor.vue'
import useDefinitionStore from '@/store/definition'
import { storeToRefs } from 'pinia'
import { useParams } from '@/hooks/useParams'
import { DefinitionSchema } from '@/components/APIEditor/types'
import { debounce, isEmpty } from 'lodash-es'
import { useGoPage } from '@/hooks/useGoPage'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useNamespace } from '@/hooks'

const { t } = useI18n()
const route = useRoute()
const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)
const [isLoading, getDefinitionDetailApi] = getDefinitionSchemaDetail()
const { project_id, shcema_id } = useParams()
const { goSchemaDetailPage } = useGoPage()
const isUpdate = shcema_id !== undefined
const isSaving = ref(false)
const schemaTree: any = inject('schemaTree')
const ns = useNamespace('document')

const definition = ref<DefinitionSchema | null>(null)
const isInvalidId = () => isNaN(parseInt(route.params.shcema_id as string, 10))

const getDetail = async () => {
  // id 无效
  if (isInvalidId()) {
    return
  }

  try {
    const data = await getDefinitionDetailApi({ project_id, def_id: route.params.shcema_id })
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

definitionStore.$onAction(({ name, after, args }) => {
  // 删除全局模型
  if (name === 'deleteDefinition' && args[1] !== parseInt(route.params.shcema_id as string, 10)) {
    after(() => getDetail())
  }
})

watch(
  definition,
  debounce(async (newVal, oldVal) => {
    if (!oldVal || !oldVal.id) {
      return
    }

    if (isEmpty(newVal.name)) {
      ElMessage.error(t('app.schema.form.title'))
      return
    }
    if (definition.value) {
      isSaving.value = true
      try {
        const data: any = unref(definition)
        if (isUpdate && definition.value.id) {
          const { id: def_id, ...rest } = data
          await definitionStore.updateDefinition({ project_id, def_id, ...rest })
          schemaTree.updateTitle(def_id, newVal.name)
        }
      } catch (e) {
        //
      } finally {
        isSaving.value = false
      }
    }
  }, 300),
  {
    deep: true,
  }
)

watch(
  () => route.params.shcema_id,
  async () => await getDetail(),
  { immediate: true }
)
</script>
