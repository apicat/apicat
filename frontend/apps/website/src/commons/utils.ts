import { HttpCodeColorMap } from '@apicat/shared'
import { compile } from 'path-to-regexp'
import { HttpMethodTypeMap } from './constant'
import { JSONSchema } from '@/components/APIEditor/types'

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
