import DefaultAjax from '../Ajax'

export async function apiForgotSendEmail(email: string): Promise<ResponseAPI.Response<void>> {
  return DefaultAjax.post('/account/retrieve-password', { email })
}

export async function apiResetPass(data: { password: string, re_password: string }, code: string): Promise<ResponseAPI.Response<void>> {
  return DefaultAjax.put(`/account/reset-password/${code}`, data)
}

export async function checkResetPassCodeIsExpired(code: string) {
  return DefaultAjax.get(`/account/reset-password/check/${code}`, undefined, { isShowErrorMsg: false })
}
