import type { JSONSchema } from '@apicat/editor'
import { guid } from '@apicat/shared'
import type { SchemaTreeNode } from '@apicat/components'
import { apiCheckReplaceModel, apiGetAISchema } from '@/api/project/definition/schema'

// AI generated json schema
export function useIntelligentSchema(projectID: string, getParams?: () => any) {
  let isLoading = false
  let requestID

  async function handleIntelligentSchema(jsonschema: JSONSchema, node: SchemaTreeNode): Promise<{ nid: string, schema: JSONSchema } | void> {
    // 避免重复请求
    if (isLoading)
      return

    requestID = `${guid()}:${node.id}`

    try {
      isLoading = true
      const { requestID: resID, schema } = await apiGetAISchema(projectID, {
        schema: JSON.stringify(jsonschema),
        requestID,
        ...getParams?.(),
      })
      isLoading = false

      if (resID !== requestID)
        return

      const nid = resID.split(':')[1]
      return { nid, schema } as { nid: string, schema: JSONSchema }
    }
    catch (error) {
      isLoading = false
    }
  }

  let isLoadingCheckReplaceModel = false
  let requestIDCheckReplaceModel

  async function handleCheckReplaceModel(jsonschema: JSONSchema): Promise<{ nid: string, schema: JSONSchema } | void> {
    if (isLoadingCheckReplaceModel)
      return

    requestIDCheckReplaceModel = guid()

    // 额外携带参数,默认为空
    const extraParams = getParams?.() || {}

    try {
      isLoadingCheckReplaceModel = true
      const params = {
        schema: JSON.stringify(jsonschema),
        requestID: requestIDCheckReplaceModel,
        title: extraParams.title || '',
      } as any

      // 如果是model类型，则需要携带modelID
      if (extraParams.type === 'model')
        params.modelID = extraParams.id

      const { requestID, schema } = await apiCheckReplaceModel(projectID, params)
      isLoadingCheckReplaceModel = false

      // not match
      if (requestIDCheckReplaceModel !== requestID)
        return

      return {
        nid: requestIDCheckReplaceModel,
        schema,
      } as { nid: string, schema: JSONSchema }
    }
    catch (error) {
      isLoadingCheckReplaceModel = false
    }
  }

  return {
    handleIntelligentSchema,
    handleCheckReplaceModel,
  }
}
