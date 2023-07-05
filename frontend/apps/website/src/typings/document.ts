import { DocumentTypeEnum } from '@/commons/constant'
import { RequestParameters } from './parameter'

/**
 * DocumentDetail
 */
export declare interface DocumentDetail {
  /**
   * 内容
   */
  content: string
  /**
   * 创建时间
   */
  created_at?: string
  /**
   * 创建人
   */
  created_by?: string
  /**
   * 集合id
   */
  id?: number
  /**
   * 父级id
   */
  parent_id?: number
  /**
   * 名称
   */
  title: string
  /**
   * 类型: category,doc,http
   */
  type: DocumentTypeEnum
  /**
   * 最后更新时间
   */
  updated_at?: string
  /**
   * 最后更新人
   */
  updated_by?: string
}

export interface ApiCatHttpRequestNode {
  id?: string | number
  type: string
  attrs: RequestParameters
}

export interface HttpDocument {
  id?: string | number
  parentid?: string | number
  tag?: Array<any>
  title: string
  type: string
  content: Array<any>
}

export interface SharedDocumentInfo {
  collection_id: string
  project_id: string
  doc_public_id?: string
}
