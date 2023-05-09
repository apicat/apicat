export interface UserInfo {
  id?: number
  email?: string
  username?: string
  password?: string
}

export interface LoginResponse {
  access_token: string
  user: UserInfo
}
