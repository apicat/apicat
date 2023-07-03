import { Menu } from '@/components/typings'
import { Language } from '@/typings/common'
import { useI18n } from 'vue-i18n'

// API请求前缀
export const API_URL = '/api'

// 权限变更状态码
export const enum PERMISSION_CHANGE_CODE {
  // 用户权限变更
  USER_PREMISSION_ERROR = 101,
  // 项目中的成员权限变更
  MEMBER_PREMISSION_ERROR = 201,
  // 目标成员在项目中的权限发生变更
  TARGET_MEMBER_PREMISSION_ERROR = 202,
}

// 请求超时时长
export const REQUEST_TIMEOUT = 1000 * 60

// 默认值
export const DEFAULT_VAL = '--'

// DEFAULT LANGUAGE
export const DEFAULT_LANGUAGE = 'zh_CN'

// SUPPORTED LANGUAGES
export const SUPPORTED_LANGUAGES: Language[] = [
  { name: '中文', lang: 'zh_CN' },
  { name: 'English', lang: 'en' },
]

// 本地存储前缀
export const STORAGE_PREFIX = 'api.cat'

// 文档类型
export const enum DocumentTypeEnum {
  DIR = 'category',
  DOC = 'doc',
  HTTP = 'http',
}
// Collection visibility private-私有项目文档，public-公开项目文档
export const enum CollectionVisibilityEnum {
  PRIVATE = 'private',
  PUBLIC = 'public',
}

export const enum ProjectVisibilityEnum {
  PRIVATE = 'private',
  PUBLIC = 'public',
}

export const enum DefinitionTypeEnum {
  DIR = 'category',
  SCHEMA = 'schema',
  RESPONSE = 'response',
}

// 导出|导出状态
export const enum ImportOrExportState {
  WAIT = 'wait',
  FINISH = 'finish',
  FAIL = 'fail',
}

/**
 * 项目布局中导航菜单
 */
export enum ProjectNavigateListEnum {
  BaseInfoSetting = 'BaseInfoSetting',
  ProjectMemberList = 'ProjectMemberList',
  ServerUrlSetting = 'ServerUrlSetting',
  GlobalParamsSetting = 'GlobalParamsSetting',
  ResponseParamsSetting = 'ResponseParamsSetting',
  ProjectExport = 'ProjectExport',
  ProjectTrash = 'ProjectTrash',
  QuitProject = 'QuitProject',
}

export type ProjectNavigateObject = { [key in ProjectNavigateListEnum]: Menu }

/**
 *
 * @returns use function in setup
 * { [key in ProjectNavigateListEnum]: { [key: string]: any }
 */
export const getProjectNavigateList = (overwrite?: any): ProjectNavigateObject => {
  const { t } = useI18n()

  const navs = {
    [ProjectNavigateListEnum.BaseInfoSetting]: { text: t('app.project.setting.baseInfo'), icon: 'ac-setting' },
    [ProjectNavigateListEnum.ProjectMemberList]: { text: t('app.project.setting.member'), icon: 'ac-members' },
    [ProjectNavigateListEnum.ServerUrlSetting]: { text: t('app.project.setting.serverUrl'), icon: 'ac-suffix-url' },
    [ProjectNavigateListEnum.GlobalParamsSetting]: { text: t('app.project.setting.globalParam'), icon: 'ac-canshuweihu' },
    [ProjectNavigateListEnum.ProjectExport]: { text: t('app.project.setting.export'), icon: 'ac-export' },
    [ProjectNavigateListEnum.ProjectTrash]: { text: t('app.project.setting.trash'), icon: 'ac-trash' },
  } as any

  if (overwrite) {
    Object.keys(navs).forEach((key: any) => {
      const item = navs[key]
      const extendItem = overwrite[key]
      navs[key] = { ...item, ...extendItem }
    })
  }

  return navs as ProjectNavigateObject
}

// 项目导出类型
export const enum ExportProjectTypes {
  Swagger = 'swagger',
  OpenAPI = 'openapi',
  HTML = 'HTML',
  MARKDOWN = 'md',
}

export const enum CommonParameterType {
  String = 'string',
  Integer = 'integer',
  Number = 'number',
  Array = 'array',
}

export const CommonParameterTypes = [
  { text: 'string', value: 'string' },
  { text: 'integer', value: 'integer' },
  { text: 'number', value: 'number' },
  { text: 'array', value: 'array' },
]

export const ContentTypes = [
  { key: 'json', value: 'application/json' },
  { key: 'xhtml', value: 'application/xhtml+xml' },
  { key: 'xml', value: 'application/xml' },
  { key: 'text', value: 'text/plain' },
  { key: 'stream', value: 'application/octet-stream' },
]

export const ContentTypesMap = {
  json: 'application/json',
  xhtml: 'application/xhtml+xml',
  xml: 'application/xml',
  text: 'text/plain',
  stream: 'application/octet-stream',
}

export const RequestContentTypesMap = {
  none: 'none',
  'form-data': 'multipart/form-data',
  'x-www-form-urlencoded': 'application/x-www-form-urlencoded',
  json: 'application/json',
  xml: 'application/xml',
  raw: 'raw',
  binary: 'application/octet-stream',
}

export const HttpMethodTypeMap = {
  get: { value: 'get', color: '#66BE74' },
  post: { value: 'post', color: '#4894FF' },
  put: { value: 'put', color: '#51B9C3' },
  patch: { value: 'patch', color: '#F1924E' },
  delete: { value: 'delete', color: '#DF4545' },
  option: { value: 'option', color: '#A973DF' },
}

export const RefPrefixKeys = {
  CommonResponse: {
    key: '#/commons/responses/',
    reg: /#\/commons\/responses\/(.*)/,
  },
  DefinitionResponse: {
    key: '#/definitions/responses/',
    reg: /#\/definitions\/responses\/(.*)/,
  },
  DefinitionSchema: {
    key: '#/definitions/schemas/',
    replaceForCodeGenerate: '#/definitions/schemas_',
    refForCodeGeneratePrefix: 'schemas_',
    reg: /#\/definitions\/schemas\/(.*)/,
  },
}

export const ProjectListCoverBgColors = ['#FF9966', '#6699CC', '#FFCC99', '#66CCCC', '#FFCCCC', '#CCCCFF', '#99CC99', '#CCCCCC', '#99CCFF']
export const ProjectListCoverIcons = ['ac-danjumoxing', 'ac-fangkuai', 'ac-home', 'ac-jiekou', 'ac-doc-text', 'ac-bijibendiannao', 'ac-diannao', 'ac-shiyan', 'ac-xiangmu']
