import { JSONSchema } from '@/components/APIEditor/types'

export declare interface APICatCommonResponse {
  code: number
  id?: number | string
  name?: string
  description?: string
  content?: Record<string, { schema: JSONSchema }>
  $ref?: string
}

export declare interface APICatCommonResponseCustom {
  id: number | string
  expand: boolean
  isLocal: boolean
  isLoading: boolean
  isLoaded: boolean
  code?: number
  description?: string
  name?: string
  detail?: APICatCommonResponse
}

export declare interface DefinitionResponse {
  id?: number
  type: string
  name: string
  description?: string
  header?: Array<{
    name: string
    description?: string
    required?: boolean
    schema: JSONSchema
  }>
  content?: Record<string, { schema: JSONSchema }>
}
