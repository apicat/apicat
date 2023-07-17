import { JSONSchema, SchemaNodeKey, SchemaOptions } from '../types'
import type SchemaStore from './SchemaStore'

let nodeIdSeed = 0

export default class SchemaNode {
  id: number = 0
  name: string = ''
  schema: JSONSchema = {}
  parent?: SchemaNode
  childNodes: SchemaNode[] = []
  level: number = 0
  isRefSchema: boolean = false
  isExpand: boolean = false
  isConstantNode: boolean = false
  store: SchemaStore
  paths: string[] = []

  constructor(options: SchemaOptions) {
    this.id = ++nodeIdSeed

    const { schema, store, parent } = options
    this.parent = parent
    this.schema = schema
    this.store = store

    if (this.parent) {
      this.level = this.parent.level + 1
    }
  }

  typeName(type: string | string[] | undefined) {
    if (type === undefined) {
      return 'any'
    }

    if (type instanceof Array) {
      return type.length > 1 ? 'other' : type[0]
    }

    return type
  }

  get key(): SchemaNodeKey {
    return this.paths.join('.')
  }

  get schemaName(): string {
    return this.name
  }

  set schemaName(value: string) {
    if (this.isRefSchema || this.isConstantNode) {
      console.log('ref | constant key 不允许修改', this.parent?.schema)
      return
    }

    if (this.parent && this.parent.schema && this.parent.schema.properties) {
      const { properties } = this.parent.schema
      properties[value] = properties[this.name]
      delete properties[this.name]
    }
    this.name = value
  }

  get isRequired(): boolean {
    if (this.parent) {
      return (this.parent.schema.required || []).includes(this.name)
    }
    return false
  }

  set isRequired(value: boolean) {
    if (!this.parent) {
      return
    }
    const schema = this.parent.schema
    if (!schema.required) {
      schema.required = []
    }
    if (value) {
      schema.required.push(this.name)
    } else {
      schema.required = schema.required.filter((item) => item !== this.name)
    }
  }

  // get $root(): SchemaNode | undefined {
  //   let parent = this
  //   if (!this.parent) {
  //     return this
  //   }

  //   while (parent) {
  //     if (parent.name === 'root') {
  //       return parent
  //     }

  //     parent = parent.parent
  //   }

  //   return undefined
  // }
}
