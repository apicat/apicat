import { JSONSchema } from '../../types'
import BasicTypeSchemaNode from '../BasicTypeSchemaNode'

export default class ArraySchemaNode extends BasicTypeSchemaNode {
  type = 'array'
  isLeaf = false

  createDefaultSchema(): JSONSchema {
    const schema = super.createDefaultSchema()
    schema.items = {}
    return schema
  }
}
