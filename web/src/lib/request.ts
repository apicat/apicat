import { MessageError } from './i18n'
import type { Endpoint, Parameter, Pair, RequestDraft, Schema } from './types'

export function sample(schema: Schema | null, depth = 0): unknown {
  if (!schema || depth > 8 || schema.Recursive || schema.Truncated) return null
  if (schema.JSONExample) {
    try {
      return JSON.parse(schema.JSONExample)
    } catch {
      /* use schema */
    }
  }
  if (schema.JSONDefault) {
    try {
      return JSON.parse(schema.JSONDefault)
    } catch {
      /* use display default */
    }
  }
  if (schema.Default !== '') {
    try {
      return JSON.parse(schema.Default)
    } catch {
      return schema.Default
    }
  }
  if (schema.Enum?.length) {
    if (schema.Type === 'boolean') return schema.Enum[0] === 'true'
    return schema.Type === 'integer' || schema.Type === 'number'
      ? Number(schema.Enum[0])
      : schema.Enum[0]
  }
  if (schema.Properties?.length || schema.Type === 'object')
    return Object.fromEntries((schema.Properties || []).map((p) => [p.Name, sample(p, depth + 1)]))
  if (schema.Type === 'array') return schema.Items ? [sample(schema.Items, depth + 1)] : []
  if (schema.AllOf?.length)
    return Object.assign({}, ...schema.AllOf.map((s) => sample(s, depth + 1)))
  if (schema.OneOf?.length || schema.AnyOf?.length)
    return sample((schema.OneOf || schema.AnyOf)![0], depth + 1)
  if (schema.Type === 'null') return null
  if (schema.Type === 'boolean') return false
  if (schema.Type === 'number' || schema.Type === 'integer') return 0
  return ''
}
function parameterValue(p: Parameter) {
  const example = p.Example || p.Schema?.JSONExample
  if (example) {
    try {
      const value = JSON.parse(example)
      return typeof value === 'string' ? value : JSON.stringify(value)
    } catch {
      return example
    }
  }
  return p.Schema?.Default || ''
}
export function initialDraft(endpoint: Endpoint): RequestDraft {
  const parameters = endpoint.Parameters || []
  const rows = (where: string): Pair[] =>
    parameters
      .filter((p) => p.In === where)
      .map((p) => ({
        key: p.Name,
        value: parameterValue(p),
        enabled: p.Required || parameterValue(p) !== '',
      }))
  const media = endpoint.RequestBody?.Content?.[0]
  const headers = rows('header')
  if (media) headers.push({ key: 'Content-Type', value: media.MediaType, enabled: true })
  const cookies = rows('cookie').filter((p) => p.enabled)
  if (cookies.length)
    headers.push({
      key: 'Cookie',
      value: cookies.map((p) => `${p.key}=${p.value}`).join('; '),
      enabled: true,
    })
  let body = ''
  if (media) {
    if (media.MediaType.includes('json'))
      body = media.Example ? pretty(media.Example) : JSON.stringify(sample(media.Schema), null, 2)
    else if (media.Example) {
      try {
        const v = JSON.parse(media.Example)
        body = typeof v === 'string' ? v : ''
      } catch {
        body = media.Example
      }
    }
  }
  return {
    method: endpoint.Method,
    url: endpoint.Path,
    environmentId: '',
    headers,
    query: rows('query'),
    path: rows('path'),
    body,
    timeoutMs: 30000,
  }
}
function pretty(value: string) {
  try {
    return JSON.stringify(JSON.parse(value), null, 2)
  } catch {
    return value
  }
}
export function serializeQuery(row: Pair, parameter?: Parameter): Pair[] {
  if (
    !row.enabled ||
    !parameter ||
    !['array', 'object'].includes(parameter.Schema?.Type || '') ||
    row.value === ''
  )
    return [row]
  let value: unknown
  try {
    value = JSON.parse(row.value)
  } catch {
    throw new MessageError('参数 {name} 需要 JSON {type}；原始查询也可直接填写到地址中。', {
      name: row.key,
      type: parameter.Schema?.Type || 'object',
    })
  }
  const style = parameter.Style || 'form'
  const explode = parameter.Explode ?? style === 'form'
  const pair = (key: string, value: unknown): Pair => ({ key, value: String(value), enabled: true })
  if (Array.isArray(value)) {
    if (value.some((v) => v !== null && typeof v === 'object'))
      throw new MessageError('参数 {name} 不支持嵌套数组或对象', { name: row.key })
    if (style === 'form' && explode) return value.map((v) => pair(row.key, v))
    if (['form', 'spaceDelimited', 'pipeDelimited'].includes(style))
      return [
        pair(
          row.key,
          value.join(style === 'spaceDelimited' ? ' ' : style === 'pipeDelimited' ? '|' : ','),
        ),
      ]
  } else if (value && typeof value === 'object') {
    const entries = Object.entries(value)
    if (entries.some(([, v]) => v !== null && typeof v === 'object'))
      throw new MessageError('参数 {name} 不支持嵌套对象', { name: row.key })
    if (style === 'deepObject') return entries.map(([k, v]) => pair(`${row.key}[${k}]`, v))
    if (style === 'form')
      return explode
        ? entries.map(([k, v]) => pair(k, v))
        : [pair(row.key, entries.flat().join(','))]
  }
  throw new MessageError('参数 {name} 的值或序列化方式 {style} 不受支持，请直接编辑请求地址', {
    name: row.key,
    style,
  })
}
export function executionDraft(
  draft: RequestDraft,
  endpoint: Endpoint,
  variables: Record<string, string> = {},
): RequestDraft {
  const resolve = (value: string) =>
    value.replace(/\{\{\s*([^{}]+?)\s*\}\}/g, (_token, key: string) => {
      const value = variables[key.trim()]
      if (value === undefined)
        throw new MessageError('变量 {name} 未定义或未启用', { name: key.trim() })
      return value
    })
  for (const p of endpoint.Parameters || []) {
    const complex = ['array', 'object'].includes(p.Schema?.Type || '')
    if (
      p.In === 'path' &&
      draft.url.includes(`{${p.Name}}`) &&
      (complex || (p.Style && p.Style !== 'simple'))
    )
      throw new MessageError('路径参数 {name} 的复杂序列化请直接填写到请求地址中', { name: p.Name })
  }
  return {
    ...draft,
    queryResolved: true,
    query: draft.query.flatMap((row) => {
      const parameter = endpoint.Parameters?.find((p) => p.Name === row.key && p.In === 'query')
      if (!row.enabled) return [row]
      return serializeQuery({ ...row, key: resolve(row.key), value: resolve(row.value) }, parameter)
    }),
  }
}
