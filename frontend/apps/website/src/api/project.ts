import Ajax, { QuietAjax } from './Ajax'
import useApi from '@/hooks/useApi'
import { MemberAuthorityInProject } from '@/typings/member'
import { ProjectInfo } from '@/typings/project'

export const getProjectList = () => Ajax.get('/projects')

export const getProjectDetail = (project_id: string) => Ajax.get(`/projects/${project_id}`)

export const createProject = async (projectInfo: Partial<ProjectInfo>): Promise<ProjectInfo> => await QuietAjax.post('/projects', projectInfo)

export const updateProjectBaseInfo = () => useApi(async ({ id: project_id, ...info }: ProjectInfo) => Ajax.put(`/projects/${project_id}`, info))

export const getProjectServerUrlList = async (project_id: any) => Ajax.get(`/projects/${project_id}/servers`)

export const saveProjectServerUrlList = async ({ project_id, urls }: any) => Ajax.put(`/projects/${project_id}/servers`, urls)

export const exportProject = ({ project_id, ...params }: any) => `/api/projects/${project_id}/data?${new URLSearchParams(params).toString()}`

export const getProjectTranshList = () => useApi(async (project_id: string) => Ajax.get(`/projects/${project_id}/trashs`))

export const restoreDoc = () =>
  useApi(async ({ project_id, ids }: any) => Ajax.put(`/projects/${project_id}/trashs?${ids.map((id: any) => `collection-id=${id}`).join('&')}`, { category: 0 }))

export const deleleProject = () => useApi(async (project_id: string) => Ajax.delete(`/projects/${project_id}`))

// 获取非此项目的成员列表
export const getMembersWithoutProject = async (project_id: string) => QuietAjax.get(`/projects/${project_id}/members/without`)
// 获取成员列表
export const getMembersInProject = (project_id: string) => async (data: Record<string, any>) => Ajax.get(`/projects/${project_id}/members?${new URLSearchParams(data).toString()}`)
// 新增成员
export const addMemberToProject = (project_id: string) => async (data: Record<string, any>) => Ajax.post(`/projects/${project_id}/members`, data)
// 删除成员
export const deleteMemberFromProject = async (project_id: string, user_id: number) => Ajax.delete(`/projects/${project_id}/members/${user_id}`)
// 修改成员权限
export const updateMemberAuthorityInProject = async (project_id: string, user_id: number, authority: MemberAuthorityInProject) =>
  QuietAjax.put(`/projects/${project_id}/members/authority/${user_id}`, { authority })
// 退出项目
export const quitProject = async (project_id: string) => Ajax.delete(`/projects/${project_id}/exit`)
// 移交项目
export const transferProject = async (project_id: string, member_id: number) => Ajax.put(`/projects/${project_id}/transfer`, { member_id })
