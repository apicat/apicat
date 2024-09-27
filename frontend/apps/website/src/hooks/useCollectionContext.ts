import type { AcEditorOptions } from '@apicat/editor'
import { languages } from '@codemirror/language-data'
import { ElMessage } from 'element-plus'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import useDefinitionSchemaStore from '@/store/definitionSchema'
import useDefinitionResponseStore from '@/store/definitionResponse'
import { useGlobalParameters } from '@/store/globalParameter'
import { useGlobalServerUrlStore } from '@/store/globalServerUrl'
import { getMockData } from '@/api/mock'
import { apiParseSchema } from '@/api/project/definition/schema'

export interface CollectionContext {
  urls: Ref<ProjectAPI.ResponseURL[]>
  parameters: Ref<ProjectAPI.ResponseGlobalParamList | any>
  responses: Ref<Definition.ResponseTreeNode[] | any>
  schemas: Ref<Definition.SchemaNode[] | any>
  acEditorOptions: AcEditorOptions
  activeUrl: Ref<number | string | null>
}

interface CollectionProviderOptions {
  projectID?: string
}

export const CollectionContextKey: InjectionKey<CollectionContext> = Symbol('CollectionContextKey')

export function useCollectionProvider({ projectID }: CollectionProviderOptions, cb?: () => Promise<void> | void) {
  const { context: collectionContext, initContextData } = useCollectionContextWithoutMounted()
  const { activeUrl } = collectionContext

  onMounted(async () => await initContextData(projectID, cb))

  onUnmounted(() => {
    activeUrl.value = null
  })
  return collectionContext
}

export function useCollectionContextWithoutMounted(): { context: CollectionContext, initContextData: (_projectID?: string, cb?: () => Promise<void> | void) => Promise<void> } {
  const { t } = useI18n()
  const definitionSchemaStore = useDefinitionSchemaStore()
  const definitionResponseStore = useDefinitionResponseStore()
  const globalParametersStore = useGlobalParameters()
  const globalServerUrlStore = useGlobalServerUrlStore()

  // 当前选中的url
  const activeUrl = ref()

  const acEditorOptions: AcEditorOptions = {
    codeBlockLanguages: languages as any,
    onCopySuccess: () => ElMessage.success(t('app.project.collection.copy.copied')),
    handleMockData: async (url, method, data) => await getMockData(url, method, data),
    handleParseSchema: async schema => await apiParseSchema(schema),
  }

  const collectionContext = {
    urls: storeToRefs(globalServerUrlStore).urls,
    parameters: storeToRefs(globalParametersStore).parameters,
    responses: storeToRefs(definitionResponseStore).responses,
    schemas: storeToRefs(definitionSchemaStore).schemas,
    acEditorOptions,
    activeUrl,
  }

  provide(CollectionContextKey, collectionContext)

  async function initContextData(_projectID?: string, cb?: () => Promise<void> | void) {
    if (!_projectID) {
      console.warn('useCollectionProvider:projectID is required')
      return
    }

    await Promise.all([
      globalServerUrlStore.getGlobalServerUrlList(_projectID),
      globalParametersStore.getGlobalParameterList(_projectID),
      definitionSchemaStore.getSchemas(_projectID),
      definitionResponseStore.getResponses(_projectID),
    ])

    // 设置默认选中的url
    if (globalServerUrlStore.urls.length)
      activeUrl.value = globalServerUrlStore.urls[0].id

    else
      activeUrl.value = null

    await cb?.()
  }

  return {
    context: collectionContext,
    initContextData,
  }
}

export function useCollectionContext(): CollectionContext
export function useCollectionContext() {
  return inject(CollectionContextKey)
}
