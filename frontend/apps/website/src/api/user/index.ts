import DefaultAjax, { RawAjax } from '@/api/Ajax'

export async function apiGetAvatar(src: string): Promise<Blob> {
  return RawAjax.get(src)
}

export async function apiGetUserInfo(): Promise<UserAPI.ResponseUserInfo> {
  return DefaultAjax.get('/user')
}

// General Page: name & language
export async function apiUpdateGeneral(data: UserAPI.RequestGeneral): Promise<void> {
  return DefaultAjax.put('/user', <UserAPI.RequestGeneral>{
    name: data.name,
    language: data.language,
  })
}

//  General Page: avatar upload
export async function apiUploadAvatar(data: UserAPI.RequestChangeAvatar): Promise<UserAPI.ResponseChangeAvatar> {
  const form = new FormData()
  for (const key in data)
    form.append(key, data[key as keyof UserAPI.RequestChangeAvatar] as File | string)

  return DefaultAjax.post('/user/avatar', form)
}

// Email Page
export async function apiUpdateEmail(email: string): Promise<void> {
  return DefaultAjax.put('/user/email', { email })
}

// OAuth Page: Connect
export async function apiConnectOAuth(platform: SignAPI.OAuthPlatform, data: { code: string }): Promise<void> {
  return DefaultAjax.post(`/user/oauth/${platform}/connect`, data)
}

// OAuth Page: Disconnect
export async function apiDisconnectOAuth(platform: SignAPI.OAuthPlatform): Promise<void> {
  return DefaultAjax.delete(`/user/oauth/${platform}/disconnect`)
}

// Reset Password Page
export async function apiResetPassword(data: UserAPI.RequestResetPassword): Promise<void> {
  return DefaultAjax.put('/user/password', data)
}

// 获取系统用户列表
export async function apiGetSystemUserList(params?: Record<string, any>): Promise<GlobalAPI.ResponseTable<UserAPI.ResponseUserInfo[]>> {
  return DefaultAjax.get('/users', { params })
}

// 修改系统用户密码
export async function apiChangeSystemUserPassword(data: { id: number, password: string, confirmPassword: string }): Promise<void> {
  const { id, ...info } = data
  return DefaultAjax.patch(`/users/${id}`, { ...info }, { isShowSuccessMsg: true })
}

// 删除系统用户
export async function apiDeleteSystemUser(id: number): Promise<void> {
  return DefaultAjax.delete(`/users/${id}`, {}, { isShowSuccessMsg: true })
}
