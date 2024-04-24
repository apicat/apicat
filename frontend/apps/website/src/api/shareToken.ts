import { Storage } from '@/commons'
import router, { COLLECTION_SHARE_PATH_NAME, PROJECT_DETAIL_PATH_NAME } from '@/router'

// DOC ------------------------------------------------
export function setCollectionSharedToken(id: string, token: string) {
  return Storage.set(`${Storage.KEYS.DOC_SHARE_TOKEN}${id}`, token, true)
}
export function getCollectionSharedToken(id: string): string {
  return Storage.get(`${Storage.KEYS.DOC_SHARE_TOKEN}${id}`, true)
}
export function clearCollectionSharedToken(id: string) {
  return Storage.remove(`${Storage.KEYS.DOC_SHARE_TOKEN}${id}`, true)
}

// PROJECT ------------------------------------------------------
// 保存项目分享后的访问token
export function setProjectSharedToken(id: string, token: string) {
  return Storage.set(`${Storage.KEYS.PROJECT_SHARE_TOKEN}${id}`, token, true)
}
// 获取项目分享后的访问token
export function getProjectSharedToken(id: string) {
  return Storage.get(`${Storage.KEYS.PROJECT_SHARE_TOKEN}${id}`, true)
}
export function clearProjectSharedToken(id: string) {
  return Storage.remove(`${Storage.KEYS.PROJECT_SHARE_TOKEN}${id}`, true)
}

// 项目和文档共用的token
export function gatherSharedTokenWithParams(params?: Record<string, any>, projectID?: string) {
  params = params || {}

  const currentRouteMatched = router.currentRoute.value.matched
  const routerParams = router.currentRoute.value.params
  const { collectionPublicID } = routerParams as Record<string, string>

  // 预览分享的文档
  if (currentRouteMatched.find(route => route.name === COLLECTION_SHARE_PATH_NAME)) {
    params.shareCode = getCollectionSharedToken(collectionPublicID)
    return params
  }

  // 预览分享的项目
  if (projectID) {
    params.shareCode = getProjectSharedToken(projectID)
    return params
  }

  return params
}

// clear share token
export function clearShareToken() {
  const currentRouteMatched = router.currentRoute.value.matched
  const routerParams = router.currentRoute.value.params
  const { collectionPublicID, project_id } = routerParams as Record<string, string>

  // 预览分享的文档
  if (currentRouteMatched.find(route => route.name === COLLECTION_SHARE_PATH_NAME))
    clearCollectionSharedToken(collectionPublicID)

  // 项目分享
  if (currentRouteMatched.find(route => route.name === PROJECT_DETAIL_PATH_NAME)) {
    project_id && clearProjectSharedToken(project_id)
    setTimeout(() => location.reload(), 1000)
  }
}
