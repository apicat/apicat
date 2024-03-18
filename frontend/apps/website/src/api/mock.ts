import { MockAjax } from './Ajax'
import { queryStringify } from '@/commons'

export interface MockRequestParams {
  mock_response_code?: string
}

export async function getMockData(requestUrl: string, method: string, params?: MockRequestParams): Promise<{ headers: any;response: any }> {
  const requestFn = (MockAjax as any)[method.toLowerCase()]
  if (!requestFn)
    throw new Error(`Method ${method} not found`)
  const { data, headers } = await requestFn(requestUrl + queryStringify(params))
  return {
    response: data,
    headers,
  }
}
