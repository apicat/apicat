import { DefinitionSchema, JSONSchema, SchemaOptions } from '../types'
import ArraySchemaNode from './ArraySchemaNode'
import ObjectSchemaNode from './ObjectSchemaNode'
import SchemaNode from './SchemaNode'

export default class RefSchemaNode extends SchemaNode {
  static SCHEMA_REF_PREFIX = '#/definitions/schemas/'
  static SCHEMA_REF_REGEX = /#\/definitions\/schemas\/(.*)/

  refSchemaId: number
  refSchema?: DefinitionSchema

  constructor(options: SchemaOptions) {
    super(options)
    this.isRefSchema = true
    const { schema } = this
    this.refSchemaId = parseInt(schema.$ref!.match(RefSchemaNode.SCHEMA_REF_REGEX)?.[1] as string, 10)
    this.refSchema = this.findRefSchema()

    if (this.refSchema) {
      this.name = this.refSchema.name
      this.doCreateSchemaNode(this.refSchema.schema, this)
    }
  }

  findRefSchema() {
    const { refSchemaId, store } = this
    if (!refSchemaId || isNaN(refSchemaId)) {
      return
    }

    return store.definitionSchemas.find((item) => item.id === refSchemaId)
  }

  doCreateSchemaNode(schema: JSONSchema, parent: SchemaNode) {
    const { store, paths } = parent
    const { nodesMap } = store

    if (schema.$ref) {
      parent.childNodes = [new RefSchemaNode({ schema, store, parent })]
    }

    // childNodes = []

    if (schema.type === 'object') {
      schema.required = schema.required || []
      const properties = (schema.properties = schema.properties || {})
      const propertiesKeys = schema['x-apicat-orders'] || Object.keys(properties)
      schema['x-apicat-orders'] = propertiesKeys

      for (let k of propertiesKeys) {
        const objectSchemaNode = new ObjectSchemaNode({ schema: properties[k], store, parent })
        objectSchemaNode.paths = [...paths, 'properties', k]
        objectSchemaNode.name = k

        nodesMap.set(objectSchemaNode.key, objectSchemaNode)
        parent.childNodes.push(objectSchemaNode)

        this.doCreateSchemaNode(objectSchemaNode.schema, objectSchemaNode)
      }
    }

    if (schema.type === 'array') {
      const arraySchemaNode = new ArraySchemaNode({ schema: schema.items as JSONSchema, store, parent })
      arraySchemaNode.paths = [...paths, 'items']
      nodesMap.set(arraySchemaNode.key, arraySchemaNode)
      parent.childNodes.push(arraySchemaNode)
      this.doCreateSchemaNode(schema.items as JSONSchema, arraySchemaNode)
    }
  }
}
