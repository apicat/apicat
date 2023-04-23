import { JSONSchema } from '@/components/APIEditor/types'

export declare interface APICatCommonResponse {
  id?: number | string
  name: string
  code: number
  description: string
  content?: Record<string, { schema: JSONSchema }>
  $ref?: string
}

export declare interface APICatCommonResponseCustom {
  id: number | string
  expand: boolean
  isLoading: boolean
  isLoaded: boolean
  code?: number
  description?: string
  detail?: APICatCommonResponse
}
