import { convertRequestPath, isJSONSchemaContentType, queryStringify } from '@/commons'
import Ajax, { QuietAjax } from './Ajax'
import useApi from '@/hooks/useApi'
import { setShareTokenToParams } from '@/store/share'

const restfulApiPath = (project_id: string | number): string => convertRequestPath('/projects/:project_id/definition/responses', { project_id })

export const getDefinitionResponseList = (project_id: string, params?: Record<string, any>) => {
  params = setShareTokenToParams(params || {})
  return Ajax.get(restfulApiPath(project_id) + queryStringify(params))
}

export const getDefinitionResponseDetail = () =>
  useApi(async ({ project_id, id, ...params }: any) => {
    params = setShareTokenToParams(params || {})
    const data: any = await Ajax.get(`${restfulApiPath(project_id)}/${id}${queryStringify(params)}`)
    const contentType: string = Object.keys(data.content || {})[0] || 'application/json'
    // 补充默认结构
    if (isJSONSchemaContentType(contentType) && data.content[contentType].schema && !data.content[contentType].schema.properties) {
      data.content[contentType].schema.properties = {}
    }
    return data
  })

export const createDefinitionResponse = async ({ project_id, ...definitionInfo }: any) => QuietAjax.post(restfulApiPath(project_id), definitionInfo)

export const updateDefinitionResponse = async ({ project_id, id, ...definitionInfo }: any) => QuietAjax.put(`${restfulApiPath(project_id)}/${id}`, definitionInfo)

export const deleteDefinitionResponse = async (project_id: string | number, id: string | number, is_unref: number) =>
  Ajax.delete(`${restfulApiPath(project_id)}/${id}?is_unref=${is_unref}`)
