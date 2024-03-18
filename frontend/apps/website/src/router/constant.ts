export const ROOT_PATH = '/'
export const ROOT_PATH_NAME = 'root'

export const MAIN_PATH = '/main'
export const MAIN_PATH_ALIAS = '/home'
export const MAIN_PATH_NAME = 'main'

export const NO_PERMISSION_PATH = '/403'
export const NOT_FOUND_PATH = '/404'

// sign
export const LOGIN_PATH = '/login'
export const LOGIN_NAME = 'login'
export const REGISTER_PATH = '/register'
export const REGISTER_NAME = 'register'
export const FORGETPASS_PATH = '/forget_pass'
export const FORGETPASS_NAME = 'forget_pass'

// send email
export const RESET_PASS_PATH = '/reset_password/:token'
export const RESET_PASS_NAME = 'reset.password'
export const REGISTER_VERIFICATION_EMAIL_PATH = '/email_verification/:token'
export const REGISTER_VERIFICATION_EMAIL_NAME = 'register.email.verification'
export const USER_CHANGE_ACCOUNT_EMAIL_PATH = '/change_email/:token'
export const USER_CHANGE_ACCOUNT_EMAIL_NAME = 'user.change.email'

export const PROJECT_SHARE_VALIDATION_NAME = 'share.proejct.verification'
export const PROJECT_SHARE_VALIDATION_PATH = '/projects/:project_id/verification'

export const DOCUMENT_SHARE_PATH = '/share/:doc_public_id'

// user
export const USER_SETTING_PATH = '/user'
export const USER_SETTING_NAME = 'user'
export const USER_PAGE_PATH = '/user/:page?'
export const USER_PAGE_NAME = 'user.page'

// oauth
export const COMPLETE_INFO_PATH = '/complete_info/:type'
export const COMPLETE_INFO_NAME = 'completeInfo'
export const OAUTH_NAME = 'oauth.login'
export const OAUTH_PATH = '/oauth/login/:type'

export const OAUTH_CONNECT_NAME = 'connect.oauth'
export const OAUTH_CONNECT_PATH = '/oauth/connect/:type'

// team
export const TEAM_PATH = '/team'
export const TEAM_NAME = 'team'
export const TEAM_CREATE_PATH = '/team/create'
export const TEAM_CREATE_NAME = 'team.create'
export const TEAM_JOIN_PATH = '/team/join/token/:token'
export const TEAM_JOIN_NAME = 'team.join'

// project
export const PROJECT_LIST_ROOT_PATH_NAME = 'projects'
export const PROJECT_LIST_ROOT_PATH = '/projects'

// project details
export const PROJECT_DETAIL_PATH_NAME = 'project.detail'
export const PROJECT_DETAIL_PATH = '/projects/:project_id'

// project collection
export const PROJECT_COLLECTION_PATH_NAME = 'project.detail.collection'
export const PROJECT_COLLECTION_PATH = '/projects/:project_id/collection/:collectionID?'

// project schema
export const PROJECT_SCHEMA_PATH_NAME = 'project.detail.schema'
export const PROJECT_SCHEMA_PATH = '/projects/:project_id/schema/:schemaID'

// project response
export const PROJECT_RESPONSE_PATH_NAME = 'project.detail.response'
export const PROJECT_RESPONSE_PATH = '/projects/:project_id/response/:responseID'

// collection share
export const COLLECTION_SHARE_PATH_NAME = 'project.detail.collection.share'
export const COLLECTION_SHARE_PATH = '/share/:collectionPublicID'

// iteration
const iterationName = 'iteration.detail'
const iterationPath = '/iteration/:iterationID'
export const ITERATION_LIST_ROOT_PATH_NAME = 'iterations'
export const ITERATION_LIST_ROOT_PATH = '/iterations'
export const ITERATION_DETAIL_PATH_NAME = iterationName
export const ITERATION_DETAIL_PATH = iterationPath
export const ITERATION_COLLECTION_PATH_NAME = `${iterationName}.collection`
export const ITERATION_COLLECTION_PATH = `${iterationPath}/collection/:collectionID`
export const ITERATION_SCHEMA_PATH_NAME = `${iterationName}.schema`
export const ITERATION_SCHEMA_PATH = `${iterationPath}/schema/:schemaID`
export const ITERATION_RESPONSE_PATH_NAME = `${iterationName}.response`
export const ITERATION_RESPONSE_PATH = `${iterationPath}/response/:responseID`

export const redirectParam = 'redirect'

// system
export const SYSTEM_SETTING_PATH = '/system'
export const SYSTEM_SETTING_NAME = 'system'
export const SYSTEM_PAGE_PATH = '/system/:page?'
export const SYSTEM_PAGE_NAME = 'system.page'
