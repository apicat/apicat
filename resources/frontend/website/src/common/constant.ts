// API请求前缀
export const API_URL = '/api'

// 请求超时
export const REQUEST_TIMEOUT = 1000 * 10

// 分页大小
export const PAGE_SIZE = 15

// 状态码
export const enum HTTP_STATUS {
    OK = 0, //成功
    FAIL = -1, //失败
    NO_LOGIN = -2, //未登录
    NO_PARENT_DIR = -102, //恢复文档，没有父级目录
    NO_EXIST_TEAM = -100, //被移除，或者切换团队
    INVALID_PREVIEW_SECRET = -103, //访问秘钥失效
    NOT_FOUND = -404, //资源不存在
}

// 文档类型
export const DOCUMENT_TYPES = {
    DIR: 0,
    DOC: 1,
    HTTP: 1,
    OTHER: 2,
    MARKDOWN: 3,
}

// 文档类型
export const DOCUMENT_TYPE = {
    API: {
        text: '接口文档',
        key: 'api.doc',
        value: 1,
    },
    DB: {
        text: '数据库文档',
        key: 'db.doc',
        value: 2,
    },
}

// 导出导出状态
export const IMPORT_EXPORT_STATE = {
    WAIT: 'wait',
    FINISH: 'finish',
    FAIL: 'fail',
}

// 项目角色
export const PROJECT_ROLES_KEYS = {
    ADMIN: 'admin',
    MANAGER: 'manage',
    DEVELOPER: 'write',
    READER: 'read',
    NONE: 'none',
}

export const PROJECT_ROLES_MAP = {
    MANAGER: 0,
    DEVELOPER: 1,
    READER: 2,
}

export const PROJECT_ALL_ROLE_LIST = [
    { text: '团队管理员', value: 999, key: PROJECT_ROLES_KEYS.ADMIN },
    { text: '管理者', value: PROJECT_ROLES_MAP.MANAGER, key: PROJECT_ROLES_KEYS.MANAGER },
    { text: '维护者', value: PROJECT_ROLES_MAP.DEVELOPER, key: PROJECT_ROLES_KEYS.DEVELOPER },
    { text: '阅读者', value: PROJECT_ROLES_MAP.READER, key: PROJECT_ROLES_KEYS.READER },
]

export const PROJECT_ROLE_LIST = PROJECT_ALL_ROLE_LIST.filter((item: any) => [PROJECT_ROLES_KEYS.DEVELOPER, PROJECT_ROLES_KEYS.READER].indexOf(item.key) !== -1)

export const PROJECT_DEFAULT_ICON = '/static/icon-project.png'

// 项目可见性类型
export const PROJECT_VISIBLE_TYPES = {
    PRIVATE: 'private',
    PUBLIC: 'public',
}

export const PROJECT_VISIBLE_LIST = [
    { text: '私有', value: 0, key: PROJECT_VISIBLE_TYPES.PRIVATE },
    { text: '公开', value: 1, key: PROJECT_VISIBLE_TYPES.PUBLIC },
]

// 团队角色
export const TEAM_ROLE = {
    ADMIN: 0,
    MANAGER: 1,
    NORMAL: 2,
}

export const TEAM_ROLE_LIST = [
    { text: '管理员', value: TEAM_ROLE.MANAGER },
    { text: '普通成员', value: TEAM_ROLE.NORMAL },
]

export const TEAM_ALL_ROLE_LIST = [{ text: '超级管理员', value: TEAM_ROLE.ADMIN }, ...TEAM_ROLE_LIST]
