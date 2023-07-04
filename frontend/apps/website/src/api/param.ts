import Ajax, { QuietAjax } from './Ajax'

const globalParamterListApiPath = (project_id: string | number): string => `/projects/${project_id}/global/parameters`
const globalParamterApiPath = (project_id: string | number, parameter_id: string | number): string => `/projects/${project_id}/global/parameters/${parameter_id}`

export const getGlobalParamList = async ({ project_id }: any) => await Ajax.get(globalParamterListApiPath(project_id))
export const createGlobalParamerter = async ({ project_id, ...params }: any) => await QuietAjax.post(globalParamterListApiPath(project_id), params)
export const updateGlobalParamerter = async ({ project_id, id, ...params }: any) => await QuietAjax.put(globalParamterApiPath(project_id, id), params)
export const deleteGlobalParamerter = async ({ project_id, id, is_unref }: any) => await Ajax.delete(globalParamterApiPath(project_id, id) + '?is_unref=' + is_unref)
