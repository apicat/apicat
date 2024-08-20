<script setup lang="ts">
import '@apicat/editor/dist/style.css'
import { storeToRefs } from 'pinia'
import { useNamespace } from '@apicat/hooks'
import { useIntelligentSchema } from '../composables/useIntelligentSchema'
import { useCollection } from './useCollection'
import useProjectStore from '@/store/project'
import { useCollectionContext } from '@/hooks/useCollectionContext'
import CollectionTestPage from '@/views/collection/CollectionTestPage.vue'

defineOptions({ inheritAttrs: false })
const props = defineProps<{ project_id: string, collectionID?: string }>()
const AcEditor = defineAsyncComponent(() => import('@apicat/editor'))

const ns = useNamespace('document')
const { activeUrl, urls, parameters, responses, schemas, acEditorOptions: options } = useCollectionContext()
const { isManager, isWriter, isPrivate, isReader, mockURL } = storeToRefs(useProjectStore())
const {
  collection,
  isLoading,
  isSaving,
  isSaveError,
  handleTitleBlur,
  handleContentUpdate,
  onShareCollectionBtnClick,
  onExportCollectionBtnClick,
  goCollectionhistory,
  readonly,
  toggleMode,
  titleInputRef,
} = useCollection(props as any)

const { handleIntelligentSchema } = useIntelligentSchema(() => {
  return {
    collectionID: collection?.value?.id,
    title: collection?.value?.title,
  }
})

options.handleIntelligentSchema = handleIntelligentSchema

const testPageRef = ref<InstanceType<typeof CollectionTestPage>>()
function showTestPage() {
  testPageRef.value?.show()
}
</script>

<template>
  <div v-if="collectionID">
    <!-- head -->
    <div v-if="collection" class="ac-header-operate">
      <div class="overflow-hidden ac-header-operate__main">
        <p v-if="readonly" class="truncate ac-header-operate__title" :title="collection?.title">
          {{ collection?.title }}
        </p>
        <p v-else class="flex-y-center" style="color: grey">
          <template v-if="isSaving">
            <Iconfont :size="18" class="mt-1px mr-4px" icon="ac-cloud-upload" />
            {{ $t('app.save.saving') }}
          </template>

          <template v-else-if="isSaveError">
            <Iconfont :size="18" class="mt-1px mr-4px" icon="ac-cloud-error" />
            {{ $t('app.save.error') }}
          </template>

          <template v-else>
            <Iconfont :size="18" class="mt-1px mr-4px" icon="ac-cloud-check" />
            {{ $t('app.save.saved') }}
          </template>
        </p>
      </div>

      <div class="ac-header-operate__btns">
        <el-button v-if="isManager || isWriter" type="primary" @click="toggleMode">
          {{ !readonly ? $t('app.common.preview') : $t('app.common.edit') }}
        </el-button>
        <template v-if="isManager || isWriter">
          <el-tooltip
            effect="dark" placement="bottom" :content="$t('app.project.collection.test.title')"
            :show-arrow="false"
          >
            <Iconfont icon="ac-shiyan cursor-pointer font-500" :size="18" @click="showTestPage" />
          </el-tooltip>
          <el-tooltip
            effect="dark" placement="bottom" :content="$t('app.project.collection.page.sharedoc')"
            :show-arrow="false"
          >
            <Iconfont icon="ac-share cursor-pointer" :size="18" @click="onShareCollectionBtnClick" />
          </el-tooltip>
        </template>

        <template v-else>
          <el-tooltip
            v-if="!isPrivate && isReader" :show-arrow="false"
            :content="$t('app.project.collection.page.sharedoc')" effect="dark" placement="bottom"
          >
            <Iconfont icon="ac-share cursor-pointer " :size="18" @click="onShareCollectionBtnClick" />
          </el-tooltip>
        </template>

        <el-tooltip
          v-if="isManager || isWriter" :content="$t('app.project.collection.page.exportdoc')"
          :show-arrow="false" effect="dark" placement="bottom"
        >
          <Iconfont icon="ac-export cursor-pointer" :size="18" @click="onExportCollectionBtnClick" />
        </el-tooltip>
        <el-tooltip
          v-if="isManager || isWriter" effect="dark" placement="bottom"
          :content="$t('app.project.collection.page.history')" :show-arrow="false"
        >
          <Iconfont class="cursor-pointer ac-history" :size="24" @click="goCollectionhistory" />
        </el-tooltip>
      </div>
    </div>

    <!-- content -->
    <div v-if="collection" v-loading="isLoading" :class="[ns.b()]">
      <div v-if="!readonly">
        <input
          ref="titleInputRef" v-model="collection.title" class="ac-document__title" type="text" maxlength="255"
          :placeholder="$t('app.schema.form.title')" @blur="handleTitleBlur"
        >
      </div>

      <AcEditor
        v-if="!isLoading" v-model:active-url="activeUrl" :mock-url="mockURL" :readonly="readonly"
        :content="collection.content!" :urls="urls" :schemas="schemas" :responses="responses" :parameters="parameters"
        :options="options" @update="handleContentUpdate"
      />
    </div>

    <CollectionTestPage
      v-if="(isManager || isWriter) && collectionID" ref="testPageRef" :project-i-d="project_id"
      :collection-i-d="collectionID"
    />
  </div>
  <el-empty v-else />
</template>
