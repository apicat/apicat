import { parseJSONWithDefault } from '@apicat/shared'
import { getDefaultSchemaStructure } from './schema'
import DefaultAjax from '@/api/Ajax'
import { gatherSharedTokenWithParams } from '@/api/shareToken'
import { ResponseTypeEnum } from '@/commons'

// FIXME: 全局化 数据默认结构
export function getDefaultEmptyResponse() {
  return {
    'application/json': {
      schema: getDefaultSchemaStructure(),
    },
  }
}

export function apiCreateResponse(projectID: string, data: Definition.CreateResponseTreeNode): Promise<Definition.ResponseTreeNode> {
  return DefaultAjax.post(`/projects/${projectID}/definition/responses`, data)
}

export async function apiGetResponseInfo(projectID: string, responseID: number): Promise<Definition.ResponseDetail> {
  const res: Definition.ResponseDetail = await DefaultAjax.get(`/projects/${projectID}/definition/responses/${responseID}`, { params: gatherSharedTokenWithParams({}, projectID) })
  res.content = parseJSONWithDefault(res.content, getDefaultEmptyResponse())
  res.header = parseJSONWithDefault(res.header, [])
  return res
}

export async function apiGetResponseTree(projectID: string): Promise<Definition.ResponseTreeNode[]> {
  const res: Definition.ResponseTreeNode[] = await DefaultAjax.get(`/projects/${projectID}/definition/responses`, { params: gatherSharedTokenWithParams({}, projectID) }, { isShowErrorMsg: false })
  function mapper(item: Definition.ResponseTreeNode) {
    if (item.items && item.items.length)
      item.items.map(mapper)

    if (item.type === ResponseTypeEnum.Response) {
      item.content = parseJSONWithDefault(item.content, getDefaultEmptyResponse())
      item.header = parseJSONWithDefault(item.header, [])
    }

    return item
  }
  return res.map(mapper)
}

export function apiRenameResponseCategory(projectID: string, response: Definition.ResponseTreeNode): Promise<void> {
  return DefaultAjax.put(`/projects/${projectID}/definition/responses/${response.id}`, response)
}

export function apiEditResponse(projectID: string, { id, ...data }: Definition.UpdateResponse): Promise<void> {
  if (data.content)
    data.content = JSON.stringify(data.content || '') as any

  if (data.header)
    data.header = JSON.stringify(data.header || '[]') as any

  return DefaultAjax.put(`/projects/${projectID}/definition/responses/${id}`, data)
}

export function apiDeleteResponse(projectID: string, responseID: number, deref: boolean): Promise<void> {
  return DefaultAjax.delete(`/projects/${projectID}/definition/responses/${responseID}?deref=${deref}`, null, {
    isShowErrorMsg: true,
    isShowSuccessMsg: true,
  })
}

export function apiMoveResponse(projectID: string, data: Definition.RequestSortParams): Promise<void> {
  return DefaultAjax.put(`/projects/${projectID}/definition/responses/move`, data)
}

export async function apiCopyResponse(projectID: string, responseID: number): Promise<Definition.ResponseDetail> {
  const res = await DefaultAjax.post(`/projects/${projectID}/definition/responses/${responseID}/copy`)
  res.content = parseJSONWithDefault(res.content, {})
  res.header = parseJSONWithDefault(res.header, [])
  return res
}

// ai for response
export async function apiGetAIResponse(data: any) {
  return await DefaultAjax.post('/response/ai', data)
}
