import { loadLocaleMessages, setHtmlPageLang } from './helper'
import { createI18n } from 'vue-i18n'
import { useLocaleStoreWithOut } from '@/store/locale'

export let i18n: ReturnType<typeof createI18n>

export const setI18nLanguage = (locale: string) => {
  ;(i18n.global.locale as any).value = locale
  const localeStore = useLocaleStoreWithOut()
  localeStore.setLocale(locale)
  setHtmlPageLang(locale)
}

export const getLocaleMessage = (locale: string) => i18n.global.getLocaleMessage(locale)

export const setLocaleMessage = (locale: string, message: any) => {
  i18n.global.setLocaleMessage(locale, message)
  setI18nLanguage(locale)
}

export default async () => {
  const localeStore = useLocaleStoreWithOut()
  const locale = localeStore.locale
  const message = (await loadLocaleMessages(locale)) || {}

  i18n = createI18n({
    legacy: false,
    locale,
    messages: {
      [locale]: message,
    },
  })

  setI18nLanguage(locale)

  return i18n
}

export * from './helper'
