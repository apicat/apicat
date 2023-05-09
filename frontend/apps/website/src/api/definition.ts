import { convertRequestPath } from '@/commons'
import Ajax, { QuietAjax } from './Ajax'
import useApi from '@/hooks/useApi'

const restfulApiPath = (project_id: string | number): string => convertRequestPath('/projects/:project_id/definition/schemas', { project_id })

export const getDefinitionList = (project_id: string) => Ajax.get(restfulApiPath(project_id))

export const getDefinitionDetail = () => useApi(async ({ project_id, def_id }: any) => Ajax.get(`${restfulApiPath(project_id)}/${def_id}`))

export const createDefinition = async ({ project_id, ...definitionInfo }: any) => QuietAjax.post(restfulApiPath(project_id), definitionInfo)

export const updateDefinition = async ({ project_id, def_id, ...definitionInfo }: any) => QuietAjax.put(`${restfulApiPath(project_id)}/${def_id}`, definitionInfo)

export const copyDefinition = async (project_id: string, def_id: string | number) => Ajax.post(`${restfulApiPath(project_id)}/${def_id}`)

export const moveDefinition = async (project_id: string, sortParams: { target: any; origin: any }) => QuietAjax.put(`${restfulApiPath(project_id)}/movement`, sortParams)

export const deleteDefinition = async (project_id: string | number, def_id: string | number, is_unref: number) =>
  Ajax.delete(`${restfulApiPath(project_id)}/${def_id}?is_unref=${is_unref}`)

export const aiGenerateDefinition = async ({ project_id, ...params }: any) => Ajax.post(`/projects/${project_id}/ai/schemas`, params)
