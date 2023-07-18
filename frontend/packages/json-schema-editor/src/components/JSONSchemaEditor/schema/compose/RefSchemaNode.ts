import { DefinitionSchema, JSONSchema, SchemaNodeOptions } from '../../types'
import SchemaNode from '../SchemaNode'

export default class RefSchemaNode extends SchemaNode {
  static SCHEMA_REF_PREFIX = '#/definitions/schemas/'
  static SCHEMA_REF_REGEX = /#\/definitions\/schemas\/(.*)/

  type = 'ref'

  isRefSchemaNode = true
  isAllowMock = false
  definitionSchema: DefinitionSchema | null = null
  refSchemaId: number = -1

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
    this.definitionSchema = this.findRefSchema()

    if (!this.definitionSchema) {
      return
    }

    const { name, schema } = this.definitionSchema
    this.level = this.level - 1
    const linkSchemaNode = this.store.createChildNodes(name, schema, this)
    this.childNodes = linkSchemaNode ? linkSchemaNode.childNodes : []
    this.level = this.level + 1
  }

  findRefSchema(): DefinitionSchema | null {
    const { refSchemaId, store } = this
    if (!refSchemaId || isNaN(refSchemaId)) {
      return null
    }

    return store.definitionSchemas.find((item) => item.id === refSchemaId) || null
  }
}
