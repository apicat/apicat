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
import { queryStringify, API_URL, Storage } from '@/commons'
import { DocumentShareInfo } from '@/typings'

const baseRestfulApiPath = (project_id: string | number): string => `/projects/${project_id}/collections`
const detailRestfulPath = (project_id: string | number, collection_id: string | number): string => `${baseRestfulApiPath(project_id)}/${collection_id}`
const shareRestfulPath = (project_id: string | number, collection_id: string | number): string => `${detailRestfulPath(project_id, collection_id)}/share`

export const getCollectionList = async (project_id: string, params?: Record<string, any>) => Ajax.get(`${baseRestfulApiPath(project_id)}${queryStringify(params)}`)

export const getCollectionDetail = () =>
  useApi(async ({ project_id, collection_id, ...params }: any) => {
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

export const copyCollection = async (project_id: string, collection_id: string | number) => Ajax.post(`${detailRestfulPath(project_id, collection_id)}`)

export const moveCollection = async (project_id: string, sortParams: { target: any; origin: any }) => QuietAjax.put(`${baseRestfulApiPath(project_id)}/movement`, sortParams)

export const deleteCollection = async (project_id: string, collection_id: string | number) => Ajax.delete(`${detailRestfulPath(project_id, collection_id)}`)

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
// 获取集合分享详情
export const getCollectionShareDetail = async ({ project_id, collection_id }: any) => QuietAjax.get(`${shareRestfulPath(project_id, collection_id)}`)
// 重置集合分享访问秘钥
export const resetSecretToCollection = async ({ project_id, collection_id }: any) => QuietAjax.put(`${shareRestfulPath(project_id, collection_id)}/reset_share_secretkey`)
// 切换集合分享状态
export const switchCollectionShareStatus = async ({ project_id, collection_id, ...params }: any) => Ajax.put(`${shareRestfulPath(project_id, collection_id)}`, params)
// 检查集合密钥是否正确
export const checkCollectionSecret = async ({ doc_public_id, secret_key }: any) => QuietAjax.post(`/projects/${doc_public_id}/share/secretkey_check`, { secret_key })

// 保存项目分享后的访问token
export const setCollectionShareInfo = (doc_public_id: string, shareInfo: DocumentShareInfo) => Storage.set(`${Storage.KEYS.SHARE_PROJECT}${doc_public_id}`, shareInfo, true)
// 获取项目分享后的访问token
export const getCollectionShareInfo = (doc_public_id: string): DocumentShareInfo => Storage.get(`${Storage.KEYS.SHARE_PROJECT}${doc_public_id}`, true)
