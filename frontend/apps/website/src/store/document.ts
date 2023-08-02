import type { CollectionNode } from '@/typings/project'
import { DocumentTypeEnum } from '@/commons/constant'
import { traverseTree } from '@apicat/shared'
import { defineStore } from 'pinia'
import { getCollectionList, getDocumentHistoryRecordList } from '@/api/collection'

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
    documentHistoryRecordTree: [] as Array<CollectionNode>,
  }),
  getters: {
    historyRecordForOptions: (state) => {
      const options = state.documentHistoryRecordTree
        .reduce((result: any, item: any) => result.concat(item.sub_nodes || []), [])
        .map((item: any) => ({ id: item.id, title: item.title }))
      return [{ id: 0, title: '最新内容' }].concat(options)
    },
  },

  actions: {
    async getApiDocTree(project_id: string) {
      const tree = await getCollectionList(project_id)
      this.apiDocTree = traverseTree((item: CollectionNode) => extendDocTreeFeild(item), tree || [], { subKey: 'items' }) as Array<CollectionNode>
      return this.apiDocTree
    },

    async refreshApiDocTree(project_id: string) {
      const tree = await getCollectionList(project_id)
      this.apiDocTree = traverseTree((item: CollectionNode) => extendDocTreeFeild(item), tree || [], { subKey: 'items' }) as Array<CollectionNode>
      return this.apiDocTree
    },

    async getDocumentHistoryRecordList(project_id: string, collection_id: string) {
      if (!project_id || !collection_id) {
        return []
      }

      try {
        const tree = await getDocumentHistoryRecordList({ project_id, collection_id })
        this.documentHistoryRecordTree = traverseTree((item: CollectionNode) => extendDocTreeFeild(item), tree || [], { subKey: 'sub_nodes' }) as Array<CollectionNode>
      } catch (error) {
        return []
      }

      return this.documentHistoryRecordTree
    },
  },
})

export default useDocumentStore
