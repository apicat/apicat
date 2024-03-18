import Ajax, { DefaultAjax } from '@/api/Ajax'

// project
export function apiGetProjectShareStatus(
  projectID: string,
): Promise<ShareAPI.ResponseProjectShareStatus> {
  return Ajax.get(`/projects/${projectID}/share/status`)
}
export function apiGetProjectShareInfo(
  projectID: string,
): Promise<ShareAPI.ResponseProjectShareInfo> {
  return Ajax.get(`/projects/${projectID}/share`)
}
export function apiChangeProjectShareStatus(
  projectID: string,
  status: boolean,
): Promise<ShareAPI.ResponseChangeProjectShareStatus> {
  return Ajax.put(`/projects/${projectID}/share`, { projectID, status })
}
export function apiResetProjectShareKey(
  projectID: string,
): Promise<ShareAPI.ResponseResetProjectShareKey> {
  return Ajax.put(`/projects/${projectID}/share/reset`)
}
export function apiSendProjectSharekey(
  projectID: string,
  secretKey: string,
): Promise<ShareAPI.ResponseCheckProjectShareKey> {
  return Ajax.post(`/projects/${projectID}/share/check`, {
    projectID,
    secretKey,
  })
}

// document
export function apiGetDocShareStatus(
  collectionPublicID: string,
): Promise<ShareAPI.ResponseDocShareStatus> {
  return DefaultAjax.get(`/collections/${collectionPublicID}/share/status`, undefined, { isShowErrorMsg: false })
}
export function apiGetDocShareInfo(
  projectID: string,
  collectionID: number,
): Promise<ShareAPI.ResponseDocShareInfo> {
  return DefaultAjax.get(
    `/projects/${projectID}/collections/${collectionID}/share`,
  )
}
export function apiChangeDocShareStatus(
  projectID: string,
  collectionID: number,
  status: boolean,
): Promise<ShareAPI.ResponseChangeDocShareStatus> {
  return DefaultAjax.put(
    `/projects/${projectID}/collections/${collectionID}/share`,
    { status },
  )
}
export function apiResetDocShareKey(
  projectID: string,
  collectionID: number,
): Promise<ShareAPI.ResponseResetDocShareStatus> {
  return DefaultAjax.put(
    `/projects/${projectID}/collections/${collectionID}/share/reset`,
  )
}
export function apiSendDocShareKey(
  projectID: string,
  collectionID: number,
  secretKey: string,
): Promise<ShareAPI.ResponseSendShareKey> {
  return DefaultAjax.post(
    `/projects/${projectID}/collections/${collectionID}/share/check`,
    { secretKey },
  )
}
