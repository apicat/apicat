import { getAllCodeGenerateSupportedLanguages } from './constant'

export const useGenerateCode = () => {
  const languages = getAllCodeGenerateSupportedLanguages()

  const currentLang = ref(languages[0].label)
  const code = ref('')
  const currentLanguageOptions: any = ref({})

  const currentLanguageOptionRender = computed(() => {
    const render = languages.find((item) => item.label === currentLang.value)
    const defaultValues: Record<string, string | boolean | undefined> = {}

    render?.options.forEach((item) => {
      defaultValues[item.name] = item.defaultValue
    })

    currentLanguageOptions.value = defaultValues

    return render
  })

  watch(
    currentLanguageOptions,
    () => {
      console.log(currentLanguageOptions.value)
    },
    {
      deep: true,
    }
  )

  return {
    languages,
    currentLang,
    currentLanguageOptionRender,
    currentLanguageOptions,
    code,
  }
}
