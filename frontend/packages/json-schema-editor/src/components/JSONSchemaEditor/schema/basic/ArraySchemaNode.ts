import { JSONSchema } from '../../types'
import BasicTypeSchemaNode from '../BasicTypeSchemaNode'

export default class ArraySchemaNode extends BasicTypeSchemaNode {
  type = 'array'
  isLeaf = false

  createDefaultSchema(override?: JSONSchema): JSONSchema {
    const schema = super.createDefaultSchema(override)
    schema.items = {}
    return schema
  }
}
