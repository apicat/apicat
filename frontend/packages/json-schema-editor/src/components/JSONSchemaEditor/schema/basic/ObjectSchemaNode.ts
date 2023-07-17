import { JSONSchema } from '../../types'
import BasicTypeSchemaNode from '../BasicTypeSchemaNode'

export default class ObjectSchemaNode extends BasicTypeSchemaNode {
  type = 'object'

  createDefaultSchema(): JSONSchema {
    const schema = super.createDefaultSchema()
    schema.properties = {}
    schema.required = []
    return schema
  }
}
