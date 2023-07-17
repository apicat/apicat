import { JSONSchema } from '../types'
import SchemaNode from './SchemaNode'

// 基础类型Schema节点
export default abstract class BasicTypeSchemaNode extends SchemaNode {
  // schema 类型
  abstract type: string

  validateSchemaType(schema: JSONSchema) {
    if (typeof schema.type !== 'string' || this.type !== schema.type) {
      throw new Error(`input schema type ${schema.type} is not equal to this schema type ${this.type}`)
    }
  }

  createDefaultSchema(): JSONSchema {
    return {
      type: this.type,
    }
  }
}
