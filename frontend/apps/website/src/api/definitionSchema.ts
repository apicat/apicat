import { queryStringify } from '@/commons'
import Ajax, { QuietAjax } from './Ajax'
import useApi from '@/hooks/useApi'
import { setShareTokenToParams } from '@/store/share'

const restfulApiPath = (project_id: string | number): string => `/projects/${project_id}/definition/schemas`
const detailRestfulPath = (project_id: string | number, def_id: string | number): string => `${restfulApiPath(project_id)}/${def_id}`

export const getDefinitionSchemaList = (project_id: string, params?: Record<string, any>) => {
  params = setShareTokenToParams(params || {})
  return Ajax.get(restfulApiPath(project_id) + queryStringify(params))
}

export const getDefinitionSchemaDetail = () =>
  useApi(async ({ project_id, def_id, ...params }: any) => {
    params = setShareTokenToParams(params || {})
    return Ajax.get(`${restfulApiPath(project_id)}/${def_id}${queryStringify(params)}`)
  })

export const createDefinitionSchema = async ({ project_id, ...definitionInfo }: any) => QuietAjax.post(restfulApiPath(project_id), definitionInfo)

export const updateDefinitionSchema = async ({ project_id, def_id, ...definitionInfo }: any) => QuietAjax.put(`${restfulApiPath(project_id)}/${def_id}`, definitionInfo)

export const copyDefinitionSchema = async (project_id: string, def_id: string | number) => Ajax.post(`${restfulApiPath(project_id)}/${def_id}`)

export const moveDefinitionSchema = async (project_id: string, sortParams: { target: any; origin: any }) => QuietAjax.put(`${restfulApiPath(project_id)}/movement`, sortParams)

export const deleteDefinitionSchema = async (project_id: string | number, def_id: string | number, is_unref: number) =>
  Ajax.delete(`${restfulApiPath(project_id)}/${def_id}?is_unref=${is_unref}`)

export const aiGenerateDefinitionSchema = async ({ project_id, ...params }: any) => Ajax.post(`/projects/${project_id}/ai/schemas`, params)

// 文档历史记录列表
export const getSchemaHistoryRecordList = ({ project_id, def_id }: Record<string, any>) => Ajax.get(`${detailRestfulPath(project_id, def_id)}/histories`)

// 文档历史记录详情
export const getSchemaHistoryRecordDetail = ({ project_id, def_id, history_id }: Record<string, any>) =>
  Ajax.get(`${detailRestfulPath(project_id, def_id)}/histories/${history_id}`)

// 文档历史记录对比
export const compareSchema = ({ project_id, def_id, ...params }: Record<string, any>) =>
  Ajax.get(`${detailRestfulPath(project_id, def_id)}/histories/diff${queryStringify(params)}`)

// 恢复文档
export const restoreSchemaByHistoryRecord = ({ project_id, def_id, history_id }: Record<string, any>) =>
  QuietAjax.put(`${detailRestfulPath(project_id, def_id)}/histories/${history_id}/restore`)
