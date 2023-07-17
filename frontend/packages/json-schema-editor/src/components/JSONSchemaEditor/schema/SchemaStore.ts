import { DefinitionSchema, JSONSchema } from '../types'
import SchemaNode from './SchemaNode'

export default class SchemaStore {
  source: JSONSchema
  schemasMap: Map<string, InstanceType<typeof SchemaNode>> = new Map()
  nodesMap: Map<string, InstanceType<typeof SchemaNode>> = new Map()
  definitionSchemas: DefinitionSchema[] = []

  constructor(source: JSONSchema) {
    this.source = source
  }

  getAllSchemas(): SchemaNode[] {
    return Array.from(this.schemasMap.values())
  }

  get schemaTypes(): string[] {
    return Array.from(this.schemasMap.keys())
  }

  setDefinitionSchemas(definitionSchemas: DefinitionSchema[]) {
    this.definitionSchemas = definitionSchemas
  }
}
