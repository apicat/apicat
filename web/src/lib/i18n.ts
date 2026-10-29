import { useSyncExternalStore } from 'react'
import { en } from './locales/en'

export type Locale = 'zh-CN' | 'en'
export type LanguagePreference = Locale | 'system'
export type MessageKey = keyof typeof en
export type MessageValues = Record<string, string | number>
const storageKey = 'apicat-language'
const listeners = new Set<() => void>()
export function detectLocale(languages: readonly string[]): Locale {
  for (const language of languages) {
    if (/^zh(?:-|$)/i.test(language)) return 'zh-CN'
    if (/^en(?:-|$)/i.test(language)) return 'en'
  }
  return 'en'
}
export function resolveLocale(saved: string | null, languages: readonly string[]): Locale {
  return saved === 'zh-CN' || saved === 'en' ? saved : detectLocale(languages)
}
function readPreference(): LanguagePreference {
  try {
    const saved = localStorage.getItem(storageKey)
    if (saved === 'zh-CN' || saved === 'en') return saved
  } catch {
    /* Use browser preference when storage is unavailable. */
  }
  return 'system'
}
let preference = typeof window === 'undefined' ? ('system' as LanguagePreference) : readPreference()
export function getLocale(): Locale {
  // A deterministic SSR/test default; the browser resolves its own language before rendering.
  return typeof navigator === 'undefined' ? 'zh-CN' : resolveLocale(preference, navigator.languages)
}
export function formatNumber(locale: Locale, value: number, options?: Intl.NumberFormatOptions) {
  return new Intl.NumberFormat(locale, options).format(value)
}
export function isMessageKey(key: string): key is MessageKey {
  return Object.hasOwn(en, key)
}
export function translate(locale: Locale, key: MessageKey, values: MessageValues = {}): string {
  const template = locale === 'en' ? en[key] : key
  // Double braces belong to request variable syntax and are never interpolated here.
  return template.replace(/(?<!\{)\{(\w+)\}(?!\})/g, (token, name: string) => {
    const value = values[name]
    return value === undefined
      ? token
      : typeof value === 'number'
        ? formatNumber(locale, value)
        : value
  })
}
function notify() {
  if (typeof document !== 'undefined') document.documentElement.lang = getLocale()
  listeners.forEach((listener) => listener())
}
export function setLanguage(value: LanguagePreference) {
  preference = value
  try {
    if (value === 'system') localStorage.removeItem(storageKey)
    else localStorage.setItem(storageKey, value)
  } catch {
    /* The in-memory preference still works. */
  }
  notify()
}
function subscribe(listener: () => void) {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}
if (typeof window !== 'undefined') {
  window.addEventListener('languagechange', notify)
  window.addEventListener('storage', (event) => {
    if (event.key === storageKey || event.key === null) {
      preference = readPreference()
      notify()
    }
  })
  document.documentElement.lang = getLocale()
}
export class MessageError extends Error {
  constructor(
    public key: MessageKey,
    public values: MessageValues = {},
  ) {
    super(key)
    this.name = 'MessageError'
  }
}
export function asError(value: unknown): Error {
  return value instanceof Error ? value : new Error(String(value))
}
export function errorText(locale: Locale, error: Error | null): string {
  if (!error) return ''
  return error instanceof MessageError ? translate(locale, error.key, error.values) : error.message
}
export function useI18n() {
  const locale = useSyncExternalStore(subscribe, getLocale, () => 'zh-CN' as Locale)
  const language = useSyncExternalStore(
    subscribe,
    () => preference,
    () => 'system' as LanguagePreference,
  )
  return {
    locale,
    language,
    setLanguage,
    t: (key: MessageKey, values?: MessageValues) => translate(locale, key, values),
    number: (value: number, options?: Intl.NumberFormatOptions) =>
      formatNumber(locale, value, options),
    errorText: (error: Error | null) => errorText(locale, error),
  }
}
