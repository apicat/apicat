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

  let requestIDCheckReplaceModel
  let abortControllerCheckReplace: AbortController | null = null

  async function handleCheckReplaceModel(jsonschema: JSONSchema): Promise<{ nid: string, schema: JSONSchema } | void> {
    // 避免重复请求
    abortControllerCheckReplace?.abort()

    requestIDCheckReplaceModel = guid()

    // 额外携带参数,默认为空
    const extraParams = getParams?.() || {}

    try {
      abortControllerCheckReplace = new AbortController()
      const params = {
        schema: JSON.stringify(jsonschema),
        requestID: requestIDCheckReplaceModel,
        title: extraParams.title || '',
      } as any

      // 如果是model类型，则需要携带modelID
      if (extraParams.type === 'model')
        params.modelID = extraParams.id

      const { requestID, schema } = await apiCheckReplaceModel(projectID, params, { signal: abortControllerCheckReplace.signal })

      // not match
      if (requestIDCheckReplaceModel !== requestID)
        return

      return {
        nid: requestIDCheckReplaceModel,
        schema,
      } as { nid: string, schema: JSONSchema }
    }
    catch (error) {
      //
    }
  }

  return {
    handleIntelligentSchema,
    handleCheckReplaceModel,
  }
}
