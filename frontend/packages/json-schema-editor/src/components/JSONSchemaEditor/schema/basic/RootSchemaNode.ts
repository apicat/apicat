import SchemaNode from '../SchemaNode'
import AnySchemaNode from './AnySchemaNode'
import ArraySchemaNode from './ArraySchemaNode'
import BooleanSchemaNode from './BooleanSchemaNode'
import NullSchemaNode from './NullSchemaNode'
import NumberSchemaNode from './NumberSchemaNode'
import ObjectSchemaNode from './ObjectSchemaNode'
import StringSchemaNode from './StringSchemaNode'

import RefSchemaNode from '../compose/RefSchemaNode'
import { JSONSchema, SchemaNodeOptions } from '../../types'

export default class RootSchemaNode extends SchemaNode {
  rootNode: SchemaNode | null

  constructor(options: SchemaNodeOptions) {
    super(options)
    this.rootNode = null
    this.initialize()
  }
  initialize() {
    this.level = -1
    this.createChildNodes('root', this.schema, this)

    // 修复层级
    if (!this.isEmpty()) {
      this.rootNode = this.childNodes[0]
      this.rootNode.parent = null
      this.rootNode.isConstantSchemaNode = true
    }
  }

  createChildNodes(schemaName: string, schema: JSONSchema, parent: SchemaNode) {
    let node = null
    const { store } = parent
    const options = { schema, store, parent }
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
      node.schemaName = schemaName
      node = reactive(node)
      console.log(node.schemaName, node.name)
      parent.insertChild(node)
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
}
