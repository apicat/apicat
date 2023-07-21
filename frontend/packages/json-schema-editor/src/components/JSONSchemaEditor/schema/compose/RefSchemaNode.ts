import { DefinitionSchema, JSONSchema, SchemaNodeOptions } from '../../types'
import SchemaNode from '../SchemaNode'

export default class RefSchemaNode extends SchemaNode {
  static SCHEMA_REF_PREFIX = '#/definitions/schemas/'
  static SCHEMA_REF_REGEX = /#\/definitions\/schemas\/(.*)/

  isRefSchemaNode = true

  definitionSchema: DefinitionSchema | null = null
  refSchemaId: number | null = null

  constructor(options: SchemaNodeOptions) {
    super(options)
    this.refSchemaId = this.getRefSchemaId()
    this.definitionSchema = this.findRefSchema()

    if (!this.definitionSchema) {
      return
    }

    if (!this.isRefSelf) {
      this.updateChildNodes()
    }
  }

  static createRefSchemaByRefId(refSchemaId: number | string): JSONSchema {
    return {
      $ref: RefSchemaNode.SCHEMA_REF_PREFIX + refSchemaId,
    }
  }

  getRefSchemaId() {
    return parseInt(this.schema.$ref!.match(RefSchemaNode.SCHEMA_REF_REGEX)?.[1] as string, 10)
  }

  updateChildNodes(schemaSource?: JSONSchema) {
    this.schema = schemaSource || this.schema

    if (!this.schema) {
      throw new Error('ref schema is required::' + JSON.stringify(schemaSource))
    }

    this.refSchemaId = this.getRefSchemaId()
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

  get isRefSelf(): boolean {
    let parent = this.parent

    while (parent) {
      if (parent instanceof RefSchemaNode) {
        if (parent.refSchemaId === this.refSchemaId) {
          return true
        }
      }
      parent = parent.parent
    }

    return false
  }
}
