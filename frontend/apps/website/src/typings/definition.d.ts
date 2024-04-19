declare module Definition {
  type ResponseTypeEnum = import('@/commons/constant').ResponseTypeEnum
  type SchemaTypeEnum = import('@/commons/constant').SchemaTypeEnum

  type ResponseContentType =
    | 'application/json'
    | 'application/xml'
    | 'text/html'
    | 'text/plain'
    | 'application/octet-stream'
    | 'none'

  interface CreateResponseTreeNode {
    name: string
    type: ResponseTypeEnum
    parentID: number
    header?: string
    content?: string
    description?: string
  }

  interface ResponseDetail {
    id: number
    name: string
    type: ResponseTypeEnum
    parentID: number
    content?: {
      [key in ResponseContentType]?: {
        schema?: object
        examples?: Record<string, any>
      }
    }
    description?: string
    header?: []
  }

  interface ResponseTreeNode {
    id: number
    name: string
    type: ResponseTypeEnum
    parentID: number
    items: ResponseTreeNode[]
    content?: string
    description?: string
    header?: string
  }

  type UpdateResponse = Pick<ResponseDetail, 'id' | 'name' | 'content' | 'description' | 'header'>

  interface RequestSortParams {
    origin: DefinitionOrder
    target: DefinitionOrder
  }

  interface DefinitionOrder {
    parentID: number
    ids: number[]
  }

  interface Schema {
    id: number
    name: string
    type: SchemaTypeEnum
    parentID: number
    description?: string
    schema?: string | object
  }

  interface RequestCreateSchemaWithAI {
    parentID: number
    prompt: string
  }

  type EditSchema = Pick<Definition.Schema, 'id', 'name' | 'description' | 'schema'>

  interface SchemaNode extends Schema {
    items?: SchemaNode[]
  }
}
