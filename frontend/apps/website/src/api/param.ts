import { setShareTokenToParams } from '@/store/share'
import Ajax, { QuietAjax } from './Ajax'
import { queryStringify } from '@/commons'

const globalParamterListApiPath = (project_id: string | number): string => `/projects/${project_id}/global/parameters`
const globalParamterApiPath = (project_id: string | number, parameter_id: string | number): string => `/projects/${project_id}/global/parameters/${parameter_id}`

export const getGlobalParamList = async ({ project_id }: any, params?: Record<string, any>) => {
  params = setShareTokenToParams(params || {})
  return Ajax.get(globalParamterListApiPath(project_id) + queryStringify(params))
}

export const createGlobalParamerter = async ({ project_id, ...params }: any) => QuietAjax.post(globalParamterListApiPath(project_id), params)
export const updateGlobalParamerter = async ({ project_id, id, ...params }: any) => QuietAjax.put(globalParamterApiPath(project_id, id), params)
export const deleteGlobalParamerter = async ({ project_id, id, is_unref }: any) => Ajax.delete(globalParamterApiPath(project_id, id) + '?is_unref=' + is_unref)
