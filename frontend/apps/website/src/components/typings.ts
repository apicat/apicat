export interface Menu {
  text: string
  icon?: string
  image?: string
  elIcon?: any
  divided?: boolean
  onClick?: () => void
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
