import { HttpCodeColorMap, traverseTree } from '@apicat/shared'
import { compile } from 'path-to-regexp'
import { HttpMethodTypeMap } from './constant'
import { JSONSchema } from '@/components/APIEditor/types'
import { memoize } from 'lodash-es'

/**
 * 创建API模块get path
 * @param path api path
 * @param params
 * @returns
 */
export const convertRequestPath = (path: string, params: { [key: string]: any }): string => compile(path)(params)

export const getResponseStatusCodeBgColor = (code: number): any => {
  const backgroundColor = (HttpCodeColorMap as any)[String(code)[0]]
  return {
    backgroundColor,
  }
}

export const getRequestMethodColor = (method: string): any => {
  const color = (HttpMethodTypeMap as any)[method].color
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
