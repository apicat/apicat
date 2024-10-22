import type { JSONSchema } from '@apicat/editor'
import { guid } from '@apicat/shared'
import type { SchemaTreeNode } from '@apicat/components'
import { apiCheckReplaceModel, apiGetAISchema } from '@/api/project/definition/schema'

// AI generated json schema
export function useIntelligentSchema(projectID: string, getParams?: () => any) {
  let requestID
  let abortController: AbortController | null = null

  async function handleIntelligentSchema(jsonschema: JSONSchema, node: SchemaTreeNode): Promise<{ nid: string, schema: JSONSchema } | void> {
    // 避免重复请求
    abortController?.abort()

    requestID = `${guid()}:${node.id}`
    abortController = new AbortController()
    try {
      const { requestID: resID, schema } = await apiGetAISchema(projectID, {
        schema: JSON.stringify(jsonschema),
        requestID,
        ...getParams?.(),
      }, { signal: abortController.signal })

      if (resID !== requestID)
        return

      const nid = resID.split(':')[1]
      return { nid, schema } as { nid: string, schema: JSONSchema }
    }
    catch (error) {
      //
    }
  }

  const requestMaps = new Map<string, {
    guid: string
    abortController: AbortController
  }>()

  async function handleCheckReplaceModel(jsonschema: JSONSchema, uuid: string): Promise<{ nid: string, schema: JSONSchema } | void> {
    // 避免重复请求

    if (requestMaps.has(uuid))
      requestMaps.get(uuid)?.abortController.abort()

    const info = {
      guid: guid(),
      abortController: new AbortController(),
    }

    requestMaps.set(uuid, info)

    // 额外携带参数,默认为空
    const extraParams = getParams?.() || {}

    try {
      const params = {
        schema: JSON.stringify(jsonschema),
        requestID: info.guid,
        title: extraParams.title || '',
      } as any

      // 如果是model类型，则需要携带modelID
      if (extraParams.type === 'model')
        params.modelID = extraParams.id

      const { requestID, schema } = await apiCheckReplaceModel(projectID, params, { signal: info.abortController.signal })

      // not match
      if (info.guid !== requestID)
        return

      return {
        nid: info.guid,
        schema,
      } as { nid: string, schema: JSONSchema }
    }
    catch (error) {
      //
    }
  }

  onBeforeUnmount(() => requestMaps.clear())

  return {
    handleIntelligentSchema,
    handleCheckReplaceModel,
  }
}
