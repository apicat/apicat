import Ajax from '../../Ajax'

// 获取成员列表
export async function apiGetProjectMembers(
  params: Partial<GlobalAPI.RequestTable>,
  id: string,
): Promise<GlobalAPI.ResponseTable<ProjectAPI.Member[]>> {
  return Ajax.get(`/projects/${id}/members`, {
    params,
  })
}

// 获取不在此项目的成员
export async function apiGetExcludedMembers(id: string): Promise<TeamAPI.TeamMember[]> {
  return Ajax.get(`/projects/${id}/members/notin`)
}

// 创建成员
export async function apiCreateProjectMember({
  id,
  ...data
}: ProjectAPI.RequestCreateMember): Promise<TeamAPI.TeamMember[]> {
  return Ajax.post(`/projects/${id}/members`, data)
}

// 修改成员
export async function apiEditProjectMember(
  projectID: string,
  memberID: number,
  data: ProjectAPI.RequestEditMember,
): Promise<void> {
  return Ajax.put(`/projects/${projectID}/members/${memberID}`, data)
}

// 删除成员
export async function apiRemoveProjectMember(id: string, memberID: number): Promise<void> {
  return Ajax.delete(`/projects/${id}/members/${memberID}`)
}
