// eslint-disable-next-line ts/ban-ts-comment
// @ts-nocheck
import elementPlusLocaleMessage from 'element-plus/es/locale/lang/en'
import { defineStore } from 'pinia'
import { createI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import localizedFormat from 'dayjs/plugin/localizedFormat'
import { DEFAULT_LANGUAGE, SUPPORTED_LANGUAGES } from '@/commons/constant'
import { pinia } from '@/plugins'
import { guessDefaultLocale, loadLocaleMessages, setHtmlPageLang, setPersistedLocale } from '@/i18n/helper'
import type { Language } from '@/typings/common'
import enUS from '@/i18n/lang/en-US'

dayjs.extend(localizedFormat)

export const i18n = createI18n({
  locale: DEFAULT_LANGUAGE,
  fallbackLocale: DEFAULT_LANGUAGE,
  messages: {
    [DEFAULT_LANGUAGE]: enUS,
  },
})

export const useLocaleStore = defineStore('locale', {
  state: () => {
    return {
      i18n,

      t: i18n.global.t,
      locale: DEFAULT_LANGUAGE,
      localeMessages: enUS,
      supportedLocales: SUPPORTED_LANGUAGES,
    }
  },

  getters: {
    languagesForSelect: state => state.supportedLocales.map(item => ({ label: item.name, value: item.lang })),
    elementPlusLocaleMessage: state => state.localeMessages.elementPlusLocale || elementPlusLocaleMessage,
    acCompLocaleMessage: state => state.localeMessages.components,
    acEditorLocaleMessage: state => state.localeMessages.editor,
    dayjsLocale: state => state.localeMessages.dayjsLocale,
  },

  actions: {
    async switchLanguage(locale: Language['lang']) {
      if (this.locale === locale || !this.supportedLocales.map(item => item.lang).includes(locale))
        return

      setHtmlPageLang(locale)

      const isExist = this.i18n.global.availableLocales.includes(locale)

      const message = isExist ? this.i18n.global.getLocaleMessage(locale) : await loadLocaleMessages(locale)

      if (message) {
        setPersistedLocale(locale)
        this.localeMessages = message

        this.i18n.global.setLocaleMessage(locale, message)
        this.i18n.global.locale = locale
        this.locale = locale

        dayjs.locale(this.dayjsLocale.name, this.dayjsLocale)

        const descriptionEl = document.querySelector('meta[name="description"]')
        descriptionEl && descriptionEl.setAttribute('content', this.t('app.description'))
      }
    },

    async initLocale() {
      await this.switchLanguage(guessDefaultLocale())
    },
  },
})

export default useLocaleStore

export const useLocaleStoreWithOut = () => useLocaleStore(pinia)
