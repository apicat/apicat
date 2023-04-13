import { DEFAULT_LANGUAGE, SUPPORTED_LANGUAGES } from '@/commons/constant'
import Storage from '@/commons/storage'
import { Language } from '@/typings/common'

export const loadLocaleMessages = async (locale: string) => {
  const message = ((await import(`./lang/${locale}.ts`)) as any).default
  if (!message) return null
  return message
}

export function setHtmlPageLang(locale: string) {
  document.querySelector('html')?.setAttribute('lang', locale)
}

export const getUserLocale = () => {
  const locale = window.navigator.language || DEFAULT_LANGUAGE

  return {
    locale: locale,
    localeNoRegion: locale.split('-')[0],
  }
}

export const isLocaleSupported = (locale: string) => {
  return SUPPORTED_LANGUAGES.find((item: Language) => item.lang === locale)
}

export const getPersistedLocale = () => {
  const locale = Storage.get(Storage.KEYS.LOCALE)
  if (isLocaleSupported(locale)) {
    return locale
  }

  return null
}

export const setPersistedLocale = (locale: string) => {
  if (!locale) {
    return
  }
  Storage.set(Storage.KEYS.LOCALE, locale)
}

export const guessDefaultLocale = () => {
  const userPersistedLocale = getPersistedLocale()
  if (userPersistedLocale) {
    return userPersistedLocale
  }

  const { locale, localeNoRegion } = getUserLocale()

  if (isLocaleSupported(locale)) {
    return locale
  }

  if (isLocaleSupported(localeNoRegion)) {
    return localeNoRegion
  }

  return DEFAULT_LANGUAGE
}
