import { DefinitionSchema, JSONSchema } from '../../types'
import SchemaNode from '../SchemaNode'
import type SchemaStore from '../SchemaStore'

export default class RefSchemaNode extends SchemaNode {
  static SCHEMA_REF_PREFIX = '#/definitions/schemas/'
  static SCHEMA_REF_REGEX = /#\/definitions\/schemas\/(.*)/

  isAllowMock = false
  refDefinitionSchema: DefinitionSchema | null = null
  refSchemaId: number = -1

  constructor(store: SchemaStore, schemaSource: JSONSchema) {
    super(store, schemaSource)
    this.initialize(schemaSource)
  }

  validateSchemaType(schema: JSONSchema) {
    if (!schema.$ref) {
      throw new Error('$ref is required::' + JSON.stringify(schema))
    }
  }

  findRefSchema(): DefinitionSchema | null {
    const { refSchemaId, store } = this
    if (!refSchemaId || isNaN(refSchemaId)) {
      return null
    }

    return store.definitionSchemas.find((item) => item.id === refSchemaId) || null
  }

  initialize(schemaSource?: JSONSchema) {
    schemaSource = schemaSource || this.schema

    this.refSchemaId = parseInt(schemaSource.$ref!.match(RefSchemaNode.SCHEMA_REF_REGEX)?.[1] as string, 10)
    this.refDefinitionSchema = this.findRefSchema()

    if (!this.refDefinitionSchema) {
      return
    }

    const { name } = this.refDefinitionSchema
    this.name = name
    // this.updateChildrenNodes(schema)
  }

  // 更新子节点，避免循环依赖问题
  // updateChildrenNodes(schema: JSONSchema) {}
}
