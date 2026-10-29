import { MessageError, isMessageKey, type MessageValues } from './i18n'
let token = ''
export function setToken(value: string) {
  token = value
}
export async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
  let response: Response
  try {
    response = await fetch(`/api${path}`, {
      ...options,
      headers: { 'Content-Type': 'application/json', 'X-Apicat-Token': token, ...options.headers },
    })
  } catch (error) {
    if (error instanceof Error && error.name === 'AbortError') throw error
    throw new MessageError('无法连接本地服务，请检查服务是否正在运行。')
  }
  let data
  try {
    data = await response.json()
  } catch {
    throw new MessageError('服务返回了无效响应 ({status})', { status: response.status })
  }
  if (!response.ok) {
    if (typeof data?.messageKey === 'string' && isMessageKey(data.messageKey)) {
      const values: MessageValues = {}
      for (const [key, value] of Object.entries(data.values || {})) {
        if (typeof value === 'string' || typeof value === 'number') values[key] = value
      }
      throw new MessageError(data.messageKey, values)
    }
    if (typeof data?.error === 'string') throw new Error(data.error)
    throw new MessageError('请求失败 ({status})', { status: response.status })
  }
  return data as T
}
