import type { JSONSchema } from '@apicat/editor'
import { guid } from '@apicat/shared'
import type { SchemaTreeNode } from '@apicat/components'
import { apiGetAISchema } from '@/api/project/definition/schema'

// AI generated json schema
export function useIntelligentSchema(getParams?: () => any) {
  let isLoading = false
  let requestID

  const handleIntelligentSchema = async (jsonschema: JSONSchema, node: SchemaTreeNode) => {
    // 避免重复请求
    if (isLoading)
      return

    requestID = `${guid()}:${node.id}`

    try {
      isLoading = true
      const { requestID: resID, schema } = await apiGetAISchema({
        schema: jsonschema,
        requestID,
        ...getParams?.(),
      })
      isLoading = false

      if (resID !== requestID)
        return

      const nid = resID.split(':')[1]
      return {
        nid,
        schema,
      }
    }
    catch (error) {
      isLoading = false
    }
  }

  return {
    handleIntelligentSchema,
  }
}
