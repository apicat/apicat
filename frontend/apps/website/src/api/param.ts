import Ajax from './Ajax'
import useApi from '@/hooks/useApi'
import { convertRequestPath } from '@/commons'
import { CommonParam, ResponseList, ResponseParamDetail } from '@/typings'

const paramterApiPath = (project_id: string | number): string => convertRequestPath('/projects/:project_id/parameters', { project_id })
const responseParamApiPath = (project_id: string | number): string => convertRequestPath('/projects/:project_id/responses', { project_id })
const responseParamDetailApiPath = (project_id: string | number, response_id: string | number): string =>
  convertRequestPath('/projects/:project_id/responses/:response_id', { project_id, response_id })

export const getCommonParamList = useApi(async ({ project_id }: any) => await Ajax.get(paramterApiPath(project_id)))

export const saveCommonParamerter = async ({ project_id, ...params }: any) => await Ajax.put(paramterApiPath(project_id), params)
export const saveHeaderParamerter = useApi(async (params: CommonParam) => saveCommonParamerter({ in: 'header', ...params }))
export const saveCookieParamerter = useApi(async (params: CommonParam) => saveCommonParamerter({ in: 'cookie', ...params }))
export const saveQueryParamerter = useApi(async (params: CommonParam) => saveCommonParamerter({ in: 'query', ...params }))

export const getResponseParamList = useApi(async ({ project_id }: any): Promise<ResponseList[]> => await Ajax.get(responseParamApiPath(project_id)))

export const addResponseParam = async ({ project_id, ...param }: any): Promise<ResponseParamDetail> => await Ajax.post(responseParamApiPath(project_id), param)

export const updateResponseParam = async ({ project_id, response_id, ...param }: any): Promise<void> => await Ajax.put(responseParamDetailApiPath(project_id, response_id), param)

export const deleteResponseParam = async ({ project_id, response_id }: any) => await Ajax.delete(responseParamDetailApiPath(project_id, response_id))

export const getResponseParam = async ({ project_id, response_id }: any): Promise<ResponseParamDetail> => await Ajax.get(responseParamDetailApiPath(project_id, response_id))
