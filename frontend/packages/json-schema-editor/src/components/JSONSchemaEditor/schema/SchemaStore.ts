import { DefinitionSchema, JSONSchema } from '../types'
import SchemaNode, { SchemaType } from './SchemaNode'
import RefSchemaNode from './compose/RefSchemaNode'

export default class SchemaStore {
  sourceSchema: JSONSchema
  nodesMap: Record<string, InstanceType<typeof SchemaNode>> = {}
  definitionSchemas: DefinitionSchema[] = []
  root: SchemaNode | null = null

  changeNotify: (schema?: JSONSchema) => void

  constructor(sourceSchema: JSONSchema, definitionSchemas: DefinitionSchema[], changeNotify?: (schema?: JSONSchema) => void) {
    this.sourceSchema = sourceSchema
    this.definitionSchemas = definitionSchemas
    this.setSchema(sourceSchema)
    this.changeNotify = () => changeNotify && changeNotify(this.root?.schema)
  }

  setSchema(schema: JSONSchema) {
    if (this.sourceSchema !== schema) {
      this.sourceSchema = schema
      this.root = this.createChildNodes('root', this.sourceSchema)
      if (this.root) {
        this.root.isConstantSchemaNode = true
      }
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
    const options = { schema, store: this, parent, name: schemaName }
    let node = new SchemaNode(options)

    if (schema.$ref !== undefined) {
      node = new RefSchemaNode(options)
    } else {
      switch (schema.type) {
        case 'object':
          this.createObjectChildNodes(schema, node)
          break
        case 'array':
          this.createArrayChildNodes(schema, node)
          break
      }
    }

    node = reactive(node)
    parent && parent.childNodes.push(node)

    return node
  }

  createObjectChildNodes(schema: JSONSchema, parent: SchemaNode) {
    const properties = schema.properties || {}
    const propertiesKeys = schema['x-apicat-orders'] || Object.keys(properties)
    schema['x-apicat-orders'] = propertiesKeys
    for (let k of propertiesKeys) {
      this.createChildNodes(k, properties[k], parent)
    }
  }

  createArrayChildNodes(schema: JSONSchema, parent: SchemaNode) {
    this.createChildNodes('items', schema.items as JSONSchema, parent)
  }

  createSchemaNode(schema: JSONSchema, parent?: SchemaNode): SchemaNode | null {
    if (schema.$ref !== undefined) {
      return new RefSchemaNode({ schema, store: this, parent, isTemp: true })
    }

    return null
  }

  changeToSchemaNodeByType(sourceSchemaNode: SchemaNode, type: SchemaType) {
    if (!sourceSchemaNode) {
      throw new Error('sourceSchemaNode is required')
    }

    // TODO valid type

    const { schema: sourceSchema, store } = sourceSchemaNode
    const schema = this.createBasicSchemaStructByType(type, sourceSchema)
    sourceSchemaNode.type = type
    sourceSchemaNode.schema = schema

    if (type === 'array') {
      const itemSchema = this.createBasicSchemaStructByType('string')
      schema.items = itemSchema
      sourceSchemaNode.childNodes = [new SchemaNode({ store, parent: sourceSchemaNode, schema: schema.items, name: 'item' })]
    } else {
      sourceSchemaNode.childNodes = []
    }

    this.changeNotify()
  }

  changeToRefSchemaNode(sourceSchemaNode: SchemaNode | RefSchemaNode, refSchemaId: number) {
    if (!sourceSchemaNode) {
      throw new Error('sourceSchemaNode is required')
    }

    if (sourceSchemaNode instanceof RefSchemaNode) {
      // sourceSchemaNode.schema =
      sourceSchemaNode.updateChildNodes(RefSchemaNode.createRefSchemaByRefId(refSchemaId))

      return
    }

    const { name, parent } = sourceSchemaNode
    const ref = new RefSchemaNode({ schema: RefSchemaNode.createRefSchemaByRefId(refSchemaId), store: this, parent, name })
    // const idx = sourceSchemaNode.remove()
    console.log(ref)
  }

  createBasicSchemaStructByType(type: SchemaType, overriteJsonSchema?: JSONSchema): JSONSchema {
    const { example, description } = overriteJsonSchema || {}
    const schema: JSONSchema = { type, example, description }

    switch (type) {
      case 'object':
        return {
          ...schema,
          properties: {},
          required: [],
          'x-apicat-orders': [],
        }
      case 'array':
        return {
          ...schema,
          items: {},
        }
      default:
        return { ...schema }
    }
  }
}
