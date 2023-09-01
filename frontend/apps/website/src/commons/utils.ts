import { HttpCodeColorMap, traverseTree } from '@apicat/shared'
import { compile } from 'path-to-regexp'
import { HttpMethodTypeMap } from './constant'
import { JSONSchema } from '@/components/APIEditor/types'
import { isEmpty, memoize } from 'lodash-es'

/**
 * 创建API模块get path
 * @param path api path
 * @param params
 * @returns
 */
export const convertRequestPath = (path: string, params: { [key: string]: any }): string => compile(path)(params)
export const createRestfulApiPath = convertRequestPath

/**
 * Returns a string representing the query parameters in the URL format.
 *
 * @param {Record<string, any>} params - An object containing the query parameters.
 * @return {string} A string representing the query parameters in the URL format.
 */
export const queryStringify = (data?: Record<string, any>): string => {
  if (!data || isEmpty(data)) {
    return ''
  }

  const params = new URLSearchParams()
  Object.entries(data).forEach(([key, value]) => {
    if (Array.isArray(value)) {
      value.forEach((value) => params.append(key, (value || '').toString()))
    } else {
      value && params.append(key, (value || '').toString())
    }
  })

  return `?${params.toString()}`
}

/**
 * 重置路径参数中空字符串问题
 * @param params
 * @returns
 */
export const resetEmptyPathParams = (params?: Record<string, any>): Record<string, any> => {
  if (!params || isEmpty(params)) {
    return {}
  }

  Object.keys(params).forEach((key) => {
    if (!params[key]) {
      params[key] = undefined
    }
  })

  return params
}

export const getResponseStatusCodeBgColor = (code: number): any => {
  const backgroundColor = (HttpCodeColorMap as any)[String(code)[0]]
  return {
    backgroundColor,
  }
}

export const getRequestMethodColor = (method: string): any => {
  const color = (HttpMethodTypeMap as any)[(method || '').toLowerCase()].color
  return color ?? HttpMethodTypeMap.get.color
}

export const hasRefInSchema = (schema: JSONSchema) => {
  if (schema.$ref != undefined) {
    return true
  }

  // check child
  switch (schema.type) {
    case 'object':
      const properties = schema.properties
      if (properties) {
        const keys = Object.keys(properties)
        for (let key of keys) {
          if (hasRefInSchema(properties[key])) {
            return true
          }
        }
      }
      break
    case 'array':
      if (hasRefInSchema(schema.items as JSONSchema)) {
        return true
      }
  }

  return false
}

export const randomArray = (arr: any[]) => arr[Math.floor(Math.random() * arr.length)]

export const uuid = () => Math.random().toString(36).substring(2, 9)

export const ROW_KEY = '_id'

export const markDataWithKey = (data: Record<string, any>, rowKey = ROW_KEY, defaultValue?: any) => {
  if (!data || data[rowKey]) return
  Object.defineProperty(data, rowKey, {
    value: defaultValue || data.id || uuid(),
    enumerable: false,
    configurable: false,
    writable: false,
  })
}

export const createTreeMaxDepthFn = (subKey: string) =>
  memoize(function (node) {
    let maxLevel = 0
    traverseTree(
      (item: any) => {
        if (!item._extend.isLeaf) {
          maxLevel++
        }
      },
      [node] as any[],
      { subKey }
    )
    return maxLevel
  })

export const isJSONSchemaContentType = (contentType: string) => contentType == 'application/json' || contentType == 'application/xml'

export const removeJsonSchemaTempProperty = (schema: JSONSchema) => {
  const tempKeys = new Set()
  function removeTempProperty(jsonSchema: JSONSchema) {
    if (jsonSchema.type === 'object' && jsonSchema.properties) {
      const ps = jsonSchema.properties || {}
      Object.keys(ps).forEach((propertyName) => {
        const subJsonSchema = ps[propertyName]
        if (subJsonSchema.type === 'object') {
          removeTempProperty(subJsonSchema)
        }
        // 移除临时属性
        if (subJsonSchema && subJsonSchema['x-apicat-temp-prop'] !== undefined) {
          delete ps[propertyName]
          tempKeys.add(propertyName)
        }
      })
    }

    // 移除临时必选属性
    jsonSchema.required = jsonSchema.required?.filter((item) => !tempKeys.has(item))
    // 移除临时排序属性
    jsonSchema['x-apicat-orders'] = jsonSchema['x-apicat-orders']?.filter((item) => !tempKeys.has(item))
  }

  removeTempProperty(schema)
}
