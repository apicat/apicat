<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useNamespace } from '@apicat/hooks'
import { ElMessage, ClickOutside as vClickOutside } from 'element-plus'
import { useI18n } from 'vue-i18n'
import type { PageModeCtx } from '../composables/usePageMode'
import ResponseForm from './components/ResponseForm.vue'
import ResponseRaw from './components/ResponseRaw.vue'
import useProjectStore from '@/store/project'
import useApi from '@/hooks/useApi'
import useDefinitionResponseStore from '@/store/definitionResponse'
import useDefinitionSchemaStore from '@/store/definitionSchema'
import { injectPagesMode } from '@/layouts/ProjectDetailLayout/composables/usePagesMode'
import { useTitleInputFocus } from '@/hooks/useTitleInputFocus'

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  project_id: string
  responseID: string
}>()

const { t } = useI18n()

const ns = useNamespace('document')
const { toggleMode, readonly } = injectPagesMode('response') as PageModeCtx
const { isManager, isWriter } = storeToRefs(useProjectStore())
const { schemas } = storeToRefs(useDefinitionSchemaStore())
const responseIDRef = toRef(props, 'responseID')
const definitionResponseStore = useDefinitionResponseStore()
const { loading, responseDetail: response } = storeToRefs(definitionResponseStore)
const [isSaving, updateResponse, isSaveError] = useApi(definitionResponseStore.updateResponse)
const { inputRef: titleInputRef, focus } = useTitleInputFocus()

let oldTitle = ''

function handleBlurNameInput() {
  const title = response.value.name || ''
  if (!title || !title.trim()) {
    response.value.name = oldTitle
    oldTitle = ''
  }
}
watchDebounced(
  response,
  async (n, o) => {
    if (readonly.value) return

    // 还原旧的title时，不需要请求接口
    if (!oldTitle) {
      oldTitle = n?.name || ''
      return
    }

    if (n && o && n.id === o.id) {
      // title is empty
      if (!n.name || !n.name.trim()) return ElMessage.error(t('app.definitionResponse.page.edit.titleNull'))

      // backup old title
      oldTitle = n.name

      try {
        await updateResponse(props.project_id, n)
      } catch (error) {
        //
      }
    }
  },
  { deep: true, debounce: 200 },
)

watch(
  responseIDRef,
  async (id, oID) => {
    if (id === oID) return
    oldTitle = ''
    const responseID = Number.parseInt(id)
    if (!Number.isNaN(responseID)) {
      await definitionResponseStore.getResponseDetail(props.project_id, responseID)
      oldTitle = response.value.name || ''
      if (!readonly.value) focus()
    }
  },
  {
    immediate: true,
  },
)
</script>

<template>
  <!-- head -->
  <div class="ac-header-operate">
    <div class="overflow-hidden ac-header-operate__main">
      <p v-if="readonly" class="truncate ac-header-operate__title" :title="response.name">
        {{ response.name }}
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

    <div v-if="isManager || isWriter" class="ac-header-operate__btns">
      <el-button type="primary" @click="toggleMode">
        {{ !readonly ? $t('app.common.preview') : $t('app.common.edit') }}
      </el-button>
    </div>
  </div>

  <!-- content -->
  <div v-loading="loading" :class="[ns.b()]">
    <div class="mb-10px">
      <h4 v-if="readonly && response.description" class="break-words">
        {{ response.description }}
      </h4>
      <div v-if="!readonly">
        <input
          ref="titleInputRef"
          v-model="response.name"
          v-click-outside="handleBlurNameInput"
          class="ac-document__title"
          type="text"
          maxlength="255"
          :placeholder="$t('app.schema.form.title')" />
        <input
          v-model="response.description"
          class="w-full ac-document__desc"
          type="text"
          maxlength="255"
          :placeholder="$t('app.schema.form.desc')" />
      </div>
    </div>
    <div v-if="!loading">
      <ResponseForm v-if="!readonly" v-model:response="response" :definition-schemas="schemas" up />
      <ResponseRaw v-else :response="response" :definition-schemas="schemas" />
    </div>
  </div>
</template>
