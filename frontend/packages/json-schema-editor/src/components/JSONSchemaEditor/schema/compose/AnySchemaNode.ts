import { JSONSchema } from '../../types'
import SchemaNode from '../SchemaNode'

export default class AnySchemaNode extends SchemaNode {
  type = 'any'

  createDefaultSchema(): JSONSchema {
    const { store } = this
    return {
      type: store.schemaTypes.filter((type: string) => type !== this.type),
    }
  }
}
