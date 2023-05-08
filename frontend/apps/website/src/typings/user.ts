export interface UserInfo {
  id?: number
  email: string
  name?: string
  password?: string
}

export interface LoginResponse {
  access_token: string
  user: UserInfo
}
