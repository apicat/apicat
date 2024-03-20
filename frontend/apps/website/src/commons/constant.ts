import type { Language } from '@/typings/common'

// API请求前缀
export const API_URL = '/api'

// 403 二级响应code字典 权限变更状态码
export enum PERMISSION_CHANGE_CODE {
  // 用户权限变更
  USER_PREMISSION_ERROR = 101,
  // 项目中的成员权限变更
  MEMBER_PREMISSION_ERROR = 201,
  // 目标成员在项目中的权限发生变更
  TARGET_MEMBER_PREMISSION_ERROR = 202,
  // 重定向403页面
  REDIRECT_UNAUTHORIZED_PAGE = 301,
}

// 404 二级响应code字典
export const RESPONSE_NOT_FOUND_MAPS = {
  REDIRECT_NOT_FOUND_PAGE: 101,
  SHOW_NOT_FOUND_MESSAGE: 201,
}

// 401 二级响应code字典
export const RESPONSE_UNAUTHORIZED_MAPS = {
  // 登录密钥失效或错误
  LOGIN_TOKEN_EXPIRED_OR_ERROR: 101,
  // 文档｜项目密钥失效或错误
  PROJECT_OR_DOCUMENT_SECRET_TOKEN_EXPIRED_OR_ERROR: 201,
}

// 请求超时时长
export const REQUEST_TIMEOUT = 1000 * 60

// 默认值
export const DEFAULT_VAL = '--'

// SUPPORTED LANGUAGES
export const SUPPORTED_LANGUAGES: Language[] = [
  { name: 'English', lang: 'en-US' },
  { name: '中文', lang: 'zh-CN' },
]

// DEFAULT LANGUAGE
export const DEFAULT_LANGUAGE = SUPPORTED_LANGUAGES[0].lang

// 本地存储前缀
export const STORAGE_PREFIX = 'api.cat'

// -------------------------------------------------------
// collection
export enum CollectionTypeEnum {
  Dir = 'category',
  Doc = 'doc',
  Http = 'http',
}
// schema
export enum SchemaTypeEnum {
  Category = 'category',
  Schema = 'schema',
}
// response
export enum ResponseTypeEnum {
  Category = 'category',
  Response = 'response',
}

// -------------------------------------------------------

// Collection visibility private-私有项目文档，public-公开项目文档
export enum CollectionVisibilityEnum {
  PRIVATE = 'private',
  PUBLIC = 'public',
}

export enum Visibility {
  Private = 'private',
  Public = 'public',
}

export enum DefinitionTypeEnum {
  DIR = 'category',
  SCHEMA = 'schema',
  RESPONSE = 'response',
}

// 导出|导出状态
export enum ImportOrExportState {
  WAIT = 'wait',
  FINISH = 'finish',
  FAIL = 'fail',
}

// 项目导出类型
export enum ExportProjectTypes {
  Swagger = 'swagger',
  OpenAPI = 'openapi',
  HTML = 'HTML',
  MARKDOWN = 'md',
  ApiCat = 'apicat',
}

// 项目导入类型
export enum ImportProjectTypes {
  ApiCat = 'apicat',
  OpenAPI = 'openapi',
  Swagger = 'swagger',
  Postman = 'postman',
}

export enum CommonParameterType {
  String = 'string',
  Integer = 'integer',
  Number = 'number',
  Array = 'array',
}

export enum Authority {
  Manage = 'manage',
  None = 'none',
  Read = 'read',
  Write = 'write',
}

export enum EditableRole {
  Admin = 'admin',
  Member = 'member',
}

export enum Role {
  // admin = 'Admin',
  // member = 'Member',
  // owner = 'Owner',
  Admin = 'admin',
  Member = 'member',
  Owner = 'owner',
}

export enum Status {
  Active = 'active',
  Deactive = 'deactive',
}

export enum OAuthPlatform {
  GITHUB = 'github',
}

export const OAuthPlatformConfig = {
  GITHUB: {
    OAUTH_URL: 'https://github.com/login/oauth/authorize',
    params: {
      scope: 'user:email',
    },
  },
}
/**
 * 成员在项目中的权限
 * 当前成员在此项目的权限:manage,write,read
 */
export enum MemberAuthorityInProject {
  MANAGER = 'manage',
  WRITE = 'write',
  READ = 'read',
  NONE = 'none',
}
export const MemberAuthorityMap = {
  [MemberAuthorityInProject.MANAGER]: 'app.memberAuth.manage',
  [MemberAuthorityInProject.WRITE]: 'app.memberAuth.write',
  [MemberAuthorityInProject.READ]: 'app.memberAuth.read',
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
  'none': 'none',
  'form-data': 'multipart/form-data',
  'x-www-form-urlencoded': 'application/x-www-form-urlencoded',
  'json': 'application/json',
  'xml': 'application/xml',
  'raw': 'raw',
  'binary': 'application/octet-stream',
}

export const ResponseContentTypesMap = {
  'application/json': 'json',
  'application/xml': 'xml',
  'text/html': 'html',
  'text/plain': 'raw',
  'application/octet-stream': 'binary',
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

export const ProjectListCoverBgColors = [
  '#FF9966',
  '#6699CC',
  '#FFCC99',
  '#66CCCC',
  '#FFCCCC',
  '#CCCCFF',
  '#99CC99',
  '#CCCCCC',
  '#99CCFF',
]
export const ProjectListCoverIcons = [
  'ac-danjumoxing',
  'ac-fangkuai',
  'ac-home',
  'ac-jiekou',
  'ac-doc-text',
  'ac-bijibendiannao',
  'ac-diannao',
  'ac-shiyan',
  'ac-xiangmu',
]

export enum SysStorage {
  Disk = 'disk',
  CF = 'cloudflare',
  Qiniu = 'qiniu',
}

export enum SysCache {
  Local = 'local',
  Redis = 'redis',
}

export enum SysEmail {
  SMTP = 'smtp',
  SendCloud = 'sendcloud',
}

export enum SysModel {
  Azure = 'azure-openai',
  OpenAI = 'openai',
}
