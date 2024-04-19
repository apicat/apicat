import DefaultAjax from '../Ajax'

// service
export async function apiGetService(): Promise<SystemAPI.ServiceData> {
  return DefaultAjax.get('/sysconfigs/service')
}
export async function apiUpdateService(data: SystemAPI.ServiceData): Promise<void> {
  return DefaultAjax.put('/sysconfigs/service', data, { isShowSuccessMsg: true })
}

// oauth
export async function apiGetOAuth(): Promise<SystemAPI.OAuthData> {
  return DefaultAjax.get('/sysconfigs/oauth')
}
export async function apiUpdateOAuth(data: SystemAPI.OAuthData): Promise<void> {
  return DefaultAjax.put('/sysconfigs/oauth', data, { isShowSuccessMsg: true })
}

// storage
export async function apiGetStorage(): Promise<SystemAPI.StorageItem[]> {
  return DefaultAjax.get('/sysconfigs/storages')
}
export async function apiUpdateStroageLocal(data: SystemAPI.StorageDisk): Promise<void> {
  return DefaultAjax.put('/sysconfigs/storages/disk', data, { isShowSuccessMsg: true })
}
export async function apiUpdateStroageCF(data: SystemAPI.StorageCF): Promise<void> {
  return DefaultAjax.put('/sysconfigs/storages/cloudflare', data)
}
export async function apiUpdateStroageQiniu(data: SystemAPI.StorageQiniu): Promise<void> {
  return DefaultAjax.put('/sysconfigs/storages/qiniu', data, { isShowSuccessMsg: true })
}

// email
export async function apiGetEmail(): Promise<SystemAPI.EmailItem[]> {
  return DefaultAjax.get('/sysconfigs/emails')
}
export async function apiUpdateEmailSMTP(data: SystemAPI.EmailSMTP): Promise<void> {
  return DefaultAjax.put('/sysconfigs/emails/smtp', data, { isShowSuccessMsg: true })
}
export async function apiUpdateEmailSendCloud(data: SystemAPI.EmailSendCloud): Promise<void> {
  return DefaultAjax.put('/sysconfigs/emails/sendcloud', data, { isShowSuccessMsg: true })
}

// model
export async function apiGetModel(): Promise<SystemAPI.ModelItem[]> {
  return DefaultAjax.get('/sysconfigs/models')
}
export async function apiUpdateModelOpenAI(data: SystemAPI.ModelOpenAI): Promise<void> {
  return DefaultAjax.put('/sysconfigs/models/openai', data, { isShowSuccessMsg: true })
}
export async function apiUpdateModelAzure(data: SystemAPI.ModelAzure): Promise<void> {
  return DefaultAjax.put('/sysconfigs/models/azure-openai', data, { isShowSuccessMsg: true })
}
