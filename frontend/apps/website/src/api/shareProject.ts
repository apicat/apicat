import { QuietAjax } from './Ajax'
import { MemberAuthorityInProject } from '@/typings/member'
import { CookieOptions, Cookies } from '@/commons/cookie'

const restfulApipath = (project_id: string) => `/projects/${project_id}/share`

// 获取项目分享详情
export const getProjectShareDetail = async (project_id: string): Promise<{ authority: MemberAuthorityInProject; visibility: string; secret_key: string }> =>
  QuietAjax.get(restfulApipath(project_id))

// 项目当前分享状态
export const getProjectAuthInfo = async (project_id: string): Promise<{ authority: MemberAuthorityInProject; visibility: string; has_shared: boolean }> =>
  QuietAjax.get(`${restfulApipath(project_id)}/status`)
// 项目分享开关
export const switchProjectShareStatus = ({ project_id, ...params }: any) => QuietAjax.put(`${restfulApipath(project_id)}/switch`, params)
// 重置分享项目访问秘钥
export const resetSecretToProject = ({ project_id }: Record<string, any>) => QuietAjax.put(`${restfulApipath(project_id)}/reset`)
// 私有项目秘钥校验
export const checkProjectSecret = ({ project_id, secret_key }: Record<string, any>): Promise<{ token: string; expiration: string }> =>
  QuietAjax.post(`${restfulApipath(project_id)}/check`, { secret_key })
// 保存项目分享后的访问token
export const setProjectSharedToken = (project_id: string, token: string, options?: CookieOptions) => Cookies.set(`${Cookies.KEYS.SHARE_PROJECT}${project_id}`, token, options)
// 获取项目分享后的访问token
export const getProjectSharedToken = (project_id: string) => Cookies.get(`${Cookies.KEYS.SHARE_PROJECT}${project_id}`)
