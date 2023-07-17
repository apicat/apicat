import { JSONSchema, SchemaNodeOptions } from '../types'
import SchemaNode from './SchemaNode'

// 基础类型Schema节点
export default abstract class BasicTypeSchemaNode extends SchemaNode {
  // schema 类型
  abstract type: string
  constructor(options: SchemaNodeOptions) {
    super(options)
    const { parent } = options
    if (parent && parent.type === 'array') {
      this.isConstantSchemaNode = true
    }
  }
  createDefaultSchema(): JSONSchema {
    return {
      type: this.type,
    }
  }
}
