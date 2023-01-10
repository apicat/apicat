import Ajax from './Ajax'
import { compile } from 'path-to-regexp'
import { ITERATE_DOCUMENT_DETAIL_PATH, ITERATE_DOCUMENT_EDIT_PATH, ITERATE_ROUTE_PATH } from '@/router/constant'

/**
 * 迭代列表
 * @param data
 * project_id
 * page
 */
export const getIterations = (data: any) => Ajax.get('/iterations', { params: data })

/**
 * 迭代详情
 * @param data
 * iteration_id
 */
export const getIterationDetail = (data: any) => Ajax.get('/iteration', { params: data })

/**
 * 新建迭代
 * @param iteration
 * project_id
 * title
 * description
 */
export const createIteration = (iteration: any) => Ajax.post('/iteration/create', iteration)

/**
 * 编辑迭代
 * @param iteration
 * iteration_id
 * title
 * description
 */
export const editIteration = (iteration: any) => Ajax.post('/iteration/edit', iteration)

/**
 * 删除迭代
 * @param data
 * iteration_id
 */
export const deleteIteration = (data: any) => Ajax.post('/iteration/remove', data)

/**
 * 收藏项目列表
 */
export const getCollectionList = () => Ajax.get('/iteration/stars')

/**
 * 收藏排序
 * @param data
 * project_ids
 */
export const sortCollectionList = (data: any) => Ajax.post('/iteration/star_order', data)

/**
 * 添加收藏
 * @param data
 * project_id
 */
export const addCollection = (data: any) => Ajax.post('/iteration/star', data)

/**
 * 取消收藏
 * @param data
 * project_id
 */
export const cancelCollection = (data: any) => Ajax.post('/iteration/unstar', data)

/**
 * 规划API到迭代
 * @param data
 * iteration_id
 * node_ids
 */
export const planApisToIterate = (data: any) => Ajax.post('/iteration/push', data)

export const toIterateDocumentPath = (iterate_id_public: string) => compile(ITERATE_ROUTE_PATH)({ iterate_id_public })
export const toIterateDocumentDetailPath = (iterate_id_public: string, node_id?: string) =>
    compile(ITERATE_DOCUMENT_DETAIL_PATH)({ iterate_id_public, node_id })
export const toIterateDocumentEditPath = compile(ITERATE_DOCUMENT_EDIT_PATH)
