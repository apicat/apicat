import { defineStore } from 'pinia'
import { guessDefaultLocale, setPersistedLocale } from '@/i18n/helper'
import { pinia } from '@/plugins'

export const useLocaleStore = defineStore('locale', {
  state: () => ({
    locale: guessDefaultLocale(),
  }),

  actions: {
    setLocale: setPersistedLocale,
  },
})

export const useLocaleStoreWithOut = () => useLocaleStore(pinia)
