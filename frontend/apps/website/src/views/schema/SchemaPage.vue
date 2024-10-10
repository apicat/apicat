<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useNamespace } from '@apicat/hooks'
import { ElMessage, ClickOutside as vClickOutside } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { JSONSchemaTable } from '@apicat/components'
import type { PageModeCtx } from '../composables/usePageMode'
import { useIntelligentSchema } from '../composables/useIntelligentSchema'
import { useAITips } from './useAITips'
import useProjectStore from '@/store/project'
import useApi from '@/hooks/useApi'
import useDefinitionSchemaStore from '@/store/definitionSchema'
import { getSchemaHistoryPath } from '@/router/history'
import { injectPagesMode } from '@/layouts/ProjectDetailLayout/composables/usePagesMode'
import { useTitleInputFocus } from '@/hooks/useTitleInputFocus'
import { injectAsyncInitTask } from '@/hooks/useWaitAsyncTask'
import { apiParseSchema } from '@/api/project/definition/schema'

defineOptions({ inheritAttrs: false })
const props = defineProps<{ project_id: string, schemaID: string }>()

const GenerateCode = defineAsyncComponent(() => import('@/components/GenerateCode/GenerateCode.vue'))

const { t } = useI18n()
const ns = useNamespace('document')
const schemaIDRef = toRef(props, 'schemaID')

const { toggleMode, readonly } = injectPagesMode('schema') as PageModeCtx
const definitionSchemaStore = useDefinitionSchemaStore()
const { isManager, isWriter } = storeToRefs(useProjectStore())
const { schemas, schemaDetail: schema, isLoading: loading } = storeToRefs(definitionSchemaStore)
const [isSaving, updateSchema, isSaveError] = useApi(definitionSchemaStore.updateSchema)
const { inputRef: titleInputRef, focus } = useTitleInputFocus()
const router = useRouter()
const { handleIntelligentSchema, handleCheckReplaceModel } = useIntelligentSchema(props.project_id, () => {
  return {
    id: schema.value?.id,
    type: 'model',
    title: schema.value?.name,
  }
})
const { jsonSchemaTableIns, isAIMode, isShowAIStyle, preSchema, handleTitleBlur } = useAITips(props.project_id, schema, readonly, updateSchema)

let oldTitle = ''

function handleBlurNameInput() {
  const title = schema.value?.name || ''
  if (!title || !title.trim()) {
    schema.value!.name = oldTitle
    oldTitle = ''
  }
}

function goSchemahistory() {
  router.push({
    path: getSchemaHistoryPath(props.project_id, Number.parseInt(props.schemaID)),
    query: { ...router.currentRoute.value.query, iterationID: router.currentRoute.value.params.iterationID },
  })
}

watchDebounced(
  schema,
  async (n, o) => {
    if (readonly.value || !n)
      return

    // 还原旧的title时，不需要请求接口
    if (!oldTitle) {
      oldTitle = n?.name || ''
      return
    }

    if (n && o && n.id === o.id) {
      // title is empty
      if (!n.name || !n.name.trim())
        return ElMessage.error(t('app.schema.page.edit.titleNull'))

      // backup old title
      oldTitle = n.name

      try {
        !isAIMode.value && await updateSchema(props.project_id, n)
      }
      catch (error) {
        //
      }
    }
  },
  {
    deep: true,
    debounce: 500,
  },
)

async function setDetail(id: string) {
  oldTitle = ''
  const schemaID = Number.parseInt(id)
  if (!Number.isNaN(schemaID)) {
    await definitionSchemaStore.getSchemaDetail(props.project_id, schemaID)
    if (schema.value)
      preSchema.value = JSON.parse(JSON.stringify(schema.value))

    oldTitle = schema.value?.name || ''
    if (!readonly.value)
      focus()
  }
}

const generateCodeRef = ref<InstanceType<typeof GenerateCode>>()
watch(schemaIDRef, async (id, oID) => {
  if (id === oID)
    return
  generateCodeRef.value?.dispose()
  await setDetail(id)
})

injectAsyncInitTask()?.addTask(setDetail(schemaIDRef.value))
</script>

<template>
  <!-- head -->
  <div class="ac-header-operate">
    <div class="overflow-hidden ac-header-operate__main">
      <p v-if="readonly" class="truncate ac-header-operate__title" :title="schema?.name">
        {{ schema?.name }}
      </p>

      <p v-else class="flex-y-center" style="color: grey">
        <template v-if="isSaving">
          <Iconfont :size="18" class="mt-1px mr-4px" icon="ac-cloud-upload" />
          {{ $t('app.common.saving') }}
        </template>

        <template v-else-if="isSaveError">
          <Iconfont :size="18" class="mt-1px mr-4px" icon="ac-cloud-error" />
          {{ 'Save Error' }}
        </template>

        <template v-else>
          <Iconfont :size="18" class="mt-1px mr-4px" icon="ac-cloud-check" />
          {{ $t('app.common.savedCloud') }}
        </template>
      </p>
    </div>

    <div class="ac-header-operate__btns">
      <el-button v-if="isManager || isWriter" type="primary" @click="toggleMode">
        {{ !readonly ? $t('app.common.preview') : $t('app.common.edit') }}
      </el-button>

      <div v-if="readonly">
        <el-tooltip
          v-if="isManager || isWriter" effect="dark" placement="bottom"
          :content="$t('app.schema.history.title')" :show-arrow="false"
        >
          <Iconfont class="cursor-pointer ac-history" :size="24" @click="goSchemahistory" />
        </el-tooltip>
      </div>
    </div>
  </div>

  <!-- content -->
  <div v-loading="loading" :class="[ns.b(), ns.is('tips', isShowAIStyle)]">
    <div class="mb-10px">
      <h4 v-if="readonly && schema?.description" class="break-words">
        {{ schema?.description }}
      </h4>
      <div v-if="!readonly && schema">
        <input
          ref="titleInputRef" v-model="schema.name" v-click-outside="handleBlurNameInput"
          class="ac-document__title"
          type="text" maxlength="255" :placeholder="$t('app.schema.form.title')" @blur="handleTitleBlur"
        >
        <input
          v-model="schema.description" class="w-full ac-document__desc" type="text" maxlength="255"
          :placeholder="$t('app.schema.form.desc')"
        >
      </div>
    </div>

    <JSONSchemaTable
      v-if="!loading && schema" ref="jsonSchemaTableIns" :key="schema.id" v-model:schema="schema.schema"
      :readonly="readonly" :root-schema-key="schema.id" :definition-schemas="schemas"
      :handle-parse-schema="apiParseSchema" :handle-intelligent-schema="handleIntelligentSchema"
      :handle-check-replace-model="handleCheckReplaceModel"
    />
    <GenerateCode v-if="readonly && schema" ref="generateCodeRef" :schema="schema" />
  </div>
</template>
