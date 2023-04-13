import { DocumentTypeEnum, ProjectNavigateObject } from '@/commons/constant'
/**
 * project_detail
 */
export declare interface ProjectInfo {
  /**
   * 创建时间
   */
  created_at?: string
  /**
   * 项目描述
   */
  description?: string
  /**
   * 项目id
   */
  id: number | string
  /**
   * 项目名称
   */
  title: string
}

/**
 * server_detail
 */
export declare interface ServerDetail {
  /**
   * 描述
   */
  description: string
  /**
   * 地址
   */
  url: string
}

export declare interface CollectionNode {
  id: number
  parent_id?: number
  title: string
  name?: string
  type: DocumentTypeEnum
  sub_nodes?: CollectionNode[]
  _oldName?: string | undefined
  _extend?: {
    isLeaf: boolean
    isEditable: boolean
    isCurrent: boolean
  }
}

/**
 * schema 节点
 */
export interface SchemaNode {
  /**
   * 模型描述
   */
  description: string
  /**
   * 模型id
   */
  id: number
  /**
   * 模型名称
   */
  name: string
  /**
   * 类型
   * category
   * schema
   */
  type: string
}

export interface SchemaDetail extends SchemaNode {}

export type { ProjectNavigateObject }

/**
 * trash_list
 */
export interface TrashModel {
  /**
   * 删除时间
   */
  deleted_at: string
  /**
   * collection id
   */
  id: number
  /**
   * 标题
   */
  title: string
  /**
   * 类型
   * doc
   * http
   */
  type: DocumentTypeEnum
}
