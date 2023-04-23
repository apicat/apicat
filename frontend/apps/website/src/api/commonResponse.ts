import Ajax from './Ajax'
import { APICatCommonResponse } from '@/typings'
import { convertRequestPath } from '@/commons'

const responseParamApiPath = (project_id: string | number): string => convertRequestPath('/projects/:project_id/definitions/responses', { project_id })
const responseParamDetailApiPath = (project_id: string | number, response_id: string | number): string =>
  convertRequestPath('/projects/:project_id/definitions/responses/:response_id', { project_id, response_id })

export const getCommonResponseList = async ({ project_id }: any): Promise<APICatCommonResponse[]> => await Ajax.get(responseParamApiPath(project_id))

export const addResponseParam = async ({ project_id, ...param }: any): Promise<APICatCommonResponse> => await Ajax.post(responseParamApiPath(project_id), param)

export const updateResponseParam = async ({ project_id, response_id, ...param }: any): Promise<void> => await Ajax.put(responseParamDetailApiPath(project_id, response_id), param)

export const deleteResponseParam = async ({ project_id, response_id }: any) => await Ajax.delete(responseParamDetailApiPath(project_id, response_id))

export const getResponseParam = async ({ project_id, response_id }: any): Promise<APICatCommonResponse> => await Ajax.get(responseParamDetailApiPath(project_id, response_id))
