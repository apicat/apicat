import { debounce } from 'lodash-es'
import { getAllCodeGenerateSupportedLanguages } from './constant'
import Storage from '@/commons/storage'
import { useDefinitionSchemaStore } from '@/store/definition'
import { DefinitionSchema } from '../APIEditor/types'
import { FetchingJSONSchemaStore, InputData, JSONSchemaInput, RendererOptions, TargetLanguage, quicktype } from 'quicktype-core'
import { useConfig } from './useConfig'

export const useGenerateCode = () => {
  const languages = getAllCodeGenerateSupportedLanguages()

  const code: Ref<string> = ref('')
  const dataModelName: Ref<string> = ref('')
  const apicatSchema: Ref<DefinitionSchema | null> = ref(null)
  const currentLanguage: Ref<string> = ref(Storage.get(Storage.KEYS.CODE_GENERATE_LANGUAGE) || languages[0].label)
  const currentLanguageOptions: Ref<Record<string, any>> = ref({})

  const currentLanguageOptionRender = computed(() => {
    const render = languages.find((item) => item.label === currentLanguage.value)
    const defaultValues: Record<string, string | boolean | undefined> = {}
    const localLangConfig = getDefaultOptions(currentLanguage.value)

    // set default values
    render?.options.forEach((item) => {
      defaultValues[item.name] = localLangConfig[item.name] !== undefined ? localLangConfig[item.name] : item.defaultValue
    })

    currentLanguageOptions.value = defaultValues

    return render
  })

  const currentTargetLanguage = computed(() => languages.find((item) => item.label === currentLanguage.value)?.targetLanguage)

  const definitionSchemaStore = useDefinitionSchemaStore()
  const { getDefaultOptions } = useConfig(currentLanguage, currentLanguageOptions)

  const getInputData = useMemoize(async (dataModelName: string, jsonSchemaString: string) => {
    const schemaInput = new JSONSchemaInput(new FetchingJSONSchemaStore())
    await schemaInput.addSource({ name: dataModelName, schema: jsonSchemaString })
    const inputData = new InputData()
    inputData.addInput(schemaInput)
    return inputData
  })

  async function quicktypeJSONSchema(dataModelName: string, jsonSchemaString: string, lang: string | TargetLanguage, rendererOptions: RendererOptions) {
    return await quicktype({
      inputData: await getInputData(dataModelName, jsonSchemaString),
      lang,
      rendererOptions,
    })
  }

  watch(
    [dataModelName, currentLanguage, currentLanguageOptions],
    debounce(async () => {
      try {
        const { lines } = await quicktypeJSONSchema(
          unref(dataModelName),
          definitionSchemaStore.transformSchemaForCodeGenerate(apicatSchema.value!),
          unref(currentTargetLanguage) as TargetLanguage,
          toRaw(currentLanguageOptions.value)
        )
        code.value = lines.join('\n')
      } catch (error) {
        code.value = JSON.stringify(error, null, 2)
      }
    }, 200),
    {
      deep: true,
    }
  )

  return {
    languages,
    currentLanguage,
    currentLanguageOptionRender,
    currentLanguageOptions,
    code,
    dataModelName,
    apicatSchema,
  }
}
