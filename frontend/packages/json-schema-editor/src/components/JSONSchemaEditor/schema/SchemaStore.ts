import { DefinitionSchema, JSONSchema } from '../types'
import SchemaNode from './SchemaNode'

import RootSchemaNode from './basic/RootSchemaNode'
import RefSchemaNode from './compose/RefSchemaNode'

export default class SchemaStore {
  sourceSchema: JSONSchema
  nodesMap: Map<number, InstanceType<typeof SchemaNode>> = new Map()
  definitionSchemas: DefinitionSchema[] = []
  root: InstanceType<typeof RootSchemaNode>

  constructor(sourceSchema: JSONSchema, definitionSchemas: DefinitionSchema[]) {
    this.sourceSchema = sourceSchema
    this.definitionSchemas = definitionSchemas
    this.root = new RootSchemaNode({ schema: this.sourceSchema, store: this })
  }

  register(node: SchemaNode) {
    this.nodesMap.set(node.id, node)
  }

  deregisterNode(node: SchemaNode) {
    node.childNodes.forEach((child: SchemaNode) => {
      this.deregisterNode(child)
    })

    this.nodesMap.delete(node.id)
  }

  setDefinitionSchemas(definitionSchemas: DefinitionSchema[]) {
    this.definitionSchemas = definitionSchemas
    this.nodesMap.forEach((node) => {
      if (node instanceof RefSchemaNode) {
        node.updateChildNodes()
      }
    })
  }
}
