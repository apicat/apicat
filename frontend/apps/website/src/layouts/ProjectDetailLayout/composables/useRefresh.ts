import { storeToRefs } from 'pinia'
import { useCollectionsStore } from '@/store/collections'
import useDefinitionResponseStore from '@/store/definitionResponse'
import useDefinitionSchemaStore from '@/store/definitionSchema'
import { useGlobalParameters } from '@/store/globalParameter'
import useProjectStore from '@/store/project'
import type { CollectionContext } from '@/hooks/useCollectionContext'
import { useGlobalServerUrlStore } from '@/store/globalServerUrl'

export function useRefresh(activeCollectionKey: Ref<number | undefined>, activeSchemaKey: Ref<number | undefined>, activeResponseKey: Ref<number | undefined>, ctx: CollectionContext) {
  const { projectID } = storeToRefs(useProjectStore())
  const collectionStore = useCollectionsStore()
  const schemaStore = useDefinitionSchemaStore()
  const responseStore = useDefinitionResponseStore()
  const globalParameters = useGlobalParameters()
  const globalServerUrlStore = useGlobalServerUrlStore()

  let collectionID: number | undefined

  const subscribes: Array<any> = []
  let unsubscribe: any

  // 删除response时，同步 collection
  unsubscribe = responseStore.$onAction(({ name, after }) => {
    if (name === 'deleteResponse' && activeCollectionKey.value) {
      after(async () => {
        if (!projectID.value)
          return

        // 刷新 response list
        await responseStore.getResponses(projectID.value!)

        // 停留在collection页面
        if (activeCollectionKey.value) {
          collectionStore.collectionDetail = null
          activeCollectionKey.value !== collectionID && await collectionStore.getCollectionDetail(projectID.value, activeCollectionKey.value)
        }
      })
    }
  })
  subscribes.push(unsubscribe)

  // 删除schema时，同步 collection | schema | response
  unsubscribe = schemaStore.$onAction(({ name, after }) => {
    if (name === 'deleteSchema') {
      after(async () => {
        if (!projectID.value)
          return

        // 刷新 schema list & response list
        await schemaStore.getSchemas(projectID.value!)
        await responseStore.getResponses(projectID.value!)

        // 停留在collection页面
        if (activeCollectionKey.value) {
          collectionStore.collectionDetail = null
          activeCollectionKey.value !== collectionID && await collectionStore.getCollectionDetail(projectID.value, activeCollectionKey.value)
        }

        // 停留在schema页面
        if (activeSchemaKey.value) {
          schemaStore.schemaDetail = null
          await schemaStore.getSchemaDetail(projectID.value, activeSchemaKey.value)
        }

        // 停留在response页面
        if (activeResponseKey.value)
          await responseStore.getResponseDetail(projectID.value, activeResponseKey.value)
      })
    }
  })
  subscribes.push(unsubscribe)

  // 删除globalParameter时，同步 collection
  unsubscribe = globalParameters.$onAction(({ name, after }) => {
    if (name === 'deleteGlobalParameter') {
      after(async () => {
        if (!projectID.value || !activeCollectionKey.value)
          return
        activeCollectionKey.value !== collectionID && await collectionStore.getCollectionDetail(projectID.value, activeCollectionKey.value)
      })
    }
  })
  subscribes.push(unsubscribe)

  // 删除｜创建全局URL
  globalServerUrlStore.$onAction(({ name, after }) => {
    if (name === 'deleteGlobalServerUrl' || name === 'createGlobalServerUrl') {
      after(() => {
        const url = globalServerUrlStore.urls.find(item => item.id === ctx.activeUrl.value)
        if (!url && globalServerUrlStore.urls.length)
          ctx.activeUrl.value = globalServerUrlStore.urls[0].id
      })
    }
  })

  // onUnmounted remove subscribes
  onUnmounted(() => {
    subscribes.forEach(unsubscribe => unsubscribe())
    unsubscribe = undefined
    collectionID = undefined
  })
}
