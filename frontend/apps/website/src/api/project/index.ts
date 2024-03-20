import { parseJSONWithDefault } from '@apicat/shared'
import Ajax from '../Ajax'
import { gatherSharedTokenWithParams } from '../shareToken'
import {
  Authority,
  ProjectListCoverBgColors,
  ProjectListCoverIcons,
} from '@/commons/constant'
import { queryStringify, randomArray } from '@/commons'

export function createProject(
  teamID: string,
  data: ProjectAPI.RequestCreateProject,
): Promise<ProjectAPI.ResponseProject> {
  data.teamID = teamID
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
