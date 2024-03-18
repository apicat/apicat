declare namespace HistoryRecord {
  interface TreeNode {
    id: number | string
    title: string
    items?: TreeNode[]
  }

  type CollectionTypeEnum = import('@/commons/constant').CollectionTypeEnum

  /**
   * 历史记录列表（用于目录显示）
   */
  interface CollectionTreeNode {
    id: number | string
    title: string
    items?: CollectionTreeNode[]
    type?: CollectionTypeEnum
  }

  /**
   * 历史记录列表（后端响应数据）
   */
  type ResponseCollectionRecordList = CollectionInfo[]

  /**
   * 历史记录信息
   */
  interface CollectionInfo {
    id: number
    createdAt: number
    createdBy: number
    type: CollectionTypeEnum
  }

  interface CollectionInfoForOptions {
    id: number
    title: string
  }

  /**
   * 集合历史记录详情
   */
  interface CollectionDetail {
    id: number
    collectionID: number
    content: Array<Record<string, any>>
    createdAt: string
    createdBy: string
    title: string
    updatedAt: string
    type: CollectionTypeEnum
  }

  /**
   * 文档对比结果
   */
  interface CollectionDiff {
    doc1: CollectionDetail
    doc2: CollectionDetail
  }

  type SchemaTypeEnum = import('@/commons/constant').SchemaTypeEnum

  interface SchemaHistoryNode extends TreeNode {
    type?: SchemaTypeEnum
  }

  interface SchemaHistory {
    createdAt: number
    id: number
    type: SchemaTypeEnum
    createdBy: string
  }

  interface SchemaHistoryInfo {
    createdAt: string
    createdBy: string
    description: string
    id: number
    name: string
    schema: string
    schemaID: number
    type: SchemaTypeEnum
    updatedAt: string
  }

  interface SchemaHistoryDiff {
    schema1: SchemaHistoryInfo
    schema2: SchemaHistoryInfo
  }

  interface SchemaInfoForOptions {
    id: number
    title: string
  }
}
