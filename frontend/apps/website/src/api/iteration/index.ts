import DefaultAjax from '../Ajax'

export function apiCreateIteration(
  teamID: string,
  data: IterationAPI.RequestCreateIteration,
): Promise<IterationAPI.ResponseIteration> {
  return DefaultAjax.post(`/teams/${teamID}/iterations`, data)
}

export function apiGetIterations(
  params: Partial<GlobalAPI.RequestTable>,
  teamID: string,
): Promise<GlobalAPI.ResponseTable<IterationAPI.ResponseIteration[]>> {
  return DefaultAjax.get(`/teams/${teamID}/iterations`, { params })
}

export function apiGetIterationInfo(iterationID: string): Promise<IterationAPI.ResponseIteration> {
  return DefaultAjax.get(`/iterations/${iterationID}`)
}

export function apiEditIterationInfo(iterationID: string, data: IterationAPI.RequestEditIteration): Promise<void> {
  return DefaultAjax.put(`/iterations/${iterationID}`, data)
}

export function apiDeleteIterationInfo(iterationID: string): Promise<void> {
  return DefaultAjax.delete(`/iterations/${iterationID}`)
}
