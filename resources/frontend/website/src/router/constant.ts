export const INDEX_PATH = '/'
export const LOGIN_PATH = '/login'
export const REGISTE_PATH = '/register'
export const MAIN_PATH = '/main'

export const NOT_FOUND = { path: '/404' }
export const DOCUMENT_ROUTE_NAME = 'document'

// 项目相关链接地址
export const PROJECT_PREVIEW_PATH = '/editor/:project_id'
export const PROJECT_SETTING_PATH = '/project/:project_id/setting'
export const PROJECT_MEMBERS_PATH = '/project/:project_id/members'
export const PROJECT_PARAMS_PATH = '/project/:project_id/params'
export const PROJECT_TRASH_PATH = '/project/:project_id/trash'

// 预览相关链接地址
export const PREVIEW_PROJECT = 'preview.project'
export const PREVIEW_PROJECT_SECRET = `${PREVIEW_PROJECT}.verification`
export const PREVIEW_DOCUMENT = 'preview.document'
export const PREVIEW_DOCUMENT_SECRET = `${PREVIEW_DOCUMENT}.verification`

// 文档相关链接地址
export const DOCUMENT_EDIT_NAME = 'document.api.edit'
export const DOCUMENT_DETAIL_NAME = 'document.api.detail'
export const DECUMENT_DETAIL_PATH = '/editor/:project_id/doc/:node_id?'
export const DECUMENT_EDIT_PATH = '/editor/:project_id/doc/:node_id/edit'
