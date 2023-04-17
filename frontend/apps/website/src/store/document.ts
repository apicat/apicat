import type { CollectionNode } from '@/typings/project'
import { DocumentTypeEnum } from '@/commons/constant'
import { traverseTree } from '@apicat/shared'
import { defineStore } from 'pinia'
import { getCollectionList } from '@/api/collection'

export const extendDocTreeFeild = (node = {} as CollectionNode, type = DocumentTypeEnum.HTTP): CollectionNode => {
  node = node || {}
  node.type = node.type === undefined ? type : node.type
  node._extend = {
    isLeaf: node.type !== DocumentTypeEnum.DIR,
    isEditable: false,
    isCurrent: false,
  }

  return node
}

export const useDocumentStore = defineStore('document', {
  state: () => ({
    apiDocTree: [] as Array<CollectionNode>,
    tempCreateDocParentId: undefined as number | undefined,
  }),
  actions: {
    async getApiDocTree(project_id: string) {
      const tree = await getCollectionList(project_id)
      this.apiDocTree = traverseTree((item: CollectionNode) => extendDocTreeFeild(item), tree || [], { subKey: 'sub_nodes' }) as Array<CollectionNode>
      return this.apiDocTree
    },

    async refreshApiDocTree(project_id: string) {
      const tree = await getCollectionList(project_id)
      this.apiDocTree = traverseTree((item: CollectionNode) => extendDocTreeFeild(item), tree || [], { subKey: 'sub_nodes' }) as Array<CollectionNode>
      return this.apiDocTree
    },
  },
})

export default useDocumentStore
