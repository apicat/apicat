/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_MOCK_SERVER: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

type Nullable<T> = T | null
type Arrayable<T> = T | T[]
type Awaitable<T> = Promise<T> | T
type PageableQueryParam<T> = {
  page?: number
  page_size?: number
} & T

declare interface CommonResponseMessageForMessageTemplate {
  accessToken?: string
  description?: string
  emoji?: string
  title?: string

  message?: string
}

declare namespace ResponseAPI {
  interface Response<T> {
    code: number
    data: T
    msg: string
  }
}

// 迭代 & 项目左侧栏选中的key
declare type IterationSelectedKey = 'all' | 'create' | string
declare type ProjectGroupSelectKey = 'all' | 'followed' | 'my' | 'create' | number | null
declare interface SwitchProjectGroupInfo {
  key: ProjectGroupSelectKey
  title: string
}
