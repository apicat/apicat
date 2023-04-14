import Ajax from './Ajax'
import useApi from '@/hooks/useApi'
import { ProjectInfo } from '@/typings/project'

export const getProjectList = () => Ajax.get('/projects')

export const getProjectDetail = (project_id: string) => Ajax.get(`/projects/${project_id}`)

export const createProject = useApi(async (projectInfo: ProjectInfo) => await Ajax.post('/projects', projectInfo))

export const updateProjectBaseInfo = useApi(async ({ id: project_id, ...info }: ProjectInfo) => Ajax.put(`/projects/${project_id}`, info))

export const getProjectServerUrlList = async (project_id: any) => Ajax.get(`/projects/${project_id}/servers`)

export const saveProjectServerUrlList = async ({ project_id, urls }: any) => Ajax.put(`/projects/${project_id}/servers`, urls)

export const exportProject = useApi(async ({ project_id, ...params }: any) => Ajax.get(`/projects/${project_id}/data`, { params }))

export const getProjectTranshList = useApi(async (project_id: string) => Ajax.get(`/projects/${project_id}/trashs`))

export const restoreDoc = useApi(async ({ project_id, ids }: any) =>
  Ajax.put(`/projects/${project_id}/trashs?${ids.map((id: any) => `collection-id=${id}`).join('&')}`, { category: 0 })
)

export const deleleProject = useApi(async (project_id: string) => Ajax.delete(`/projects/${project_id}`))
