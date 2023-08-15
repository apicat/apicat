import Ajax, { QuietAjax } from './Ajax'
import useApi from '@/hooks/useApi'
import { MemberAuthorityInProject } from '@/typings/member'
import { ProjectCover, ProjectInfo } from '@/typings/project'
import { API_URL, ProjectListCoverBgColors, ProjectListCoverIcons } from '@/commons/constant'
import { queryStringify, randomArray } from '@/commons'
import { setShareTokenToParams } from '@/store/share'

export const getProjectDefaultCover = (overwrite?: Partial<ProjectCover>): ProjectCover => ({
  type: 'icon',
  coverBgColor: randomArray(ProjectListCoverBgColors),
  coverIcon: randomArray(ProjectListCoverIcons),
  ...overwrite,
})

export const convertProjectCover = (project: ProjectInfo): ProjectInfo => {
  try {
    project.cover = JSON.parse(project.cover as string)
  } catch (error) {
    project.cover = getProjectDefaultCover({ coverBgColor: ProjectListCoverBgColors[1], coverIcon: ProjectListCoverIcons[0], type: 'icon' })
  }
  return project
}

export const getProjectList = (): Promise<ProjectInfo[]> => Ajax.get('/projects')

export const getProjectDetail = async (project_id: string, params?: Record<string, any>): Promise<ProjectInfo> => Ajax.get(`/projects/${project_id}${queryStringify(params)}`)

export const createProject = async (projectInfo: Partial<ProjectInfo>): Promise<ProjectInfo> => await QuietAjax.post('/projects', projectInfo)

export const updateProjectBaseInfo = async ({ id: project_id, ...info }: ProjectInfo): Promise<ProjectInfo> => Ajax.put(`/projects/${project_id}`, info)

export const getProjectServerUrlList = async (project_id: string, params?: Record<string, any>) => {
  params = setShareTokenToParams(params || {})
  return Ajax.get(`/projects/${project_id}/servers${queryStringify(params)}`)
}

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
// 获取已关注的项目列表
export const getFollowedProjectList = (): Promise<ProjectInfo[]> => Ajax.get('/projects/follow')
// 关注项目
export const followProject = (project_id: string) => QuietAjax.post(`/projects/${project_id}/follow`)
// 取消关注项目
export const unFollowProject = (project_id: string) => QuietAjax.delete(`/projects/${project_id}/follow`)
// toggle 关注项目
export const toggleFollowProject = (project: ProjectInfo) => (project.is_followed ? unFollowProject(project.id as string) : followProject(project.id as string))
