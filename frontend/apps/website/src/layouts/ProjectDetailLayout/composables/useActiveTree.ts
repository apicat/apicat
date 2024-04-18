import type { CurrentNodeContext } from '../constants'
import { CurrentNodeContextKey } from '../constants'
import {
  ITERATION_COLLECTION_PATH_NAME,
  ITERATION_RESPONSE_PATH_NAME,
  ITERATION_SCHEMA_PATH_NAME,
  PROJECT_COLLECTION_PATH_NAME,
  PROJECT_RESPONSE_PATH_NAME,
  PROJECT_SCHEMA_PATH_NAME,
} from '@/router'

export const NODE_TYPE_COLLECTION = 'collection'
export const NODE_TYPE_SCHEMA = 'schema'
export const NODE_TYPE_RESPONSE = 'response'

export function useActiveTree() {
  const route = useRoute()
  const currentNode = ref<{ id?: number; type: string }>({ id: undefined, type: 'none' })

  const activeCollectionKey = computed(() => {
    if (currentNode.value.type !== NODE_TYPE_COLLECTION)
      return undefined
    else return currentNode.value.id
  })

  const activeSchemaKey = computed(() => {
    if (currentNode.value.type !== NODE_TYPE_SCHEMA)
      return undefined
    else return currentNode.value.id
  })

  const activeResponseKey = computed(() => {
    if (currentNode.value.type !== NODE_TYPE_RESPONSE)
      return undefined
    else return currentNode.value.id
  })

  watchEffect(() => {
    // 路由
    try {
      const name = (route.name || '') as string
      let id: string | undefined
      let type: string = 'none'

      switch (name) {
        case PROJECT_COLLECTION_PATH_NAME:
        case ITERATION_COLLECTION_PATH_NAME:
          type = NODE_TYPE_COLLECTION
          id = (route.params.collectionID || '') as string
          break

        case PROJECT_SCHEMA_PATH_NAME:
        case ITERATION_SCHEMA_PATH_NAME:
          type = NODE_TYPE_SCHEMA
          id = (route.params.schemaID || '') as string
          break

        case PROJECT_RESPONSE_PATH_NAME:
        case ITERATION_RESPONSE_PATH_NAME:
          type = NODE_TYPE_RESPONSE
          id = (route.params.responseID || '') as string
          break

        default:
          id = undefined
          type = 'none'
      }

      if (id)
        currentNode.value = { id: Number.parseInt(id), type }
      else
        currentNode.value = { id: undefined, type }
    }
    catch (e) {
    //
    }
  })

  // 设置
  provide(CurrentNodeContextKey, {
    activeCollectionNode(id: number): void {
      currentNode.value = {
        id,
        type: NODE_TYPE_COLLECTION,
      }
    },
    activeSchemaNode(id: number): void {
      currentNode.value = {
        id,
        type: NODE_TYPE_SCHEMA,
      }
    },
    activeResponseNode(id: number): void {
      currentNode.value = {
        id,
        type: NODE_TYPE_RESPONSE,
      }
    },
    currentActiveNode: currentNode,
    activeCollectionKey,
    activeSchemaKey,
    activeResponseKey,
  })

  return {
    currentNode,
    activeCollectionKey,
    activeSchemaKey,
    activeResponseKey,
    currentNodeKey: computed(() => `${currentNode.value.id}-${currentNode.value.type}`),
  }
}

export function useActiveTreeContext(): CurrentNodeContext {
  return inject(CurrentNodeContextKey) || {} as any
}
