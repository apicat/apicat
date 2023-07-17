import { SchemaOptions } from '../types'
import ObjectSchemaNode from './ObjectSchemaNode'

export default class RootSchemaNode extends ObjectSchemaNode {
  isRoot: boolean = true
  constructor(options: SchemaOptions) {
    super(options)
    this.isConstantNode = true
    this.name = 'root'

    const { schema } = this
    if (!schema.$ref) {
      schema.type = this.typeName(schema.type)
    }
  }
}
