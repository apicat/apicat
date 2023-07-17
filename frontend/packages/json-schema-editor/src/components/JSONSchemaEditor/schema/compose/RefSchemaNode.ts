import { DefinitionSchema, JSONSchema, SchemaNodeOptions } from '../../types'
import SchemaNode from '../SchemaNode'

export default class RefSchemaNode extends SchemaNode {
  static SCHEMA_REF_PREFIX = '#/definitions/schemas/'
  static SCHEMA_REF_REGEX = /#\/definitions\/schemas\/(.*)/

  isRefSchema = true
  isAllowMock = false
  refDefinitionSchema: DefinitionSchema | null = null
  refSchemaId: number = -1
  isRootRefNode: boolean = false
  rootRefSchemaNode: RefSchemaNode | null = null

  constructor(options: SchemaNodeOptions) {
    super(options)
    this.updateChildNodes()
  }

  createDefaultSchema(definitionSchema: DefinitionSchema): JSONSchema {
    return {
      $ref: RefSchemaNode.SCHEMA_REF_PREFIX + definitionSchema.id,
    }
  }

  updateChildNodes(schemaSource?: JSONSchema) {
    schemaSource = schemaSource || this.schema

    if (!schemaSource) {
      throw new Error('ref schema is required::' + JSON.stringify(schemaSource))
    }

    this.refSchemaId = parseInt(schemaSource.$ref!.match(RefSchemaNode.SCHEMA_REF_REGEX)?.[1] as string, 10)
    this.refDefinitionSchema = this.findRefSchema()

    if (!this.refDefinitionSchema) {
      return
    }

    const { name } = this.refDefinitionSchema
    this.schemaName = name

    this.createChildNodes()
  }

  findRefSchema(): DefinitionSchema | null {
    const { refSchemaId, store } = this
    if (!refSchemaId || isNaN(refSchemaId)) {
      return null
    }

    return store.definitionSchemas.find((item) => item.id === refSchemaId) || null
  }

  createChildNodes() {}
}
