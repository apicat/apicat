import { DocumentTypeEnum, ProjectNavigateObject } from '@/commons/constant'
import { MemberAuthorityInProject } from './member'

export declare interface ProjectCover {
  type: 'icon' | 'url'
  coverBgColor?: string
  coverIcon?: string
  coverUrl?: string
}

/**
 * 项目信息
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
  /**
   * 项目封面
   */
  cover?: string | ProjectCover
  /**
   * 项目权限
   */
  visibility?: 'private' | 'public'

  /**
   * 当前成员在此项目的权限
   */
  authority?: MemberAuthorityInProject

  /**
   * 是否关注
   */
  is_followed?: boolean
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
  items?: CollectionNode[]
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
export declare interface SchemaNode {
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
export declare interface TrashModel {
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
