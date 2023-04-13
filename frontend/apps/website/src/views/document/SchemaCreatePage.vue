<template>
  <div class="ac-header-operate">
    <div class="ac-header-operate__main">
      <p>
        <i class="ac-iconfont"></i>
        {{ isSaving ? '保存中...' : '已保存在云端' }}
      </p>
    </div>

    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="onPreviewBtnClick">预览</el-button>
    </div>
  </div>

  <SchmaEditor v-model="definition" :definitions="definitions" />
</template>
<script setup lang="ts">
import SchmaEditor from './components/SchemaEditor.vue'
import useDefinitionStore from '@/store/definition'
import { storeToRefs } from 'pinia'
import { useParams } from '@/hooks/useParams'
import { Definition } from '@/components/APIEditor/types'
import createDefaultDefinition from './components/createDefaultDefinition'
import { debounce, isEmpty } from 'lodash-es'
import { ElMessage } from 'element-plus'
import { useGoPage } from '@/hooks/useGoPage'

const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)
const { project_id } = useParams()
const isSaving = ref(false)
const definition = ref<Definition>(createDefaultDefinition())
const { goSchemaDetailPage } = useGoPage()
const schemaTree: any = inject('schemaTree')
let def_id: any = null

watch(
  definition,
  debounce(async () => {
    const data = unref(definition)
    if (isEmpty(data.name)) {
      ElMessage.error('请输入模型标题')
      return
    }

    isSaving.value = true
    try {
      // update
      if (definition.value.id) {
        const { id: def_id, ...rest } = data
        await definitionStore.updateDefinition({ project_id, def_id, ...rest })
      } else {
        const res: any = await definitionStore.createDefinition({ project_id, ...data })
        def_id = res.id
        definition.value.id = res.id
        schemaTree.activeNode(res.id)
      }
    } finally {
      isSaving.value = false
    }
  }, 400),
  {
    deep: true,
  }
)

const onPreviewBtnClick = () => {
  if (!def_id) {
    return
  }

  goSchemaDetailPage(def_id)
}
</script>
