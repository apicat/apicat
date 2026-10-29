let token = ''
export function setToken(value: string) {
  token = value
}
export async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`/api${path}`, {
    ...options,
    headers: { 'Content-Type': 'application/json', 'X-Apicat-Token': token, ...options.headers },
  })
  const data = await response
    .json()
    .catch(() => ({ error: `服务返回了无效响应 (${response.status})` }))
  if (!response.ok) throw new Error(data.error || `请求失败 (${response.status})`)
  return data as T
}
