import Ajax, { QuietAjax } from './Ajax'
import { queryStringify } from '@/commons'
import { Pageable, Iteration } from '@/typings'

const restfulRootPath = '/iterations'
const restfulDetailPath = (iteration_public_id: string | number): string => `${restfulRootPath}${iteration_public_id}`

// 获取迭代列表
export const getIterationList = async (params: { project_id?: string | number; page?: number; page_size?: number }): Promise<Pageable<{ iterations: Iteration[] }>> => {
  return QuietAjax.get(`${restfulRootPath}${queryStringify(params)}`)
}

// 获取迭代详情
export const getIterationDetail = async ({ iteration_public_id }: { iteration_public_id: string | number }): Promise<Iteration> => {
  return QuietAjax.get(`${restfulDetailPath(iteration_public_id)}`)
}

// 更新迭代
export const updateIteration = async (params: { iteration_public_id: string | number; title: string; desciption?: string; collection_ids: string[] }): Promise<Iteration> => {
  const { iteration_public_id, ...body } = params
  return Ajax.put(`${restfulDetailPath(iteration_public_id)}`, body)
}

// 创建迭代
export const createIteration = async (params: { project_id: string | number; title: string; desciption?: string; collection_ids: string[] }): Promise<Iteration> => {
  return Ajax.post(restfulRootPath, params)
}

// 删除迭代
export const deleteIteration = async ({ iteration_public_id }: { iteration_public_id: string | number }): Promise<void> => {
  return Ajax.delete(`${restfulDetailPath(iteration_public_id)}`)
}
