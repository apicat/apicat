import { HttpCodeColorMap } from '@apicat/shared'
import { compile } from 'path-to-regexp'
import { HttpMethodTypeMap } from './constant'

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
