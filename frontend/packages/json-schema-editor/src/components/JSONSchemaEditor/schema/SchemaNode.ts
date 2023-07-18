import { ConstNodeType, JSONSchema, SchemaNodeOptions } from '../types'
import type SchemaStore from './SchemaStore'

let nodeIdSeed = 0
export default abstract class SchemaNode {
  id: number = 0
  // 层级，渲染阶梯UI
  level: number = 0
  // 管理所有node
  store: SchemaStore
  // schema 名称
  schemaName: string = ''
  // 是否常量节点 root | item | normal
  isConstantSchemaNode: boolean = false
  // 允许mock
  isAllowMock = true
  // 允许必选
  isDisabledRequired = false
  // 是否展开
  isExpand: boolean = false
  // 是否叶子
  isLeaf: boolean = true
  // 原始schema数据
  schema: JSONSchema
  // 父级Schema节点
  parent: SchemaNode | null = null
  // 子级Schema节点
  childNodes: SchemaNode[] = []
  // 是否临时节点
  isTempSchemaNode: boolean = false
  // 是否引用类型节点
  isRefSchemaNode: boolean = false

  changeNotify: (schema?: JSONSchema) => void

  abstract type: string

  constructor(options: SchemaNodeOptions) {
    this.id = nodeIdSeed++
    const { store, schema, parent } = options
    this.store = store
    this.parent = parent ?? null
    // raw schema
    this.schema = markRaw(schema ?? this.createDefaultSchema())

    // 保存节点
    store.register(this)

    this.changeNotify = () => {
      store.changeNotify && store.changeNotify()
    }

    this.initialize()
  }

  initialize() {
    this.isConstantSchemaNode = false

    if (this.parent) {
      this.level = this.parent.level + 1
    }

    // root
    if (!this.parent && this.level === 0) {
      this.isConstantSchemaNode = true
    }

    // array items
    if (this.parent && this.parent.type === 'array') {
      this.isConstantSchemaNode = true
    }
  }

  createDefaultSchema(...args: any): JSONSchema
  createDefaultSchema(): JSONSchema {
    return {}
  }

  get name(): string {
    return this.schemaName
  }

  set name(value: string) {
    if (!this.schemaName && this.isTempSchemaNode && value) {
      this.isTempSchemaNode = false
      this.schemaName = value
      if (this.parent) {
        const { schema } = this.parent
        const properties = (schema.properties = schema.properties || {})
        properties[value] = this.schema
        const orders = schema['x-apicat-orders'] || []
        orders.push(value)
        schema['x-apicat-orders'] = orders
        this.changeNotify()
      }

      return
    }

    if (this.isConstantSchemaNode || this.isInRefSchemaNode) {
      return
    }

    // validate
    // if(!this.validate(value)) {
    //   return
    // }

    if (this.parent) {
      const { schema } = this.parent
      const properties = (schema.properties = schema.properties || {})
      properties[value] = properties[this.schemaName]

      // 更新必选
      if (this.isRequired) {
        schema.required = schema.required?.map((one) => (one === this.schemaName ? value : one))
      }

      // 更新排序
      const orders = schema['x-apicat-orders'] || []
      orders[orders?.indexOf(this.schemaName)] = value
      schema['x-apicat-orders'] = orders

      // 移除旧key
      delete properties[this.schemaName]
    }

    this.schemaName = value

    this.changeNotify()
  }

  // 是否必须必选
  get isRequired(): boolean {
    if (!this.parent) {
      return false
    }

    return (this.parent.schema.required || []).includes(this.schemaName)
  }

  set isRequired(value: boolean) {
    if (!this.parent || this.isDisabledRequired || this.isInRefSchemaNode) {
      return
    }

    if (!this.parent.schema.required) {
      this.parent.schema.required = []
    }

    if (value) {
      this.parent.schema.required = [this.schemaName, ...this.parent.schema.required]
    } else {
      this.parent.schema.required = this.parent.schema.required.filter((item) => item !== this.schemaName)
    }

    this.changeNotify()
  }

  get example(): string {
    return this.schema.example
  }

  set example(value: string) {
    if (this.isInRefSchemaNode) {
      return
    }
    this.schema.example = value
    this.changeNotify()
  }

  get description(): string {
    return this.schema.description || ''
  }

  set description(value: string) {
    if (this.isInRefSchemaNode) {
      return
    }
    this.schema.description = value
    this.changeNotify()
  }

  get nextSibling(): SchemaNode | null {
    const parent = this.parent
    if (parent) {
      const index = parent.childNodes.indexOf(this)
      if (index > -1) {
        return parent.childNodes[index + 1]
      }
    }
    return null
  }

  get previousSibling(): SchemaNode | null {
    const parent = this.parent
    if (parent) {
      const index = parent.childNodes.indexOf(this)
      if (index > -1) {
        return index > 0 ? parent.childNodes[index - 1] : null
      }
    }
    return null
  }

  insertChild(child: SchemaNode, index?: number) {
    child.parent = this
    child.initialize()

    if (typeof index === 'undefined' || index < 0) {
      this.childNodes.push(child)
    } else {
      this.childNodes.splice(index, 0, child)
    }

    this.updateLeafState()
    if (!child.isTempSchemaNode) {
      this.changeNotify()
    }
  }

  insertBefore(child: SchemaNode, ref: SchemaNode): void {
    let index
    if (ref) {
      index = this.childNodes.indexOf(ref)
    }
    this.insertChild(child, index)
  }

  insertAfter(child: SchemaNode, ref: SchemaNode): void {
    let index
    if (ref) {
      index = this.childNodes.indexOf(ref)
      if (index !== -1) index += 1
    }
    this.insertChild(child, index)
  }

  remove(): void {
    const parent = this.parent
    if (parent) {
      parent.removeChild(this)
    }
  }

  removeChild(child: SchemaNode): void {
    const index = this.childNodes.indexOf(child)

    if (index > -1) {
      this.store && this.store.deregisterNode(child)
      console.log(child, child.parent)
      child.parent = null

      this.childNodes.splice(index, 1)

      this.changeNotify()
    }

    this.updateLeafState()
  }

  expand(): void {
    this.isExpand = true
  }

  collapse(): void {
    this.isExpand = false
  }

  updateLeafState(): void {
    const childNodes = this.childNodes
    this.isLeaf = !childNodes || childNodes.length === 0
  }

  isEmpty() {
    return this.childNodes.length === 0
  }

  get isInRefSchemaNode() {
    let result = false
    let parent = this.parent
    while (parent) {
      if (parent && parent.isRefSchemaNode) {
        result = true
        break
      }
      parent = parent.parent
    }
    return result
  }

  addProperty() {}

  updateProperty() {}

  deleteProperty() {}
}
