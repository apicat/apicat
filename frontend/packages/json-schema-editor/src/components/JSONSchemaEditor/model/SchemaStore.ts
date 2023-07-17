import { DefinitionSchema, JSONSchema, SchemaTreeStoreOptions } from '../types'
import SchemaNode from './SchemaNode'
import ObjectSchemaNode from './ObjectSchemaNode'
import RootSchemaNode from './RootSchemaNode'
import ArraySchemaNode from './ArraySchemaNode'
import RefSchemaNode from './RefSchemaNode'

export default class SchemaTreeStore {
  // 根节点
  root: InstanceType<typeof SchemaNode>

  // 模型集合
  definitionSchemas: DefinitionSchema[] = []

  // 展开的节点Key集合
  expandKeys: string[] = []

  // nodesMap
  nodesMap: Map<string, InstanceType<typeof SchemaNode>> = new Map()

  constructor({ schema, definitionSchemas }: SchemaTreeStoreOptions) {
    this.root = new RootSchemaNode({ schema, store: this })
    this.definitionSchemas = definitionSchemas

    SchemaTreeStore.doCreateSchemaNode(this.root)
  }

  register(node: SchemaNode) {
    this.nodesMap.set(node.key, node)
  }

  deregisterNode(node: SchemaNode) {
    node.childNodes.forEach((child: SchemaNode) => {
      this.deregisterNode(child)
    })

    this.nodesMap.delete(node.key)
  }

  static doCreateSchemaNode(parent: SchemaNode) {
    const { store, schema, paths } = parent
    const { nodesMap } = store

    if (schema.$ref) {
      parent.childNodes = [new RefSchemaNode({ schema, store, parent })]
    }

    if (schema.type === 'object') {
      const childNodes = []
      schema.required = schema.required || []
      const properties = (schema.properties = schema.properties || {})
      const propertiesKeys = schema['x-apicat-orders'] || Object.keys(properties)
      schema['x-apicat-orders'] = propertiesKeys

      for (let k of propertiesKeys) {
        let schemeNode = new ObjectSchemaNode({ schema: properties[k], store, parent })

        schemeNode.paths = [...paths, 'properties', k]
        schemeNode.name = k
        nodesMap.set(schemeNode.key, schemeNode)
        childNodes.push(schemeNode)

        parent.childNodes = childNodes

        SchemaTreeStore.doCreateSchemaNode(schemeNode)
      }
    }

    if (schema.type === 'array') {
      let schemeNode = new ArraySchemaNode({ schema: schema.items as JSONSchema, store, parent })
      schemeNode.paths = [...paths, 'items']
      nodesMap.set(schemeNode.key, schemeNode)
      parent.childNodes = [schemeNode]
      SchemaTreeStore.doCreateSchemaNode(schemeNode)
    }
  }
}
