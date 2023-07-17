import { SchemaOptions } from '../types'
import SchemaNode from './SchemaNode'

export default class ArraySchemaNode extends SchemaNode {
  constructor(options: SchemaOptions) {
    super(options)
    this.isConstantNode = true
    this.name = 'items'
  }
}
