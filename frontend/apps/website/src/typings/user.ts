/**
 * 用户在团队中的权限
 * 权限：superadmin,admin,user
 */
export const enum UserRoleInTeam {
  SUPER_ADMIN = 'superadmin',
  ADMIN = 'admin',
  USER = 'user',
}

export const UserRoleInTeamMap = {
  [UserRoleInTeam.SUPER_ADMIN]: '超级管理员',
  [UserRoleInTeam.ADMIN]: '管理员',
  [UserRoleInTeam.USER]: '普通用户',
}

export declare interface UserInfo {
  id?: number
  email?: string
  role: UserRoleInTeam
  username?: string
  password?: string
  start_using?: 0 | 1
  is_enabled?: 0 | 1
  created_at?: string
  updated_at?: string

  accountStatus?: string
  accountStatusType?: string
  isSelf?: boolean
  isSuperAdmin?: boolean
}

export declare interface LoginResponse {
  access_token: string
  user: UserInfo
}
