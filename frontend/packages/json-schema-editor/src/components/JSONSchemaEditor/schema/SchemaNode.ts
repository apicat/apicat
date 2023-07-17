import { JSONSchema, SchemaNodeOptions } from '../types'
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

  constructor(options: SchemaNodeOptions) {
    this.id = nodeIdSeed++
    const { store, schema, parent } = options
    this.store = store
    this.parent = parent ?? null
    this.schema = schema ?? this.createDefaultSchema()

    // 保存节点
    store.register(this)

    if (parent) {
      this.level = parent.level + 1
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
    if (this.isConstantSchemaNode) {
      console.log(this)
      return
    }
    this.schemaName = value
  }

  // 是否必须必选
  get isRequired(): boolean {
    if (!this.parent) {
      return false
    }

    return (this.parent.schema.required || []).includes(this.schemaName)
  }

  set isRequired(value: boolean) {
    if (!this.parent || this.isDisabledRequired) {
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
    child.level = this.level + 1
    if (typeof index === 'undefined' || index < 0) {
      this.childNodes.push(child)
    } else {
      this.childNodes.splice(index, 0, child)
    }

    this.updateLeafState()
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
      child.parent = null
      this.childNodes.splice(index, 1)
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
}
