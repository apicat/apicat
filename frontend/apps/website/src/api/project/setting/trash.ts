import Ajax from '../../Ajax'

export function apiGetTrashList(projectID: string): Promise<ProjectAPI.Trash[]> {
  return Ajax.get(`/projects/${projectID}/collections/trashes`)
}

export function apiRestoreTrash(projectID: string, data: ProjectAPI.RequestRestoreTrash): Promise<ProjectAPI.Trash> {
  return Ajax.put(`/projects/${projectID}/collections/restore`, data)
}
