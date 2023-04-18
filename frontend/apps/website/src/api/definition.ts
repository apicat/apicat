import Ajax, { QuietAjax } from './Ajax'
import useApi from '@/hooks/useApi'

export const getDefinitionList = (project_id: string) => Ajax.get(`/projects/${project_id}/definitions`)

export const getDefinitionDetail = useApi(async ({ project_id, def_id }: any) => Ajax.get(`/projects/${project_id}/definitions/${def_id}`))

export const createDefinition = async ({ project_id, ...definitionInfo }: any) => QuietAjax.post(`/projects/${project_id}/definitions`, definitionInfo)

export const updateDefinition = async ({ project_id, def_id, ...definitionInfo }: any) => QuietAjax.put(`/projects/${project_id}/definitions/${def_id}`, definitionInfo)

export const copyDefinition = async (project_id: string, def_id: string | number) => Ajax.post(`/projects/${project_id}/definitions/${def_id}`)

export const moveDefinition = async (project_id: string, sortParams: { target: any; origin: any }) => QuietAjax.put(`/projects/${project_id}/definitions/movement`, sortParams)

export const deleteDefinition = async (project_id: string | number, def_id: string | number) => Ajax.delete(`/projects/${project_id}/definitions/${def_id}`)
