import { HttpMethodTypeMap } from '@/commons/constant'
import type { APICatResponse } from '@/components/ResponseForm.vue'
import { HttpDocument } from '@/typings'
import { createDefaultResponseContent, createDefaultSchema } from './createDefaultDefinition'

/**
 * 创建默认HTTP文档数据格式
 */
export const HTTP_REQUEST_NODE_KEY = 'apicat-http-request'
export const HTTP_RESPONSE_NODE_KEY = 'apicat-http-response'
export const HTTP_URL_NODE_KEY = 'apicat-http-url'

export const createRequestDefaultValue = (overwrite?: any) => ({
  globalExcepts: {
    header: [],
    path: [],
    cookie: [],
    query: [],
  },
  parameters: {
    header: [],
    path: [],
    cookie: [],
    query: [],
  },
  content: {},
  ...overwrite,
})

export const createHttpUrlNode = () => ({
  type: HTTP_URL_NODE_KEY,
  attrs: {
    path: '',
    method: HttpMethodTypeMap.get.value,
  },
})

export const createHttpRequestNode = () => ({
  type: HTTP_REQUEST_NODE_KEY,
  attrs: createRequestDefaultValue(),
})

export const createHttpResponse = (overwrite?: Partial<APICatResponse>) => ({
  name: '',
  code: 200,
  description: '',
  content: createDefaultResponseContent(),
  ...overwrite,
})

export const createHttpResponseNode = () => ({
  type: HTTP_RESPONSE_NODE_KEY,
  attrs: {
    list: [createHttpResponse({ name: 'Response Name' })],
  },
})

export const createHttpDocument = (overwrite?: Partial<HttpDocument>): HttpDocument => {
  const content: any = []

  // 添加默认的http请求地址节点
  content.push(createHttpUrlNode())

  // 添加默认的http请求节点
  content.push(createHttpRequestNode())

  // 添加默认的http响应节点
  content.push(createHttpResponseNode())

  return {
    title: '',
    type: 'http',
    content,
    ...overwrite,
  }
}
