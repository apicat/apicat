import { JSONSchema } from '@/components/APIEditor/types'
import type { APICatResponse } from '@/components/ResponseForm.vue'

export interface GlobalParameter {
  id?: string | number
  name: string
  required?: boolean
  schema: JSONSchema
}

export interface GlobalParameters {
  header: GlobalParameter[]
  cookie: GlobalParameter[]
  query: GlobalParameter[]
  path: GlobalParameter[]
}

export interface CommonParam {
  in?: string
  _id?: number | string
  name: string
  required: boolean
  schema: {
    type: string
    default: string
    example: string
    description: string
  }
}

export type ResponseParamDetail = APICatResponse

export interface ResponseList {
  id?: number | string
  code?: number
  description?: string
}

export interface ResponseListCustom extends ResponseList {
  _id: number | string
  expand: boolean
  isLoading: boolean
  isLoaded: boolean
  detail?: APICatResponse
}

export interface RequestParameters {
  parameters: {
    header: CommonParam[]
    path: CommonParam[]
    cookie: CommonParam[]
    query: CommonParam[]
  }
  content?: unknown
}
