export interface Menu {
  text?: string
  content?: HTMLElement
  icon?: string
  image?: string
  elIcon?: any
  divided?: boolean
  onClick?: (...args: any) => void
  iconify?: string
  refText?: Ref<string>
  [key: string]: any
}

export interface SimpleJSONSchema {
  type: string
  default: unknown
  description: string
}

export interface CommonParameter {
  _id?: number
  name: string
  required?: boolean | undefined
  schema?: SimpleJSONSchema
}
