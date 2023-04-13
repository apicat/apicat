import { SUPPORTED_LANGUAGES } from '@/commons/constant'
import { i18n, setI18nLanguage, getLocaleMessage, loadLocaleMessages, setLocaleMessage } from '@/i18n'
import { useLocaleStoreWithOut } from '@/store/locale'

export function useLocale() {
  const localeStore = useLocaleStoreWithOut()
  const locale = localeStore.locale

  const currentLocal = ref(locale)

  const elementPlusLocale: any = computed(() => (getLocaleMessage(unref(currentLocal)) as any).elementPlusLocale)

  /**
   * 切换多语言
   * @param newLang
   * @param oldlang
   */
  const switchLanguage = async (newLang: string, oldlang: string) => {
    if (i18n.global.availableLocales.includes(newLang)) {
      setI18nLanguage(newLang)
      return
    }

    const message = await loadLocaleMessages(newLang)
    message && setLocaleMessage(newLang, message)
  }

  watch(currentLocal, switchLanguage)

  return {
    currentLocal,
    elementPlusLocale,
    supporteLanguages: SUPPORTED_LANGUAGES,
  }
}
