import { HttpCodeColorMap, traverseTree } from '@apicat/shared'
import { compile } from 'path-to-regexp'
import { isEmpty } from 'lodash-es'
import type { JSONSchema } from '@apicat/editor'
import { HttpMethodTypeMap } from './constant'

/**
 * 创建API模块get path
 * @param path api path
 * @param params
 */
export function convertRequestPath(path: string, params: { [key: string]: any }): string {
  return compile(path)(params)
}
export const createRestfulApiPath = convertRequestPath

/**
 * Returns a string representing the query parameters in the URL format.
 *
 * @param data - An object containing the query parameters.
 * @return A string representing the query parameters in the URL format.
 */
export function queryStringify(data?: Record<string, any>): string {
  if (!data || isEmpty(data)) return ''

  const params = new URLSearchParams()
  Object.entries(data).forEach(([key, value]) => {
    if (Array.isArray(value)) value.forEach((value) => params.append(key, (value || '').toString()))
    else value && params.append(key, (value || '').toString())
  })

  return `?${params.toString()}`
}

/**
 * 重置路径参数中空字符串问题
 * @param params
 */
export function resetEmptyPathParams(params?: Record<string, any>): Record<string, any> {
  if (!params || isEmpty(params)) return {}

  Object.keys(params).forEach((key) => {
    if (!params[key]) params[key] = undefined
  })

  return params
}

export function getResponseStatusCodeBgColor(code: number): any {
  const backgroundColor = (HttpCodeColorMap as any)[String(code)[0]]
  return {
    backgroundColor,
  }
}

export function getRequestMethodColor(method: string): any {
  const color = (HttpMethodTypeMap as any)[(method || '').toLowerCase()].color
  return color ?? HttpMethodTypeMap.get.color
}

export function hasRefInSchema(schema: JSONSchema) {
  if (schema.$ref !== undefined) return true

  const properties = schema.properties

  // check child
  switch (schema.type) {
    case 'object':
      if (properties) {
        const keys = Object.keys(properties)
        for (const key of keys) {
          if (hasRefInSchema(properties[key])) return true
        }
      }
      break
    case 'array':
      if (hasRefInSchema(schema.items as JSONSchema)) return true
  }

  return false
}

export function randomArray(arr: any[]) {
  return arr[Math.floor(Math.random() * arr.length)]
}

export const uuid = () => Math.random().toString(36).substring(2, 9)

export const ROW_KEY = '_id'

export function markDataWithKey(data: Record<string, any>, rowKey = ROW_KEY, defaultValue?: any) {
  if (!data || data[rowKey]) return
  Object.defineProperty(data, rowKey, {
    value: defaultValue || data.id || uuid(),
    enumerable: false,
    configurable: false,
    writable: false,
  })
}

export function checkDepthOver<T>(node: T, subKey: keyof T, maxDepth: number) {
  function getin(node: T, depth: number): boolean {
    depth++
    if (depth > maxDepth) return true
    const li = node[subKey] as T[]
    for (let i = 0; i < li.length; i++) {
      const val = li[i]
      if (getin(val, depth)) return true
    }
    return false
  }
  return getin(node, 0)
}
export function checkDepth<T>(node: T, subKey: keyof T, isLeaf: (node: T) => boolean) {
  function getin(node: T, depth: number): number {
    depth++
    const li = node[subKey] as T[]
    const des: number[] = []
    if (!li || li.length <= 0) des.push(isLeaf(node) ? depth - 1 : depth)
    for (let i = 0; i < li.length; i++) {
      des.push(getin(li[i], depth))
    }
    return Math.max(...des)
  }
  return getin(node, 0)
}

export function createTreeMaxDepthFn<T>(subKey: keyof T) {
  return (node: T) => {
    let maxLevel = 0
    traverseTree(
      () => {
        maxLevel++
        return true
      },
      [node] as T[],
      { subKey },
    )
    return maxLevel
  }
}

export function isJSONSchemaContentType(contentType: string) {
  return contentType === 'application/json' || contentType === 'application/xml'
}

export function removeJsonSchemaTempProperty(schema: JSONSchema) {
  const tempKeys = new Set()
  function removeTempProperty(jsonSchema: JSONSchema) {
    if (jsonSchema.type === 'object' && jsonSchema.properties) {
      const ps = jsonSchema.properties || {}
      Object.keys(ps).forEach((propertyName) => {
        const subJsonSchema = ps[propertyName]
        if (subJsonSchema.type === 'object') removeTempProperty(subJsonSchema)

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

export function waitFor(millisec: number): Promise<void> {
  return new Promise((resolve) => setTimeout(() => resolve(), millisec))
}

export function flattenObjectDot(obj: any, prefix = '') {
  return Object.keys(obj).reduce((acc: any, key) => {
    const pre = prefix.length ? `${prefix}.` : ''
    if (typeof obj[key] === 'object') acc = acc.concat(flattenObjectDot(obj[key], pre + key))
    else acc.push({ key: pre + key, value: obj[key] })

    return acc
  }, [])
}

export function flattenObject(obj: any): object {
  let a: any = {}
  Object.keys(obj).forEach((key) => {
    if (typeof obj[key] === 'object') a = { ...a, ...flattenObject(obj[key]) }
    else a[key] = obj[key]
  })
  return a
}

export function notNullRule(msg: string) {
  return [
    {
      required: true,
      message: msg,
      trigger: 'blur',
    },
  ]
}

import isURL from '@/commons/util/isURL'

export function isUrlRule(msg: string, allow_localhost = true, require_protocol = true) {
  return [
    {
      validator(_: any, v: any, c: any) {
        if (!v || !isURL(v, { allow_localhost, require_protocol })) return c(new Error(msg))
        return c()
      },
      trigger: 'blur',
    },
  ]
}
