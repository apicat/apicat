import type { CurrentNodeContext } from '../constants'
import { CurrentNodeContextKey } from '../constants'
import { useCollectionsStore } from '@/store/collections'
import {
  ITERATION_COLLECTION_PATH_NAME,
  ITERATION_RESPONSE_PATH_NAME,
  ITERATION_SCHEMA_PATH_NAME,
  PROJECT_COLLECTION_PATH_NAME,
  PROJECT_RESPONSE_PATH_NAME,
  PROJECT_SCHEMA_PATH_NAME,
} from '@/router'
import useDefinitionSchemaStore from '@/store/definitionSchema'
import useDefinitionResponseStore from '@/store/definitionResponse'
import { CollectionTypeEnum, ResponseTypeEnum, SchemaTypeEnum } from '@/commons'

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
      else currentNode.value = { id: undefined, type }
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

  function existNode<T>(
    nodes: T[],
    id: number | undefined,
    isLeaf: (node: T) => boolean,
    childKey: string = 'items',
  ): boolean {
    const getin = (items: any[]): boolean => {
      for (let i = 0; i < (items || []).length; i++) {
        const node = items[i]
        if (node.id === id) {
          return isLeaf(node)
        }
        else {
          if (!isLeaf(node)) {
            const exist = getin(node[childKey])
            if (exist)
              return exist
          }
        }
      }
      return false
    }
    return getin(nodes)
  }

  const collectionStore = useCollectionsStore()
  const schemaStore = useDefinitionSchemaStore()
  const responseStore = useDefinitionResponseStore()

  return {
    currentNode,
    activeCollectionKey,
    activeSchemaKey,
    activeResponseKey,
    currentNodeKey: computed(() => `${currentNode.value.id}-${currentNode.value.type}`),
    nodeExist: computed(() => {
      let data: any[]
      let isLeaf: any
      switch (currentNode.value.type) {
        case NODE_TYPE_COLLECTION:
          data = collectionStore.collections
          isLeaf = (node: CollectionAPI.ResponseCollection) => node.type !== CollectionTypeEnum.Dir
          break
        case NODE_TYPE_SCHEMA:
          data = schemaStore.schemas
          isLeaf = (node: Definition.SchemaNode) => node.type !== SchemaTypeEnum.Category
          break
        case NODE_TYPE_RESPONSE:
          data = responseStore.responses
          isLeaf = (node: Definition.ResponseTreeNode) => node.type !== ResponseTypeEnum.Category
          break
        default:
          return false
      }
      return existNode(data, currentNode.value.id, isLeaf)
    }),
  }
}

export function useActiveTreeContext(): CurrentNodeContext | undefined {
  return inject(CurrentNodeContextKey)
}
