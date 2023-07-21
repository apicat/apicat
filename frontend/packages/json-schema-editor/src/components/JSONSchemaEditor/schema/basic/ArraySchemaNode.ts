import { JSONSchema, SchemaNodeOptions } from '../../types'
import BasicTypeSchemaNode from '../BasicTypeSchemaNode'
import SchemaNode from '../SchemaNode'
import StringSchemaNode from './StringSchemaNode'

export default class ArraySchemaNode extends BasicTypeSchemaNode {
  type = 'array'

  constructor(options: SchemaNodeOptions) {
    super(options)
    const { schema, store } = options
    // 创建item默认节点
    if (!schema) {
      const item = new StringSchemaNode({ store, name: 'items', parent: this })
      this.insertChild(item)
    }
  }

  createDefaultSchema(override?: JSONSchema): JSONSchema {
    const schema = super.createDefaultSchema(override)
    schema.items = {}
    return schema
  }

  insertChild(child: SchemaNode, index?: number): void {
    child.parent = this
    child.initialize()

    if (typeof index === 'undefined' || index < 0) {
      this.childNodes.push(child)
    } else {
      this.childNodes.splice(index, 0, child)
    }

    if (!this.isTempSchemaNode) {
      this.changeNotify()
    }
  }
}
