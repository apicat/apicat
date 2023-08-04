export declare interface JSONSchema {
  _id?: number | string
  type?: string | string[]
  description?: string
  required?: string[]
  format?: string
  pattern?: string
  properties?: Record<string, JSONSchema>
  additionalProperties?: boolean | JSONSchema
  items?: JSONSchema | boolean
  enum?: unknown[]
  example?: any
  default?: any
  $ref?: string
  'x-apicat-orders'?: string[]
  'x-apicat-temp-prop'?: boolean
  'x-apicat-mock'?: string
}

export const basicTypes = ['string', 'boolean', 'number', 'integer', 'object', 'array']
export const allowMockTypes = basicTypes.filter((t) => t !== 'object')
export function typename(type: string | string[] | undefined) {
  if (type === undefined) {
    return 'any'
  }
  if (type instanceof Array) {
    return type.length > 1 ? 'other' : type[0]
  }
  return type
}

export declare interface APICatSchemaObject {
  id?: number | string
  name: string
  schema: JSONSchema
  required?: boolean
}

export declare type APICatSchemaObjectCustom = APICatSchemaObject & {
  _id?: number | string
  _name?: string
  [key: string]: any
}

export declare interface DefinitionSchema extends APICatSchemaObject {
  id?: number
  parent_id: number
  type: string
  description?: string
}

export enum constNodeType {
  root = '<root>',
  items = '<items>',
}

export declare interface Tree {
  key: string
  label: string
  tempLabel?: string
  type: string
  schema: JSONSchema
  refObj?: APICatSchemaObject
  isSelf?: boolean
  children?: Tree[]
  parent?: Tree
}
