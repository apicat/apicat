import { ProjectGroup } from '@/typings'
import Ajax, { QuietAjax } from './Ajax'
// import { queryStringify } from '@/commons'

const restfulApiPath = (subPath?: string): string => `$/project_group${subPath ? `/${subPath}` : ''}`
const restfulOperationApiPath = (group_id: number): string => `${restfulApiPath()}/${group_id}`

// 获取分组列表
export const getProjectGroupList = async (): Promise<ProjectGroup[]> => Ajax.get(restfulApiPath())

// 新增分组
export const createProjectGroup = async (params: ProjectGroup): Promise<ProjectGroup> => QuietAjax.post(restfulApiPath(), params)

// 删除分组
export const deleteProjectGroup = async (group_id: number) => Ajax.delete(restfulOperationApiPath(group_id))

// 重命名分组
export const renameProjectGroup = async (group_id: number, params: ProjectGroup) => QuietAjax.put(`${restfulOperationApiPath(group_id)}/rename`, params)

// 分组排序
export const sortProjectGroup = async (ids: number[]) => QuietAjax.put(restfulApiPath('order'), { ids })
