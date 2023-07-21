import { JSONSchema } from '../types'
import SchemaNode from './SchemaNode'

// 基础类型Schema节点
export default class BasicTypeSchemaNode extends SchemaNode {
  createDefaultSchema(override?: JSONSchema): JSONSchema {
    if (override && override.type === this.type) {
      return {
        ...override,
        type: this.type,
      }
    }

    return {
      type: this.type,
    }
  }
}
