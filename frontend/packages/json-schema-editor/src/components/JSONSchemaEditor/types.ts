import SchemaNode from './schema/SchemaNode'
import type SchemaStore from './schema/SchemaStore'

import RefSchemaNode from './schema/compose/RefSchemaNode'
import StringSchemaNode from './schema/basic/StringSchemaNode'
import NumberSchemaNode from './schema/basic/NumberSchemaNode'
import ObjectSchemaNode from './schema/basic/ObjectSchemaNode'
import ArraySchemaNode from './schema/basic/ArraySchemaNode'
import BooleanSchemaNode from './schema/basic/BooleanSchemaNode'
import NullSchemaNode from './schema/basic/NullSchemaNode'
import AnySchemaNode from './schema/basic/AnySchemaNode'

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

export declare interface SchemaNodeOptions {
  name?: string
  store: SchemaStore
  schema: JSONSchema
  parent?: SchemaNode | RefSchemaNode | null
  isTemp?: boolean
}
