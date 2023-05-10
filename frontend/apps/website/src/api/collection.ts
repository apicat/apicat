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

export const getCollectionList = (project_id: string) => Ajax.get(`/projects/${project_id}/collections`)

export const getCollectionDetail = () =>
  useApi(async ({ project_id, collection_id }: any) => {
    const doc: any = await Ajax.get(`/projects/${project_id}/collections/${collection_id}`)
    try {
      doc.content = JSON.parse(doc.content)
      mergeDocumentContent(doc.content)
    } catch (error) {
      doc.content = createHttpDocument().content
    }
    return doc
  })

export const createCollection = async ({ project_id, ...data }: any) => QuietAjax.post(`/projects/${project_id}/collections`, data)

export const updateCollection = async ({ project_id, collection_id, ...collectionInfo }: any) =>
  QuietAjax.put(`/projects/${project_id}/collections/${collection_id}`, collectionInfo)

export const copyCollection = async (project_id: string, collection_id: string | number) => Ajax.post(`/projects/${project_id}/collections/${collection_id}`)

export const moveCollection = async (project_id: string, sortParams: { target: any; origin: any }) => QuietAjax.put(`/projects/${project_id}/collections/movement`, sortParams)

export const deleteCollection = async (project_id: string, collection_id: string | number) => Ajax.delete(`/projects/${project_id}/collections/${collection_id}`)

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

export const createCollectionByAI = async ({ project_id, ...params }: any, axiosConfig?: any) => Ajax.post(`/projects/${project_id}/ai/collections`, params, axiosConfig)

export const createCollectionWithSchemaByAI = async ({ project_id, schema_id }: any) => Ajax.get(`/projects/${project_id}/ai/collections/name?schema_id=${schema_id}`)
