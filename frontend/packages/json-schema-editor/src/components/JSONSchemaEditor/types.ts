import SchemaNode from './model/SchemaNode'
import SchemaNodeV2 from './schema/SchemaNode'
import type SchemaTreeStore from './model/SchemaStore'
import type SchemaStore from './schema/SchemaStore'

import RefSchemaNode from './schema/compose/RefSchemaNode'

export declare interface JSONSchema {
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
}

export const basicTypes = ['string', 'boolean', 'number', 'integer', 'object', 'array']

export declare interface APICatSchemaObject {
  name: string
  schema: JSONSchema
  required?: boolean
}

export declare interface DefinitionSchema extends APICatSchemaObject {
  id?: number
  parent_id?: number
  type: string
  description?: string
}

export enum ConstNodeType {
  root = 'root',
  items = 'items',
}

export declare interface SchemaTreeStoreOptions {
  schema: JSONSchema
  definitionSchemas: DefinitionSchema[]
  expandKeys?: string[]
}

export declare interface SchemaTree {
  key: string
  label: string
  type: string
  schema: JSONSchema
  refObj?: APICatSchemaObject
  children?: SchemaTree[]
  parent?: SchemaTree
}

export declare type SchemaNodeKey = string

export declare interface SchemaData {
  [key: string]: any
}

export declare interface SchemaOptions {
  schema: JSONSchema
  store: SchemaTreeStore
  parent?: SchemaNode
}

export declare interface SchemaNodeOptions {
  store: SchemaStore
  schema?: JSONSchema
  parent?: SchemaNodeV2 | RefSchemaNode
}
