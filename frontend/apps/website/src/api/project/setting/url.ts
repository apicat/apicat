import Ajax from '../../Ajax'
import { gatherSharedTokenWithParams } from '@/api/shareToken'

export async function apiCreateProjectURL(
  projectID: string,
  data: ProjectAPI.RequestCreateURL,
): Promise<ProjectAPI.ResponseURL> {
  return Ajax.post(`/projects/${projectID}/servers`, data)
}

export async function apiGetProjectURLList(projectID: string): Promise<ProjectAPI.ResponseURL[]> {
  return Ajax.get(`/projects/${projectID}/servers`, { params: gatherSharedTokenWithParams({}, projectID) }, { isShowErrorMsg: false })
}

export async function apiEditProjectURL(
  projectID: string,
  serverID: number,
  data: ProjectAPI.RequestEditURL,
): Promise<void> {
  return Ajax.put(`/projects/${projectID}/servers/${serverID}`, data)
}

export async function apiDeleteProjectURL(projectID: string, serverID: number): Promise<void> {
  return Ajax.delete(`/projects/${projectID}/servers/${serverID}`)
}

export async function apiSortProjectURLList(projectID: string, data: ProjectAPI.RequestSortURL): Promise<void> {
  return Ajax.put(`/projects/${projectID}/servers/sort`, data)
}
