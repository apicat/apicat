import { parseJSONWithDefault } from '@apicat/shared'
import Ajax, { QuietAjax } from '../Ajax'
import { gatherSharedTokenWithParams } from '../shareToken'
import useApi from '@/hooks/useApi'
import {
  API_URL,
  Authority,
  MemberAuthorityInProject,
  ProjectListCoverBgColors,
  ProjectListCoverIcons,
} from '@/commons/constant'
import { queryStringify, randomArray } from '@/commons'

export function createProject(
  teamID: string,
  data: ProjectAPI.RequestCreateProject,
): Promise<ProjectAPI.ResponseProject> {
  data.teamID = teamID
  // return QuietAjax.post(`/teams/${teamID}/projects`, data)
  return Ajax.post(`/teams/${teamID}/projects`, data)
}

// Project list ----------------------------------------------------
// COVER
export function getProjectDefaultCover(overwrite?: Partial<ProjectAPI.ProjectCover>): ProjectAPI.ProjectCover {
  return {
    type: 'icon',
    coverBgColor: randomArray(ProjectListCoverBgColors),
    coverIcon: randomArray(ProjectListCoverIcons),
    ...overwrite,
  }
}
export function convertProjectCover(project: ProjectAPI.ResponseProject): ProjectAPI.ResponseProject {
  project.cover = parseJSONWithDefault(
    project.cover as string,
    getProjectDefaultCover({
      coverBgColor: ProjectListCoverBgColors[1],
      coverIcon: ProjectListCoverIcons[0],
      type: 'icon',
    }),
  )
  return project
}

// LIST
export async function apiGetProjectList(
  teamID: string,
  params?: ProjectAPI.RequestProject,
  isShowErrorMsg = true,
): Promise<ProjectAPI.ResponseProject[]> {
  params = queryStringify(params) as any
  const d: ProjectAPI.ResponseProject[] = await Ajax.get(`/teams/${teamID}/projects${params}`, {}, { isShowErrorMsg })
  return (d || []).map((item: ProjectAPI.ResponseProject) => convertProjectCover(item))
}
export function apiGetMyProjectList(teamID: string): Promise<ProjectAPI.ResponseProject[]> {
  return apiGetProjectList(teamID, { permissions: Authority.Manage })
}
export function apiGetMyFollowedProjectList(teamID: string): Promise<ProjectAPI.ResponseProject[]> {
  return apiGetProjectList(teamID, { isFollowed: true })
}
export function apiGetProjectListByGroupId(
  teamID: string,
  groupID: number,
  showErrorMsg = true,
): Promise<ProjectAPI.ResponseProject[]> {
  return apiGetProjectList(teamID, { groupID }, showErrorMsg)
}
// Project list ----------------------------------------------------

export function apiGetProject(id: string, params: Record<string, any> = {}): Promise<ProjectAPI.ResponseProject> {
  return Ajax.get(`/projects/${id}`, { params: gatherSharedTokenWithParams(params, id) })
}

export function apiChangeProjectGroup(projectID: string, data: ProjectAPI.RequestChangeGroup): Promise<void> {
  return Ajax.put(`/projects/${projectID}/group`, data)
}

export function apiFollowProject(projectID: string): Promise<void> {
  return Ajax.post(`/projects/${projectID}/follow`)
}

export function apiUnfollowProject(projectID: string): Promise<void> {
  return Ajax.delete(`/projects/${projectID}/follow`)
}

export function apiSetProjectGeneral(projectID: string, data: ProjectAPI.RequestSetProjectGeneral): Promise<void> {
  return Ajax.put(`/projects/${projectID}`, data)
}

export function apiDeleteProject(id: string): Promise<void> {
  return Ajax.delete(`/projects/${id}`)
}

export function apiTransferProject(projectID: string, memberID: number): Promise<void> {
  return Ajax.put(`/projects/${projectID}/transfer`, { memberID })
}

export function apiQuitProject(projectID: string): Promise<void> {
  return Ajax.delete(`/projects/${projectID}/exit`)
}

