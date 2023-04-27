<template>
  <div class="ac-header-operate">
    <div class="ac-header-operate__main">
      <p class="flex-y-center">
        <el-icon :size="18" class="mt-1px mr-4px"><ac-icon-ic-sharp-cloud-queue /></el-icon>
        {{ isSaving ? '保存中...' : '已保存在云端' }}
      </p>
    </div>
    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="() => goSchemaDetailPage()">预览</el-button>
    </div>
  </div>
  <SchmaEditor v-loading="isLoading" v-model="definition" :definitions="definitions" />
</template>
<script setup lang="ts">
import { getDefinitionDetail } from '@/api/definition'
import SchmaEditor from './components/SchemaEditor.vue'
import useDefinitionStore from '@/store/definition'
import { storeToRefs } from 'pinia'
import { useParams } from '@/hooks/useParams'
import { Definition } from '@/components/APIEditor/types'
import createDefaultDefinition from './components/createDefaultDefinition'
import { debounce, isEmpty } from 'lodash-es'
import { useGoPage } from '@/hooks/useGoPage'
import { ElMessage } from 'element-plus'

const route = useRoute()
const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)
const [isLoading, getDefinitionDetailApi] = getDefinitionDetail()
const { project_id, shcema_id } = useParams()
const { goSchemaDetailPage } = useGoPage()
const isUpdate = shcema_id !== undefined
const isSaving = ref(false)
const schemaTree: any = inject('schemaTree')

const definition = ref<Definition>(createDefaultDefinition())
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
    console.error(error)
  }
}

watch(
  definition,
  debounce(async (newVal, oldVal) => {
    if (!oldVal.id) {
      return
    }

    if (isEmpty(newVal.name)) {
      ElMessage.error('请输入模型标题')
      return
    }

    isSaving.value = true
    try {
      const data = unref(definition)
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
