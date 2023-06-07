/**
 * 成员在项目中的权限
 * 当前成员在此项目的权限:manage,write,read
 */
export const enum MemberAuthorityInProject {
  MANAGER = 'manage',
  WRITE = 'write',
  READ = 'read',
}

export declare interface ProjectMember {
  /**
   * 成员id
   */
  id?: number
  /**
   * 用户id
   */
  user_id?: number
  authority: MemberAuthorityInProject
  username?: string
  created_at?: string
}
