import { QuietAjax } from './Ajax'
import { CookieOptions, Cookies } from '@/commons/cookie'
import { CheckDocumentSecretParams, SharedDocumentInfo } from '@/typings/document'

const shareRestfulPath = (project_id: string | number, collection_id: string | number): string => `/projects/${project_id}/collections/${collection_id}/share`

// 文档分享状态
export const getCollectionShareStatus = async (doc_public_id: string): Promise<SharedDocumentInfo | null> => QuietAjax.get(`/collections/${doc_public_id}/share/status`)
// 文档分享详情
export const getCollectionShareDetail = async ({ project_id, collection_id }: any): Promise<{ visibility: string; collection_public_id: string; secret_key: string }> =>
  QuietAjax.get(`${shareRestfulPath(project_id, collection_id)}`)

// 文档分享开关
export const switchCollectionShareStatus = async ({ project_id, collection_id, ...params }: any): Promise<{ collection_public_id: string; secret_key: string }> =>
  QuietAjax.put(`${shareRestfulPath(project_id, collection_id)}/switch`, params)

// 重置文档分享密钥
export const resetSecretToCollection = async ({ project_id, collection_id }: any): Promise<{ secret_key: string }> =>
  QuietAjax.put(`${shareRestfulPath(project_id, collection_id)}/reset`)

// 文档秘钥校验
export const checkCollectionSecret = async ({ project_id, collection_id, secret_key }: CheckDocumentSecretParams): Promise<{ token: string; expiration: string }> =>
  QuietAjax.post(`${shareRestfulPath(project_id, collection_id)}/check`, { secret_key })

// 保存文档分享后的访问token
export const setCollectionSharedToken = (doc_public_id: string, token: string, options?: CookieOptions) =>
  Cookies.set(`${Cookies.KEYS.SHARE_DOCUMENT}${doc_public_id}`, token, options)

// 获取文档分享后的访问token
export const getCollectionSharedToken = (doc_public_id: string): string => Cookies.get(`${Cookies.KEYS.SHARE_DOCUMENT}${doc_public_id}`)
