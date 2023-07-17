export type OptionKind = 'primary' | 'secondary'

export type LegalValue = {
  label: string
  value: any
}

export interface OptionDefinition {
  name: string
  type: StringConstructor | BooleanConstructor | ArrayConstructor | NumberConstructor
  kind?: OptionKind
  defaultValue?: any // 默认值
  label?: string // 表单label
  description: string // 描述
  legalValues?: LegalValue[] // 合法值
}
