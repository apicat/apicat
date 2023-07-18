import { DefinitionSchema, JSONSchema } from '../types'
import SchemaNode from './SchemaNode'
import AnySchemaNode from './basic/AnySchemaNode'
import ArraySchemaNode from './basic/ArraySchemaNode'
import BooleanSchemaNode from './basic/BooleanSchemaNode'
import NullSchemaNode from './basic/NullSchemaNode'
import NumberSchemaNode from './basic/NumberSchemaNode'
import ObjectSchemaNode from './basic/ObjectSchemaNode'
import StringSchemaNode from './basic/StringSchemaNode'
import RefSchemaNode from './compose/RefSchemaNode'

export default class SchemaStore {
  sourceSchema: JSONSchema
  nodesMap: Record<string, InstanceType<typeof SchemaNode>> = {}
  definitionSchemas: DefinitionSchema[] = []
  root: SchemaNode | null

  changeNotify?: (schema?: JSONSchema) => void

  constructor(sourceSchema: JSONSchema, definitionSchemas: DefinitionSchema[], changeNotify?: (schema?: JSONSchema) => void) {
    this.sourceSchema = sourceSchema
    this.definitionSchemas = definitionSchemas
    this.root = this.createChildNodes('root', this.sourceSchema)
    this.changeNotify = () => changeNotify && changeNotify(this.root?.schema)
  }

  setSchema(schema: JSONSchema) {
    if (this.sourceSchema !== schema) {
      this.sourceSchema = schema
      this.root = this.createChildNodes('root', this.sourceSchema)
    }
  }

  register(node: SchemaNode) {
    this.nodesMap[node.id] = node
  }

  deregisterNode(node: SchemaNode) {
    node.childNodes.forEach((child: SchemaNode) => {
      this.deregisterNode(child)
    })
    delete this.nodesMap[node.id]
  }

  setDefinitionSchemas(definitionSchemas: DefinitionSchema[]) {
    if (this.definitionSchemas === definitionSchemas) {
      return
    }
    this.definitionSchemas = definitionSchemas
    Object.keys(this.nodesMap).forEach((key: string) => {
      const node = this.nodesMap[key]
      if (node instanceof RefSchemaNode) {
        node.updateChildNodes()
      }
    })
  }

  createChildNodes(schemaName: string, schema: JSONSchema, parent?: SchemaNode): SchemaNode | null {
    let node = null
    const options = { schema, store: this, parent }
    if (schema.$ref !== undefined) {
      node = new RefSchemaNode(options)
    } else {
      switch (schema.type) {
        case 'object':
          node = new ObjectSchemaNode(options)
          this.createObjectChildNodes(schema, node)
          break

        case 'array':
          node = new ArraySchemaNode(options)
          this.createArrayChildNodes(schema, node)
          break

        case 'boolean':
          node = new BooleanSchemaNode(options)
          break

        case 'integer':
        case 'number':
          node = new NumberSchemaNode(options)
          break

        case 'string':
          node = new StringSchemaNode(options)
          break

        case 'null':
          node = new NullSchemaNode(options)
          break

        case 'any':
          node = new AnySchemaNode(options)
          break
      }
    }

    if (node) {
      node.schemaName = schemaName || ''
      node = reactive(node)
      parent && parent.childNodes.push(node)
    }

    return node
  }

  createObjectChildNodes(schema: JSONSchema, parent: SchemaNode) {
    schema.required = schema.required || []
    const properties = (schema.properties = schema.properties || {})
    const propertiesKeys = schema['x-apicat-orders'] || Object.keys(properties)
    schema['x-apicat-orders'] = propertiesKeys
    for (let k of propertiesKeys) {
      this.createChildNodes(k, properties[k], parent)
    }
  }

  createArrayChildNodes(schema: JSONSchema, parent: SchemaNode) {
    this.createChildNodes('items', schema.items as JSONSchema, parent)
  }

  createSchemaNodeByJsonSchema(jsonSchema: JSONSchema) {
    const node = this.createChildNodes('', jsonSchema)
    if (node) {
      node.isTempSchemaNode = true
    }
    return node
  }
}
