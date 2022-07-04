import Ajax from './Ajax'
import { wrapperOrigin } from '@/common/utils'
import { compile } from 'path-to-regexp'
import { Storage } from '@natosoft/shared'
import ApicatLogo from '@/assets/image/logo-apicat@2x.png'
import MarkdownLogo from '@/assets/image/logo-markdown@2x.png'
import PostmanLogo from '@/assets/image/logo-postman@2x.png'
import { DECUMENT_DETAIL_PATH } from '@/router/constant'

// 常用URL地址
export const getUrlTipList = (project_id: any) => Ajax.get('/api_url/list', { params: { project_id } })
export const deleteUrlTip = (project_id: any, url_id: any) => Ajax.post('/api_url/remove', { project_id, url_id })

// 文档详情综合
export const getDocumentDetail = (project_id: any, doc_id: any, format = 'json') => {
    const token = Storage.get(Storage.KEYS.SECRET_PROJECT_TOKEN + project_id || '', true)
    return Ajax.get('/api_doc', { params: { project_id, doc_id, format, token } })
}
export const searchDocuments = (project_id: any, keywords: any) => {
    const token = Storage.get(Storage.KEYS.SECRET_PROJECT_TOKEN + project_id || '', true)
    return Ajax.get('/api_doc/search', { params: { project_id, keywords, token } })
}
export const createDoc = (doc = {}) => Ajax.post('/api_doc/create', { ...doc })
export const createHttpDoc = (doc = {} as any) => Ajax.post('/api_doc/http_template', { ...doc })

export const updateDoc = (doc = {} as any) => {
    if (typeof doc.content === 'object') {
        doc.content = JSON.stringify(doc.content)
    }
    return Ajax.post('/api_doc/update', doc)
}

export const createDocByTemplate = (doc = {}) => Ajax.post('/api_doc/template', doc)
export const copyDoc = (doc = {}) => Ajax.post('/api_doc/copy', doc)
export const renameDoc = (doc = {}) => Ajax.post('/api_doc/rename', doc)
export const deleteDoc = (doc = {}) => Ajax.post('/api_doc/remove', doc)
export const shareDoc = (doc = {}) => Ajax.post('/api_doc/share', { ...doc })
export const shareDetailDoc = (doc = {}) => Ajax.post('/api_doc/share_detail', { ...doc })
export const resetDocShareSecretkey = (doc = {}) => Ajax.post('/api_doc/share_secretkey', { ...doc })
export const getTemplateList = (project_id: any) => Ajax.get('/project/templates', { params: { project_id } })
export const templateCheckName = (doc = {}) => Ajax.post('/api_doc/check_template_name', doc)
export const saveTemplate = (doc = {}) => Ajax.post('/project/add_template', doc)
export const importDocument = (data: any) => Ajax.post('/api_doc/import', data)
export const getImportDocumentResult = (project_id: any, job_id: any) => Ajax.get('/api_doc/import_result', { params: { project_id, job_id } })

// 可导入选项
export const API_DOCUMENT_IMPORT_ACTION_MAPPING = [
    { text: 'ApiCat', icon: ApicatLogo, type: 'apicat', action: importDocument, getJobResult: getImportDocumentResult, maxSize: 2, accept: '.json' },
    { text: 'Markdown', icon: MarkdownLogo, type: 'markdown', action: importDocument, getJobResult: getImportDocumentResult, maxSize: 0.5, accept: '.md' },
    { text: 'Postman(v2.1)', icon: PostmanLogo, type: 'postman', action: importDocument, getJobResult: getImportDocumentResult, maxSize: 2, accept: '.json' },
]

// 生成文档详情路由地址
export const generateDocumentDetailPath = (project_id: any, node_id: any, hasOrigin?: boolean) =>
    wrapperOrigin(hasOrigin) + compile(DECUMENT_DETAIL_PATH)({ project_id, node_id })
