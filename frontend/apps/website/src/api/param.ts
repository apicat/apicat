import Ajax, { QuietAjax } from './Ajax'
import useApi from '@/hooks/useApi'
import { convertRequestPath } from '@/commons'
import { ResponseList, ResponseParamDetail } from '@/typings'

const globalParamterListApiPath = (project_id: string | number): string => convertRequestPath('/projects/:project_id/global/parameters', { project_id })
const globalParamterApiPath = (project_id: string | number, parameter_id: string | number): string =>
  convertRequestPath('/projects/:project_id/global/parameters/:parameter_id', { project_id, parameter_id })

const responseParamApiPath = (project_id: string | number): string => convertRequestPath('/projects/:project_id/responses', { project_id })
const responseParamDetailApiPath = (project_id: string | number, response_id: string | number): string =>
  convertRequestPath('/projects/:project_id/responses/:response_id', { project_id, response_id })

export const getGlobalParamList = async ({ project_id }: any) => await Ajax.get(globalParamterListApiPath(project_id))
export const createGlobalParamerter = async ({ project_id, ...param }: any) => await QuietAjax.post(globalParamterListApiPath(project_id), param)
export const updateGlobalParamerter = async ({ project_id, id, ...param }: any) => await QuietAjax.put(globalParamterApiPath(project_id, id), param)
export const deleteGlobalParamerter = async ({ project_id, id, ...param }: any) => await Ajax.delete(globalParamterApiPath(project_id, id), param)

export const getResponseParamList = useApi(async ({ project_id }: any): Promise<ResponseList[]> => await Ajax.get(responseParamApiPath(project_id)))

export const addResponseParam = async ({ project_id, ...param }: any): Promise<ResponseParamDetail> => await Ajax.post(responseParamApiPath(project_id), param)

export const updateResponseParam = async ({ project_id, response_id, ...param }: any): Promise<void> => await Ajax.put(responseParamDetailApiPath(project_id, response_id), param)

export const deleteResponseParam = async ({ project_id, response_id }: any) => await Ajax.delete(responseParamDetailApiPath(project_id, response_id))

export const getResponseParam = async ({ project_id, response_id }: any): Promise<ResponseParamDetail> => await Ajax.get(responseParamDetailApiPath(project_id, response_id))
