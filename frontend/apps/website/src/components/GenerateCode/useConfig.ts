import Storage from '@/commons/storage'

export function useConfig(currentLanguageRef: Ref<string>, langOptionsRef: Ref<Record<string, any>>) {
  const getDefaultOptions = (langKey: string) => (Storage.get(Storage.KEYS.CODE_GENERATE_CONFIG) || {})[langKey] || {}

  const updateDefaultOptions = () => {
    if (currentLanguageRef.value) {
      const all = Storage.get(Storage.KEYS.CODE_GENERATE_CONFIG) || {}
      all[currentLanguageRef.value] = toRaw(langOptionsRef.value)
      Storage.set(Storage.KEYS.CODE_GENERATE_CONFIG, all)
      Storage.set(Storage.KEYS.CODE_GENERATE_LANGUAGE, currentLanguageRef.value)
    }
  }

  watch(langOptionsRef, updateDefaultOptions, { deep: true })

  return {
    getDefaultOptions,
  }
}
