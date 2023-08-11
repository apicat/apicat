import { MemberAuthority } from './member'

export interface Iteration {
  /**
   * 涉及的api数量
   */
  api_num: number
  /**
   * 当前成员在此项目的权限:manage-管理员,write-写,read-读,none-无权限
   */
  authority: MemberAuthority
  /**
   * 创建日期
   */
  created_at: string
  /**
   * 创建人
   */
  created_by: string
  /**
   * 迭代描述
   */
  description: string
  /**
   * 迭代id
   */
  id: number
  /**
   * 是否关注了此项目 True or False
   */
  is_followed: boolean
  /**
   * 项目公开id
   */
  project_public_id: string
  /**
   * 项目标题
   */
  project_title: string
  /**
   * 迭代公开id
   */
  public_id: string
  /**
   * 迭代标题
   */
  title: string
  /**
   * 迭代ids
   */
  collection_ids?: number[]

  [property: string]: any
}

export type SelectedProjectKey = number | string | null
