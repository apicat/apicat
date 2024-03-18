import { DEFAULT_LANGUAGE, SUPPORTED_LANGUAGES } from '@/commons/constant'
import Storage from '@/commons/storage'
import type { Language } from '@/typings/common'

// 异步加载语言文件
export async function loadLocaleMessages(locale: string) {
  const message = ((await import(`./lang/${locale}.ts`)) as any).default
  if (!message)
    return null
  return message
}

// 设置html lang属性
export function setHtmlPageLang(locale: string) {
  document.querySelector('html')?.setAttribute('lang', locale)
}

export function isLocaleSupported(locale: string) {
  return SUPPORTED_LANGUAGES.find((item: Language) => item.lang === locale)
}

export function getPersistedLocale() {
  const locale = Storage.get(Storage.KEYS.LOCALE)
  if (isLocaleSupported(locale))
    return locale

  return null
}

export function setPersistedLocale(locale: string) {
  if (!locale)
    return

  Storage.set(Storage.KEYS.LOCALE, locale)
}

export function guessDefaultLocale() {
  return getPersistedLocale() || DEFAULT_LANGUAGE
}
