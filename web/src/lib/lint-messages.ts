import { translate, type Locale, type MessageKey } from './i18n'
import type { LintHint } from './types'
const summaries: Record<string, MessageKey> = {
  'api-no-description': 'API 缺少 description',
  'endpoint-no-summary': '接口缺少 summary 或 description',
  'endpoint-no-operation-id': '接口缺少 operationId',
  'endpoint-no-error-response': '接口缺少错误响应（4xx/5xx/default）',
  'param-no-description': '参数缺少 description',
  'field-no-description': 'Schema 字段缺少 description',
}
const details: Record<string, MessageKey> = {
  'endpoint-no-summary': '缺少 summary 或 description',
  'endpoint-no-operation-id': '缺少 operationId',
  'endpoint-no-error-response': '缺少错误响应（4xx/5xx/default）',
  'param-no-description': '参数 {name} 缺少 description',
  'field-no-description': '{count} 个 Schema 字段缺少 description',
}
export function lintSummary(locale: Locale, rule: string, fallback: string) {
  return summaries[rule] ? translate(locale, summaries[rule]) : fallback
}
export function lintDetail(locale: Locale, hint: LintHint) {
  const key = details[hint.Rule]
  if (locale === 'en' && hint.Rule === 'field-no-description' && hint.Count === 1)
    return '1 Schema field without description'
  return key ? translate(locale, key, { name: hint.Name, count: hint.Count }) : hint.Rule
}
