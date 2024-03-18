import DefaultAjax from '../Ajax'

export async function apiLogin(data: SignAPI.RequestLogin): Promise<SignAPI.ResponseLogin> {
  return DefaultAjax.post('/account/login', data)
}

export async function apiRegister(data: SignAPI.RequestRegister): Promise<SignAPI.ResponseRegister> {
  return DefaultAjax.post('/account/register', data)
}

export async function activeAccountByEmail(code: string): Promise<ResponseAPI.Response<void>> {
  return DefaultAjax.put(`/account/email-verification/${code}`, undefined, {
    isShowSuccessMsg: false,
    isShowErrorMsg: false,
  })
}

export async function changeUserEmail(code: string) {
  return DefaultAjax.put(`/user/email/${code}`, undefined, { isShowErrorMsg: false })
}
