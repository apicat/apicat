import { JSONSchema, SchemaNodeOptions } from '../types'
import type SchemaStore from './SchemaStore'
import RefSchemaNode from './compose/RefSchemaNode'

export type SchemaType = 'string' | 'boolean' | 'number' | 'integer' | 'object' | 'array' | 'null' | 'any' | 'ref' | string

let nodeIdSeed = 0
export default class SchemaNode {
  id: number = 0
  // 层级，渲染阶梯UI
  level: number = 0
  // schema 类型
  type: SchemaType
  // schema 名称
  schemaName: string = ''
  oldSchemaName: string = ''

  // 管理所有node
  store: SchemaStore
  // 是否常量节点 root | item | normal
  isConstantSchemaNode: boolean = false
  // 允许mock
  isAllowMock = true
  // 允许必选
  isDisabledRequired = false
  // 是否展开
  isExpand: boolean = false
  // 原始schema数据
  schema: JSONSchema
  // 父级Schema节点
  parent: SchemaNode | RefSchemaNode | null = null
  // 子级Schema节点
  childNodes: SchemaNode[] = []
  // 是否临时节点
  isTempSchemaNode: boolean = false
  // 是否引用类型节点
  isRefSchemaNode: boolean = false

  isEmptyName = false

  changeNotify: (schema?: JSONSchema) => void

  constructor(options: SchemaNodeOptions) {
    // TODO valid schema type

    this.id = nodeIdSeed++
    const { store, schema, parent, name, isTemp = false } = options
    this.store = store
    this.parent = parent ?? null

    // 临时节点
    this.isTempSchemaNode = isTemp

    this.schema = markRaw(this.mergeDefaultSchemaStruct(schema))
    this.type = this.getType(this.schema)
    this.schemaName = name || ''

    // 保存节点
    store.register(this)

    this.changeNotify = () => store.changeNotify && store.changeNotify()

    // 初始化
    this.initialize()
  }

  initialize() {
    if (this.parent) {
      this.level = this.parent.level + 1
    }

    // array items
    if (this.parent && this.parent.type === 'array') {
      this.isConstantSchemaNode = true
    }
  }

  mergeDefaultSchemaStruct(schema?: JSONSchema): JSONSchema {
    return schema ?? {}
  }

  get name() {
    return this.isEmptyName ? '' : this.schemaName
  }

  set name(newName: string) {
    if (this.isConstantSchemaNode || this.isInRefSchemaNode) {
      return
    }

    this.isEmptyName = !newName

    if (this.isEmptyName) {
      return
    }

    this.oldSchemaName = this.schemaName
    this.schemaName = newName

    // 临时节点 -> 新增节点时与原schema进行关联引用
    if (this.isTempSchemaNode) {
      this.isTempSchemaNode = false
      this.parent && this.parent.addSchemaPropertyField(this)
    }
    // 修改节点名称
    else {
      this.parent && this.parent.updateSchemaPropertyField(this)
    }

    this.oldSchemaName = ''

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
    if (this.isInRefSchemaNode || this.isTempSchemaNode) {
      return
    }
    this.schema.example = value
    this.changeNotify()
  }

  get description(): string {
    return this.schema.description || ''
  }

  set description(value: string) {
    if (this.isInRefSchemaNode || this.isTempSchemaNode) {
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

  get isLeaf(): boolean {
    const childNodes = this.childNodes
    return !childNodes || childNodes.length === 0
  }

  insertChild(child: SchemaNode | RefSchemaNode, index?: number): void {
    child = reactive(child)
    child.parent = this
    child.initialize()

    if (typeof index === 'undefined' || index < 0) {
      this.childNodes.push(child)
    } else {
      this.childNodes.splice(index, 0, child)
    }

    if (!child.isTempSchemaNode) {
      this.changeNotify()
    }
  }

  remove(): number {
    if (this.parent) {
      return this.parent.removeChild(this)
    }
    return -1
  }

  removeChild(child: SchemaNode | RefSchemaNode): number {
    const index = this.childNodes.indexOf(child)

    if (index > -1) {
      this.store && this.store.deregisterNode(child)
      child.parent?.removeSchemaPropertyField(child)
      child.parent = null
      this.childNodes.splice(index, 1)
    }

    return index
  }

  expand(): void {
    this.isExpand = true
  }

  collapse(): void {
    this.isExpand = false
  }

  isEmpty() {
    return this.childNodes.length === 0
  }

  getType(schema: JSONSchema) {
    if (schema.type === undefined && schema.$ref !== undefined) {
      return 'ref'
    }

    if (schema.type instanceof Array) {
      return 'any'
    }

    return schema.type as string
  }

  updateSchemaPropertyField(child: SchemaNode) {
    const { properties = {} } = this.schema
    const { oldSchemaName, name } = child
    if (!properties[oldSchemaName]) {
      throw new Error('updateSchemaPropertyField error')
    }

    properties[name] = child.schema
    delete properties[oldSchemaName]

    this.updateSchemaPropertyRequired(child)
    this.updateSchemaPropertyOrder(child)
  }

  addSchemaPropertyField(child: SchemaNode): this {
    const { schema } = this
    const property = schema.properties || {}
    property[child.name] = child.schema
    schema.properties = property

    this.updateSchemaPropertyRequired(child)
    this.updateSchemaPropertyOrder(child)
    return this
  }

  removeSchemaPropertyField(child: SchemaNode): this {
    const { schema } = this
    const { name } = child
    if (!name) {
      throw new Error('removeSchemaPropertyField failed, name is empty')
    }

    const { properties = {}, required = [] } = schema
    const order = this.schema['x-apicat-orders'] || []

    delete properties[name]

    // remove required
    schema.required = required.filter((item) => item !== name)
    // remove order
    schema['x-apicat-orders'] = order.filter((item) => item !== name)

    this.changeNotify()
    return this
  }

  updateSchemaPropertyRequired(child: SchemaNode): this {
    const { required } = this.schema
    const { oldSchemaName, name } = child
    this.schema.required = (required || []).map((field) => (field === oldSchemaName ? name : field))

    return this
  }

  updateSchemaPropertyOrder(child: SchemaNode): this {
    const { oldSchemaName, name } = child
    let orders = this.schema['x-apicat-orders'] || []

    // 移除旧schemaname
    orders = orders.filter((one) => one !== oldSchemaName)

    const index = this.childNodes.indexOf(child)

    // 添加新schemaname
    if (index < 0) {
      orders.push(name)
    } else {
      orders.splice(index, 0, name)
    }

    this.schema['x-apicat-orders'] = orders
    return this
  }
}
