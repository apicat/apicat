import Ajax from './Ajax'

// 获取项目常用参数列表
export const getProjectCommonParamList = (project_id: any) => Ajax.get('/project/params', { params: { project_id } })
export const getApiParamList = getProjectCommonParamList
export const addApiParam = (param = {}) => Ajax.post('/project/add_param', { ...param })
export const deleteApiParam = (project_id: any, param_id: any) => Ajax.post('/project/remove_param', { param_id, project_id })
export const updateApiParam = (param = {}) => Ajax.post('/project/edit_param', param)
