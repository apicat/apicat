import {
  createHttpDocument,
  createHttpRequestNode,
  createHttpResponseNode,
  createHttpUrlNode,
  HTTP_REQUEST_NODE_KEY,
  HTTP_RESPONSE_NODE_KEY,
  HTTP_URL_NODE_KEY,
} from '@/views/document/components/createHttpDocument'
import Ajax, { QuietAjax } from './Ajax'
import useApi from '@/hooks/useApi'
import { isEmpty } from 'lodash-es'
import { queryStringify, API_URL } from '@/commons'
import { setShareTokenToParams } from '@/store/share'

const baseRestfulApiPath = (project_id: string | number): string => `/projects/${project_id}/collections`
const detailRestfulPath = (project_id: string | number, collection_id: string | number): string => `${baseRestfulApiPath(project_id)}/${collection_id}`

export const getCollectionList = async (project_id: string, params?: Record<string, any>) => {
  params = setShareTokenToParams(params || {})
  return Ajax.get(`${baseRestfulApiPath(project_id)}${queryStringify(params)}`)
}

export const getCollectionDetail = () =>
  useApi(async ({ project_id, collection_id, ...params }: any) => {
    params = setShareTokenToParams(params || {})
    const doc: any = await Ajax.get(`${detailRestfulPath(project_id, collection_id)}${queryStringify(params)}`)
    try {
      doc.content = JSON.parse(doc.content)
      mergeDocumentContent(doc.content)
    } catch (error) {
      doc.content = createHttpDocument().content
    }
    return doc
  })

export const createCollection = async ({ project_id, ...data }: any) => QuietAjax.post(`${baseRestfulApiPath(project_id)}`, data)

export const updateCollection = async ({ project_id, collection_id, ...params }: any) => QuietAjax.put(`${detailRestfulPath(project_id, collection_id)}`, params)

export const copyCollection = async (project_id: string, collection_id: string | number, params?: any) => Ajax.post(`${detailRestfulPath(project_id, collection_id)}`, params)

export const moveCollection = async (project_id: string, sortParams: { target: any; origin: any }) => QuietAjax.put(`${baseRestfulApiPath(project_id)}/movement`, sortParams)

export const deleteCollection = async (project_id: string, collection_id: string | number, params?: any) =>
  Ajax.delete(`${detailRestfulPath(project_id, collection_id)}${queryStringify(params)}`)

export const exportCollection = ({ project_id, collection_id, ...params }: any) => `${API_URL}${detailRestfulPath(project_id, collection_id)}/data${queryStringify(params)}`

const mergeHttpMethod = (node: any) => {
  const defaultVal = createHttpUrlNode().attrs
  node.attrs = { ...defaultVal, ...node.attrs }
}

const mergeHttpRequest = (node: any) => {
  const defaultVal = createHttpRequestNode().attrs

  if (!node.attrs || isEmpty(node.attrs)) {
    node.attrs = defaultVal
    return
  }

  if (!node.attrs.parameters || isEmpty(node.attrs.parameters)) {
    node.attrs.parameters = defaultVal.parameters
  }

  if (!node.attrs.globalExcepts || isEmpty(node.attrs.globalExcepts)) {
    node.attrs.globalExcepts = defaultVal.globalExcepts
  }

  Object.keys(node.attrs.globalExcepts).forEach((key) => {
    if (isEmpty(node.attrs.globalExcepts[key])) {
      node.attrs.globalExcepts[key] = []
    }
  })

  Object.keys(node.attrs.parameters).forEach((key) => {
    if (isEmpty(node.attrs.parameters[key])) {
      node.attrs.parameters[key] = []
    }
  })

  if (!node.attrs.content) {
    node.attrs.content = defaultVal.content
  }

  node.attrs.parameters = { ...defaultVal.parameters, ...node.attrs.parameters }

  // todo 多个content
  const firstKey = Object.keys(node.attrs.content)[0]
  if (firstKey) {
    node.attrs.content = {
      [firstKey]: node.attrs.content[firstKey],
    }
  }
}

const mergeHttpResponse = (node: any) => {
  const defaultVal = createHttpResponseNode().attrs
  if (!node.attrs || isEmpty(node.attrs)) {
    node.attrs = defaultVal
    return
  }

  if (!node.attrs.list || !node.attrs.list.length) {
    node.attrs.list = defaultVal.list
  }

  node.attrs.list = node.attrs.list.map((item: any) => {
    if (item.$ref) {
      return item
    }

    if (!item.content) {
      item.content = defaultVal.list[0].content
    }

    // todo 多个content
    const firstKey = Object.keys(item.content)[0]
    if (firstKey) {
      item.content = {
        [firstKey]: item.content[firstKey],
      }
    }
    return item
  })
}

const mergeDocumentContent = (content: any) => {
  const defaultContent = createHttpDocument().content
  // empty content
  if (!content || !content.length) return defaultContent

  content.forEach((node: any) => {
    if (node.type === HTTP_URL_NODE_KEY) {
      mergeHttpMethod(node)
    }

    if (node.type === HTTP_REQUEST_NODE_KEY) {
      mergeHttpRequest(node)
    }

    if (node.type === HTTP_RESPONSE_NODE_KEY) {
      mergeHttpResponse(node)
    }
  })
}

// AI创建集合
export const createCollectionByAI = async ({ project_id, ...params }: any, axiosConfig?: any) => Ajax.post(`/projects/${project_id}/ai/collections`, params, axiosConfig)
// AI通过schema创建集合
export const createCollectionWithSchemaByAI = async ({ project_id, schema_id }: any) => Ajax.get(`/projects/${project_id}/ai/collections/name?schema_id=${schema_id}`)

// 文档历史记录列表
export const getDocumentHistoryRecordList = ({ project_id, collection_id }: Record<string, any>) => Ajax.get(`${detailRestfulPath(project_id, collection_id)}/histories`)

// 文档历史记录详情
export const getDocumentHistoryRecordDetail = async ({ project_id, collection_id, history_id }: Record<string, any>) => {
  const doc: any = await Ajax.get(`${detailRestfulPath(project_id, collection_id)}/histories/${history_id}`)
  parseDocumentContent(doc)
  return doc
}

// 文档历史记录对比
export const compareDocument = async ({ project_id, collection_id, ...params }: Record<string, any>) => {
  const data: any = await Ajax.get(`${detailRestfulPath(project_id, collection_id)}/histories/diff${queryStringify(params)}`)
  parseDocumentContent(data.doc1 || {})
  parseDocumentContent(data.doc2 || {})
  return data
}

// 恢复文档
export const restoreDocumentByHistoryRecord = ({ project_id, collection_id, history_id }: Record<string, any>) =>
  QuietAjax.put(`${detailRestfulPath(project_id, collection_id)}/histories/${history_id}/restore`)

const parseDocumentContent = (doc: any) => {
  try {
    doc.content = JSON.parse(doc.content)
    mergeDocumentContent(doc.content)
  } catch (error) {
    doc.content = createHttpDocument().content
  }
}
