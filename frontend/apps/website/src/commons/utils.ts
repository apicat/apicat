import { HttpCodeColorMap } from '@apicat/shared'
import { compile } from 'path-to-regexp'

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
