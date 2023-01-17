export const INDEX_PATH = '/'
export const ROOT_PATH = INDEX_PATH
export const LOGIN_PATH = '/login'
export const REGISTE_PATH = '/register'
export const MAIN_PATH = '/home'

export const NOT_FOUND = { path: '/404' }

// 项目相关链接地址
export const PROJECT_PREVIEW_PATH = '/project/:project_id'
export const PROJECT_DETAIL_PATH = PROJECT_PREVIEW_PATH
export const PROJECT_SETTING_PATH = '/project/:project_id/setting'
export const PROJECT_MEMBERS_PATH = '/project/:project_id/members'
export const PROJECT_PARAMS_PATH = '/project/:project_id/params'
export const PROJECT_TRASH_PATH = '/project/:project_id/trash'

// 预览相关链接地址
export const PREVIEW_PROJECT = 'preview.project'
export const PREVIEW_PROJECT_SECRET = `${PREVIEW_PROJECT}.verification`
export const PREVIEW_DOCUMENT = 'preview.document'
export const PREVIEW_DOCUMENT_PATH = '/doc/:doc_id'
export const PREVIEW_DOCUMENT_SECRET = `${PREVIEW_DOCUMENT}.verification`

// 文档相关链接地址
export const DOCUMENT_ROUTE_NAME = 'document'
export const DOCUMENT_EDIT_NAME = 'document.api.edit'
export const DOCUMENT_EDIT_PATH = '/project/:project_id/doc/:node_id/edit'
export const DOCUMENT_DETAIL_NAME = 'document.api.detail'
export const DOCUMENT_DETAIL_PATH = '/project/:project_id/doc/:node_id?'

// 迭代相关链接地址
export const ITERATE_ROUTE_PATH = '/iteration/:iterate_id'
export const ITERATE_ROUTE_NAME = 'iteration.document'
export const ITERATE_DOCUMENT_EDIT_PATH = '/iteration/:iterate_id/doc/:node_id/edit'
export const ITERATE_DOCUMENT_EDIT_NAME = 'iteration.document.api.edit'
export const ITERATE_DOCUMENT_DETAIL_NAME = 'iteration.document.api.detail'
export const ITERATE_DOCUMENT_DETAIL_PATH = '/iteration/:iterate_id/doc/:node_id?'

// 历史文档相关链接地址
export const DOCUMENT_HISTORY_PATH = '/history/:project_id_public'
export const DOCUMENT_HISTORY_NAME = 'history.doc'
export const DOCUMENT_HISTORY_DETAIL_PATH = `${DOCUMENT_HISTORY_PATH}/:doc_id/:id?`
export const DOCUMENT_HISTORY_DETAIL_NAME = 'history.doc.detail'
