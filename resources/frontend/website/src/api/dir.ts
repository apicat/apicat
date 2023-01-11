import Ajax from './Ajax'

/**
 * @param {*} data {project_id,name,parent_id}
 * @return node_id
 */
export const createDir = (data = {}) => Ajax.post('/dir/create', data)

/**
 * @param {*} data {project_id,node_id}
 */
export const deleteDir = (data = {}) => Ajax.post('/dir/remove', data)

/**
 * @param {*} data {project_id,node_id,name}
 */
export const renameDir = (data = {}) => Ajax.post('/dir/rename', data)

/**
 * @param {*} data {node_id,parent_node_id,position}
 */
export const sortTree = (data = {}) => Ajax.post('/api_tree/sort', data)

/**
 * 获取目录树
 * @param params
 * @param token
 */
export const treeList = (params: any, token?: string) => Ajax.get('/api_tree', { params: { ...(params || {}), token } })

export const getDirList = (project_id: string) => Ajax.get('/dir/list', { params: { project_id } })
