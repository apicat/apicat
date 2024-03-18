import Ajax, { QuietAjax } from '../Ajax'

/**
 * create
 */
export async function createProjectGroup(
  teamID: string,
  data: ProjectAPI.RequestCreateGroup,
): Promise<ProjectAPI.ResponseGroup> {
  data.teamID = teamID
  return Ajax.post(`/teams/${teamID}/project-groups`, data)
}

/**
 * get list
 */
export async function getProjectGroupList(teamID: string): Promise<ProjectAPI.ResponseGroup[]> {
  return Ajax.get(`/teams/${teamID}/project-groups`)
}

/**
 * delete one
 */
export async function deleteProjectGroup(groupID: number): Promise<void> {
  return Ajax.delete(`/project-groups/${groupID}`)
}

/**
 * rename one
 */
export async function renameProjectGroup(groupID: string, data: ProjectAPI.RequestRenameGroup): Promise<void> {
  return Ajax.put(`/project-groups/${groupID}`, data)
}

/**
 * rename one
 */
export const updateProjectGroup = renameProjectGroup

export async function apiSortProjectGroup(teamID: string, data: ProjectAPI.RequestSortGroup): Promise<void> {
  return QuietAjax.put(`/teams/${teamID}/project-groups/sort`, data)
}
