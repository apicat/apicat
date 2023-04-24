import Ajax, { QuietAjax } from './Ajax'
import { convertRequestPath } from '@/commons'

const globalParamterListApiPath = (project_id: string | number): string => convertRequestPath('/projects/:project_id/global/parameters', { project_id })
const globalParamterApiPath = (project_id: string | number, parameter_id: string | number): string =>
  convertRequestPath('/projects/:project_id/global/parameters/:parameter_id', { project_id, parameter_id })

export const getGlobalParamList = async ({ project_id }: any) => await Ajax.get(globalParamterListApiPath(project_id))
export const createGlobalParamerter = async ({ project_id, ...param }: any) => await QuietAjax.post(globalParamterListApiPath(project_id), param)
export const updateGlobalParamerter = async ({ project_id, id, ...param }: any) => await QuietAjax.put(globalParamterApiPath(project_id, id), param)
export const deleteGlobalParamerter = async ({ project_id, id, ...param }: any) => await Ajax.delete(globalParamterApiPath(project_id, id), param)
