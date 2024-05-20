import { dayjs } from 'element-plus'
import { parseJSONWithDefault } from '@apicat/shared'
import type { JSONSchema } from '@apicat/editor'
import DefaultAjax from '@/api/Ajax'
import { gatherSharedTokenWithParams } from '@/api/shareToken'
import { SchemaTypeEnum } from '@/commons'

// FIXME: 全局化 数据默认结构
export function getDefaultSchemaStructure() {
  return {
    'type': 'object',
    'properties': {},
    'required': [],
    'x-apicat-orders': [],
    'example': '',
  }
}

export function apiCreateSchema(projectID: string, data: Omit<Definition.Schema, 'id'>): Promise<Definition.Schema> {
  return DefaultAjax.post(`/projects/${projectID}/definition/schemas`, data)
}
export function apiAICreateSchema(projectID: string, data: Definition.RequestCreateSchemaWithAI): Promise<Definition.Schema> {
  return DefaultAjax.post(`/projects/${projectID}/definition/ai/schemas`, data)
}

export async function apiGetSchemaInfo(projectID: string, schemaID: number): Promise<Definition.Schema> {
  const res: Definition.Schema = await DefaultAjax.get(`/projects/${projectID}/definition/schemas/${schemaID}`, { params: gatherSharedTokenWithParams({}, projectID) })
  res.schema = parseJSONWithDefault(res.schema, getDefaultSchemaStructure())
  return res
}

export async function apiGetSchemaTree(projectID: string): Promise<Definition.SchemaNode[]> {
  try {
    const res: Definition.SchemaNode[] = await DefaultAjax.get(
      `/projects/${projectID}/definition/schemas`,
      { params: gatherSharedTokenWithParams({}, projectID) },
      { isShowErrorMsg: false },
    )
    function mapper(item: Definition.SchemaNode) {
      if (item.items && item.items.length)
        item.items.map(mapper)

      if (item.type === SchemaTypeEnum.Schema)
        item.schema = parseJSONWithDefault(item.schema, getDefaultSchemaStructure())

      return item
    }
    return res.map(mapper)
  }
  catch (error) {
    return []
  }
}

export function apiEditSchema(projectID: string, { id, ...schema }: Definition.Schema): Promise<void> {
  return DefaultAjax.put(`/projects/${projectID}/definition/schemas/${id}`, {
    ...schema,
    schema: JSON.stringify(schema.schema || {}),
  })
}

export function apiDeleteSchema(projectID: string, schemaID: number, deref: boolean): Promise<void> {
  return DefaultAjax.delete(`/projects/${projectID}/definition/schemas/${schemaID}?deref=${deref}`, null, {
    isShowErrorMsg: true,
    isShowSuccessMsg: true,
  })
}

export function apiMoveSchema(projectID: string, data: Definition.RequestSortParams): Promise<void> {
  return DefaultAjax.put(`/projects/${projectID}/definition/schemas/move`, data)
}

export async function apiCopySchema(projectID: string, schemaID: number): Promise<Definition.SchemaNode> {
  const res = await DefaultAjax.post(`/projects/${projectID}/definition/schemas/${schemaID}/copy`)
  res.schema = parseJSONWithDefault(res.schema, getDefaultSchemaStructure())
  return res
}

export function apiGetSchemaHistories(projectID: string, schemaID: number): Promise<HistoryRecord.SchemaHistory[]> {
  return DefaultAjax.get(`/projects/${projectID}/definition/schemas/${schemaID}/histories`, {
    params: { startTime: dayjs().subtract(3, 'month').unix(), endTime: dayjs().unix() },
  })
}

export async function apiGetSchemaHistoryInfo(projectID: string, schemaID: number, historyID: number): Promise<HistoryRecord.SchemaHistoryInfo> {
  const res = await DefaultAjax.get(`/projects/${projectID}/definition/schemas/${schemaID}/histories/${historyID}`)
  res.schema = parseJSONWithDefault(res.schema, {})
  return res
}

export function apiRestoreSchemaHistory(projectID: string, schemaID: number, historyID: number): Promise<void> {
  return DefaultAjax.put(`/projects/${projectID}/definition/schemas/${schemaID}/histories/${historyID}/restore`)
}

export async function apiDiffSchemaHistory(projectID: string, schemaID: number, originalID: number, targetID: number): Promise<HistoryRecord.SchemaHistoryDiff> {
  const res: HistoryRecord.SchemaHistoryDiff = await DefaultAjax.get(`/projects/${projectID}/definition/schemas/${schemaID}/histories/diff`, {
    params: {
      originalID,
      targetID,
    },
  })

  res.schema1.schema = parseJSONWithDefault(res.schema1.schema, {})
  res.schema2.schema = parseJSONWithDefault(res.schema2.schema, {})

  return res
}

// parse jsonschema
export async function apiParseSchema(jsonschema: JSONSchema): Promise<JSONSchema> {
  const { jsonschema: str } = await DefaultAjax.post<{ jsonschema: string }>('/jsonschema/parse', { jsonschema: typeof jsonschema === 'string' ? jsonschema : JSON.stringify(jsonschema) })
  if (str)
    return JSON.parse(str)
  return getDefaultSchemaStructure()
}
