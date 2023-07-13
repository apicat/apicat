/**
 * 成员在项目中的权限
 * 当前成员在此项目的权限:manage,write,read
 */
export const enum MemberAuthorityInProject {
  MANAGER = 'manage',
  WRITE = 'write',
  READ = 'read',
  NONE = 'none',
}

export const MemberAuthorityMap = {
  [MemberAuthorityInProject.MANAGER]: '管理',
  [MemberAuthorityInProject.WRITE]: '编辑',
  [MemberAuthorityInProject.READ]: '只读',
}

export declare interface ProjectMember {
  id?: number
  user_id?: number
  is_enabled?: number
  authority?: MemberAuthorityInProject
  username?: string
  email?: string
  created_at?: string

  isSelf?: boolean
  accountStatus?: string
  accountStatusType?: string
}
