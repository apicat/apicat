import { describe, expect, it, vi, afterEach } from 'vitest'
import {
  detectLocale,
  resolveLocale,
  translate,
  MessageError,
  errorText,
  formatNumber,
  setLanguage,
  getLocale,
} from './i18n'
import { api } from './api'
import { lintDetail } from './lint-messages'

afterEach(() => {
  setLanguage('system')
  vi.unstubAllGlobals()
})
describe('interface localization', () => {
  it('keeps a manual choice usable when localStorage is unavailable', () => {
    vi.stubGlobal('navigator', { languages: ['zh-CN'] })
    vi.stubGlobal('localStorage', {
      setItem: () => {
        throw new Error('storage blocked')
      },
      removeItem: () => {
        throw new Error('storage blocked')
      },
    })
    setLanguage('en')
    expect(getLocale()).toBe('en')
    setLanguage('system')
    expect(getLocale()).toBe('zh-CN')
  })
  it('follows browser language with explicit preferences taking priority', () => {
    expect(detectLocale(['zh-CN', 'en'])).toBe('zh-CN')
    expect(detectLocale(['en-GB', 'zh-CN'])).toBe('en')
    expect(detectLocale(['fr-FR', 'zh-TW'])).toBe('zh-CN')
    expect(detectLocale(['de-DE'])).toBe('en')
    expect(resolveLocale('en', ['zh-CN'])).toBe('en')
    expect(resolveLocale('zh-CN', ['en-US'])).toBe('zh-CN')
    expect(resolveLocale('invalid', ['en-US'])).toBe('en')
  })
  it('interpolates once and leaves variable syntax and user values unchanged', () => {
    const name = 'root API integer string {{token}} {type} 中文'
    expect(translate('en', '变量 {name} 未定义或未启用', { name })).toBe(
      `Variable ${name} is undefined or disabled`,
    )
    expect(translate('en', 'Token 或 {{token}}', { token: 'must-not-expand' })).toBe(
      'Token or {{token}}',
    )
    expect(translate('zh-CN', '变量 {name} 未定义或未启用', { name })).toBe(
      `变量 ${name} 未定义或未启用`,
    )
    expect(formatNumber('en', 1234.5, { minimumFractionDigits: 2 })).toBe('1,234.50')
  })
  it('can retranslate an existing error without changing technical values', () => {
    const error = new MessageError('参数 {name} 需要 JSON {type}；原始查询也可直接填写到地址中。', {
      name: 'root',
      type: 'array',
    })
    expect(errorText('en', error)).toContain('root requires a JSON array')
    expect(errorText('zh-CN', error)).toContain('root 需要 JSON array')
    const raw = new Error('API 返回的中文错误 string')
    expect(errorText('en', raw)).toBe(raw.message)
  })
  it('translates structured lint messages and keeps parameter names', () => {
    expect(lintDetail('en', { Rule: 'param-no-description', Name: '中文_query', Count: 0 })).toBe(
      'Parameter 中文_query has no description',
    )
    expect(lintDetail('en', { Rule: 'field-no-description', Name: '', Count: 1 })).toBe(
      '1 Schema field without description',
    )
  })
  it('translates service error metadata but never rewrites response data', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            error: '变量 x 重复',
            messageKey: '变量 {name} 重复',
            values: { name: '中文_x' },
          }),
          { status: 400 },
        ),
      ),
    )
    const error = await api('/workspace').catch((e) => e)
    expect(error).toBeInstanceOf(MessageError)
    if (!(error instanceof Error)) throw new Error('expected a service error')
    expect(errorText('en', error)).toBe('Duplicate variable: 中文_x')
    const payload = { body: '{"string":"字段原文"}', statusText: 'OK' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify(payload))))
    expect(await api('/requests/execute')).toEqual(payload)
  })
  it('rejects invalid service JSON and preserves request cancellation', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('bad-json')))
    await expect(api('/document')).rejects.toBeInstanceOf(MessageError)
    const abort = new DOMException('aborted', 'AbortError')
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(abort))
    await expect(api('/document')).rejects.toBe(abort)
  })
})
