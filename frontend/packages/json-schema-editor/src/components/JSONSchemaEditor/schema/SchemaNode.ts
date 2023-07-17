import { JSONSchema } from '../types'
import type SchemaStore from './SchemaStore'

export default abstract class SchemaNode {
  store: SchemaStore

  // schema 名称
  schemaName: string = ''
  // 是否常量节点，item|root
  isConstantSchemaNode: boolean = false
  // 允许mock
  isAllowMock = true
  // 允许必选
  isDisabledRequired = false
  // 原始schema数据
  schema: JSONSchema
  // 父级Schema节点
  parentSchemaNode: SchemaNode | null = null
  // 子级Schema节点
  childNodes: SchemaNode[] = []
  // schema 的所有选项
  renderOptions = []

  constructor(store: SchemaStore, schema?: JSONSchema) {
    this.store = store
    schema = schema ?? this.createDefaultSchema()
    this.validateSchemaType(schema)
    this.schema = schema
  }

  validateSchemaType(schema: JSONSchema): JSONSchema | void {
    return schema
  }

  createDefaultSchema(): JSONSchema {
    return {}
  }

  // schema 名称
  set name(value: string) {
    // todo validate
    this.schemaName = value
  }

  get name(): string {
    return this.schemaName
  }

  // 是否必须必选
  get isRequired(): boolean {
    if (this.parentSchemaNode) {
      return (this.parentSchemaNode.schema.required || []).includes(this.schemaName)
    }
    return false
  }

  set isRequired(value: boolean) {
    if (!this.parentSchemaNode) {
      return
    }

    if (!this.parentSchemaNode.schema.required) {
      this.parentSchemaNode.schema.required = []
    }

    if (value) {
      this.parentSchemaNode.schema.required = [this.schemaName, ...this.parentSchemaNode.schema.required]
    } else {
      this.parentSchemaNode.schema.required = this.parentSchemaNode.schema.required.filter((item) => item !== this.schemaName)
    }
  }
}
