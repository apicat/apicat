import { API_URL } from '@/commons'
import Ajax from './Ajax'

export interface MockRequestParams {
  mock_response_code?: string
}

export const mockApiPath = (project_id: string, mock_path: string): string => `/mock/${project_id}${mock_path}`

export const mockServerPath = location.origin + API_URL

export const getMockData = (requestPath: string, data?: MockRequestParams) => Ajax.get(requestPath + '?' + new URLSearchParams(data as any).toString())