export function apiExportProject(projectID: string, type: string, download: boolean): Promise<{ path: string }> {
  return Ajax.get(`/projects/${projectID}/export`, { params: { type, download } })
}

// ---------------------------------------------------------------------------

export async function getProjectList(params?: Record<string, any>): Promise<any[]> {
  let projects: any = await Ajax.get(`/projects${queryStringify(params)}`)
  projects = (projects || []).map((item: any) => convertProjectCover(item as any))
  return projects as any[]
}

export async function getMyProjectList(): Promise<any[]> {
  return await getProjectList({ auth: [MemberAuthorityInProject.MANAGER] })
}

export async function getMyFollowedProjectList(): Promise<any[]> {
  return await getProjectList({ is_followed: true })
}
export async function getProjectListByGroupId(group_id: number | null): Promise<any[]> {
  return await getProjectList({ group_id })
}

export async function getProjectDetail(project_id: string, params?: Record<string, any>): Promise<any> {
  return Ajax.get(`/projects/${project_id}${queryStringify(params)}`)
}

export async function updateProjectBaseInfo({ id: project_id, ...info }: any): Promise<any> {
  return Ajax.put(`/projects/${project_id}`, info)
}

export async function getProjectServerUrlList(project_id: string, params?: Record<string, any>) {
  return Ajax.get(`/projects/${project_id}/servers${queryStringify(params)}`)
}

export function saveProjectServerUrlList({ project_id, urls }: any) {
  return Ajax.put(`/projects/${project_id}/servers`, urls)
}

export function exportProject({ project_id, ...params }: any) {
  return `${API_URL}/projects/${project_id}/data${queryStringify(params)}`
}

export function getProjectTranshList() {
  return useApi((project_id: string) => Ajax.get(`/projects/${project_id}/trashs`))
}

export function restoreDoc() {
  return useApi(({ project_id, ids }: any) =>
    Ajax.put(`/projects/${project_id}/trashs?${ids.map((id: any) => `collection-id=${id}`).join('&')}`, {
      category: 0,
    }),
  )
}

export function deleleProject() {
  return useApi((project_id: string) => Ajax.delete(`/projects/${project_id}`))
}

// 获取非此项目的成员列表
export function getMembersWithoutProject(project_id: string) {
  return QuietAjax.get(`/projects/${project_id}/members/without`)
}

// 获取成员列表
export function getMembersInProject(project_id: string) {
  return (params: any): Promise<GlobalAPI.ResponseTable<any[]>> =>
    Ajax.get(`/projects/${project_id}/members${queryStringify(params)}`)
}
// 新增成员
export function addMemberToProject(project_id: string) {
  return (params: Record<string, any>) => Ajax.post(`/projects/${project_id}/members`, params)
}
// 删除成员
export function deleteMemberFromProject(project_id: string, user_id: number) {
  return Ajax.delete(`/projects/${project_id}/members/${user_id}`)
}
// 修改成员权限
export async function updateMemberAuthorityInProject(
  project_id: string,
  user_id: number,
  authority: MemberAuthorityInProject,
) {
  return QuietAjax.put(`/projects/${project_id}/members/authority/${user_id}`, {
    authority,
  })
}
// 退出项目
export function quitProject(project_id: string) {
  return Ajax.delete(`/projects/${project_id}/exit`)
}
// 移交项目
export function transferProject(project_id: string, member_id: number) {
  return Ajax.put(`/projects/${project_id}/transfer`, { member_id })
}
// 获取已关注的项目列表
export const getFollowedProjectList = getMyFollowedProjectList
// 关注项目
export function followProject(project_id: string) {
  return QuietAjax.post(`/projects/${project_id}/follow`)
}
// 取消关注项目
export function unFollowProject(project_id: string) {
  return QuietAjax.delete(`/projects/${project_id}/follow`)
}
// toggle 关注项目
export function toggleFollowProject(project: any) {
  return project.is_followed ? unFollowProject(project.id as string) : followProject(project.id as string)
}

// 项目分组设置
export function settingProjectGroup(project_id: string, target_group_id: number) {
  return Ajax.put(`/projects/${project_id}/change_group`, { target_group_id })
}
