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
 * @param {*} project_id
 */
export const treeList = (project_id: string, token?: string) => Ajax.get('/api_tree', { params: { project_id, token } })

export const getDirList = (project_id: string) => Ajax.get('/dir/list', { params: { project_id } })
