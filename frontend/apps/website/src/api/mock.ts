import { MockAjax } from './Ajax'

export interface MockRequestParams {
  mock_response_code?: string
}

export const mockApiPath = (project_id: string): string => `/mock/${project_id}`

export const mockServerPath = location.origin

export const getMockData = async (requestPath: string, method: string, data?: MockRequestParams) => {
  const requestFn = (MockAjax as any)[method.toLowerCase()]
  if (!requestFn) {
    throw Error(`Method ${method} not found`)
  }
  return await requestFn(requestPath + '?' + new URLSearchParams(data as any).toString())
}
