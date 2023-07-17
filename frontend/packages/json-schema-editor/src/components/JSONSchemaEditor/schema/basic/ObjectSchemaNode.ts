import { JSONSchema } from '../../types'
import BasicTypeSchemaNode from '../BasicTypeSchemaNode'

export default class ObjectSchemaNode extends BasicTypeSchemaNode {
  type = 'object'
  isLeaf = false

  createDefaultSchema(): JSONSchema {
    const schema = super.createDefaultSchema()
    schema.properties = {}
    schema.required = []
    schema['x-apicat-orders'] = []
    return schema
  }

  set name(value: string) {
    // todo valid name
    this.schemaName = value
  }
}
