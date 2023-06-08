/**
 * 成员在项目中的权限
 * 当前成员在此项目的权限:manage,write,read
 */
export const enum MemberAuthorityInProject {
  MANAGER = 'manage',
  WRITE = 'write',
  READ = 'read',
}

export const MemberAuthorityMap = {
  [MemberAuthorityInProject.MANAGER]: '管理',
  [MemberAuthorityInProject.WRITE]: '可写',
  [MemberAuthorityInProject.READ]: '可读',
}

export declare interface ProjectMember {
  id?: number
  user_id?: number
  authority?: MemberAuthorityInProject
  username?: string
  email?: string
  created_at?: string

  isSelf?: boolean
}
