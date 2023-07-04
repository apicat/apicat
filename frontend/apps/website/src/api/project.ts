import Ajax, { QuietAjax } from './Ajax'
import useApi from '@/hooks/useApi'
import { MemberAuthorityInProject } from '@/typings/member'
import { ProjectInfo } from '@/typings/project'
import { API_URL } from '@/commons/constant'
import { queryStringify, Storage } from '@/commons'

export const getProjectList = () => Ajax.get('/projects')

export const getProjectDetail = (project_id: string) => Ajax.get(`/projects/${project_id}`)

export const createProject = async (projectInfo: Partial<ProjectInfo>): Promise<ProjectInfo> => await QuietAjax.post('/projects', projectInfo)

export const updateProjectBaseInfo = () => useApi(({ id: project_id, ...info }: ProjectInfo) => Ajax.put(`/projects/${project_id}`, info))

export const getProjectServerUrlList = (project_id: any) => Ajax.get(`/projects/${project_id}/servers`)

export const saveProjectServerUrlList = ({ project_id, urls }: any) => Ajax.put(`/projects/${project_id}/servers`, urls)

export const exportProject = ({ project_id, ...params }: any) => `${API_URL}/projects/${project_id}/data${queryStringify(params)}`

export const getProjectTranshList = () => useApi((project_id: string) => Ajax.get(`/projects/${project_id}/trashs`))

export const restoreDoc = () =>
  useApi(({ project_id, ids }: any) => Ajax.put(`/projects/${project_id}/trashs?${ids.map((id: any) => `collection-id=${id}`).join('&')}`, { category: 0 }))

export const deleleProject = () => useApi((project_id: string) => Ajax.delete(`/projects/${project_id}`))

// 获取非此项目的成员列表
export const getMembersWithoutProject = (project_id: string) => QuietAjax.get(`/projects/${project_id}/members/without`)
// 获取成员列表
export const getMembersInProject = (project_id: string) => (params: Record<string, any>) => Ajax.get(`/projects/${project_id}/members${queryStringify(params)}`)
// 新增成员
export const addMemberToProject = (project_id: string) => (params: Record<string, any>) => Ajax.post(`/projects/${project_id}/members`, params)
// 删除成员
export const deleteMemberFromProject = (project_id: string, user_id: number) => Ajax.delete(`/projects/${project_id}/members/${user_id}`)
// 修改成员权限
export const updateMemberAuthorityInProject = async (project_id: string, user_id: number, authority: MemberAuthorityInProject) =>
  QuietAjax.put(`/projects/${project_id}/members/authority/${user_id}`, { authority })
// 退出项目
export const quitProject = (project_id: string) => Ajax.delete(`/projects/${project_id}/exit`)
// 移交项目
export const transferProject = (project_id: string, member_id: number) => Ajax.put(`/projects/${project_id}/transfer`, { member_id })

// 获取项目分享详情
export const getProjectShareDetail = (project_id: string) => QuietAjax.get(`/projects/${project_id}/status`)
// 项目当前状态
export const getProjectStatus = getProjectShareDetail
// 重置分享项目访问秘钥
export const resetSecretToProject = ({ project_id }: Record<string, any>) => QuietAjax.put(`/projects/${project_id}/share/reset_share_secretkey`)
// 项目分享开关
export const switchProjectShareStatus = ({ project_id, ...params }: any) => Ajax.put(`/projects/${project_id}/share`, params)
// 私有项目秘钥校验
export const checkProjectSecret = ({ project_id, secret_key }: Record<string, any>) => QuietAjax.post(`/projects/${project_id}/share/secretkey_check`, { secret_key })

// 保存项目分享后的访问token
export const setProjectSharedToken = (project_id: string, token: string) => Storage.set(`${Storage.KEYS.SHARE_PROJECT}${project_id}`, token, true)
// 获取项目分享后的访问token
export const getProjectSharedToken = (project_id: string) => Storage.get(`${Storage.KEYS.SHARE_PROJECT}${project_id}`, true)
