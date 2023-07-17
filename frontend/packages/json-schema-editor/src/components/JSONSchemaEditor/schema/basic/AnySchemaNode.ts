// import { JSONSchema } from '../../types'
import BasicTypeSchemaNode from '../BasicTypeSchemaNode'

export default class AnySchemaNode extends BasicTypeSchemaNode {
  type = 'any'

  // createDefaultSchema(): JSONSchema {
  //   const { store } = this
  //   return {
  //     type: store.schemaTypes.filter((type: string) => type !== this.type),
  //   }
  // }
}
