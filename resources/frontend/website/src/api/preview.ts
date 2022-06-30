import Ajax from './Ajax'

export const getDocumentInfo = (doc_id: any) => Ajax.get('/preview/doc_info', { params: { doc_id } })

export const checkDocumentSecretKey = (data: any) => Ajax.post('/api_doc/secretkey_check', data)

export const getDocumentStatus = (doc_id: any) => Ajax.get('/api_doc/has_shared', { params: { doc_id } })

export const getSingleApiDocumentDetail = (doc_id: any, token: any) => Ajax.get('/api_doc', { params: { doc_id, token, format: 'html' } })

export const getTrashNormalDocumentDetail = (project_id: any, doc_id: any) =>
    Ajax.get('/api_doc', { params: { project_id, doc_id, deleted: 1, format: 'html' } })

export const getProjectInfo = (project_id: any) => Ajax.get('/preview/project', { params: { project_id } })

export const getProjectCatalog = (token: any, project_id: any) => Ajax.get('/preview/api_nodes', { params: { project_id, token } })

export const checkProjectSecretKey = (data: any) => Ajax.post('/project/secretkey_check', data)

export const getApiDocumentDetail = (token: any, project_id: any, node_id: any) => Ajax.get('/preview/api_doc', { params: { project_id, node_id, token } })

export const searchProjectInfo = (token: any, project_id: any, keywords: any) => Ajax.get('/preview/search', { params: { project_id, keywords, token } })
