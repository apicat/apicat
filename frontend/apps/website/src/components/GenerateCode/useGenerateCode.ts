import type { RendererOptions, TargetLanguage } from 'quicktype-core'
import { FetchingJSONSchemaStore, InputData, JSONSchemaInput, quicktype } from 'quicktype-core'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import type { JSONSchema } from '@apicat/editor'
import { getAllCodeGenerateSupportedLanguages } from './constant'
import { useConfig } from './useConfig'
import { generateJSONDataByJSONSchema } from './generateJSONDataByJSONSchema'
import Storage from '@/commons/storage'
import { useDefinitionSchemaStore } from '@/store/definitionSchema'
import { RefPrefixKeys } from '@/commons'

export function useGenerateCode() {
  const { locale } = useI18n()
  const languages = computed(() => getAllCodeGenerateSupportedLanguages(locale.value))
  const { definitionsForCodeGenerate, flatSchemas } = storeToRefs(useDefinitionSchemaStore())

  const code: Ref<string> = ref('')
  const dataModelName: Ref<string> = ref('')
  const apicatSchema: Ref<Definition.Schema | null> = ref(null)
  const currentLanguage: Ref<string> = ref(Storage.get(Storage.KEYS.CODE_GENERATE_LANGUAGE) || languages.value[0].label)
  const currentLanguageOptions: Ref<Record<string, any>> = ref({})
  const currentLanguageForCodeMirror: Ref<string> = ref('')

  const currentLanguageOptionRender = computed(() => {
    const render = languages.value.find(item => item.label === currentLanguage.value)
    const defaultValues: Record<string, string | boolean | undefined> = {}
    const { getDefaultOptions } = useConfig(currentLanguage, currentLanguageOptions)
    const localLangConfig = getDefaultOptions(currentLanguage.value)

    // set default values
    render?.options.forEach((item) => {
      defaultValues[item.name]
        = localLangConfig[item.name] !== undefined ? localLangConfig[item.name] : item.defaultValue
    })

    currentLanguageOptions.value = defaultValues
    currentLanguageForCodeMirror.value = render?.name || ''

    return render
  })

  const currentTargetLanguage = computed(
    () => languages.value.find(item => item.label === currentLanguage.value)?.targetLanguage,
  )

  const getInputData = useMemoize(async (dataModelName: string, jsonSchemaString: string) => {
    const schemaInput = new JSONSchemaInput(new FetchingJSONSchemaStore())
    await schemaInput.addSource({ name: dataModelName || 'ApiCat', schema: jsonSchemaString })
    const inputData = new InputData()
    inputData.addInput(schemaInput)
    return inputData
  })

  async function quicktypeJSONSchema(
    dataModelName: string,
    jsonSchemaString: string,
    lang: string | TargetLanguage,
    rendererOptions: RendererOptions,
  ) {
    return await quicktype({
      inputData: await getInputData(dataModelName, jsonSchemaString),
      lang,
      rendererOptions,
    })
  }

  function transformSchemaForCodeGenerate(definitionSchema: Definition.Schema) {
    try {
      const schema = definitionSchema.schema as any
      let json = JSON.stringify({
        ...schema,
        description: definitionSchema.description || schema.description,
        definitions: definitionsForCodeGenerate.value,
      })
      json = json.replaceAll(RefPrefixKeys.DefinitionSchema.key, RefPrefixKeys.DefinitionSchema.replaceForCodeGenerate)
      return json
    }
    catch (error) {
      return JSON.stringify({ type: 'object' })
    }
  }

  async function _generateCode() {
    try {
      const targetLanguage = unref(currentTargetLanguage) as TargetLanguage

      if (!targetLanguage) {
        if (currentLanguage.value === 'JSON') {
          code.value = JSON.stringify(
            generateJSONDataByJSONSchema(
              apicatSchema.value!.schema as JSONSchema,
              flatSchemas.value,
              apicatSchema.value!.id,
            ),
            null,
            2,
          )
          return
        }

        if (currentLanguage.value === 'JSONSchema') {
          code.value = JSON.stringify(apicatSchema.value!.schema, null, 2)
          return
        }

        code.value = ''
        return
      }

      const { lines } = await quicktypeJSONSchema(
        unref(dataModelName),
        transformSchemaForCodeGenerate(apicatSchema.value!),
        targetLanguage,
        toRaw(currentLanguageOptions.value),
      )
      code.value = lines.join('\n')
    }
    catch (error) {
      code.value = JSON.stringify(error, null, 2)
    }
  }
  function generateCode() {
    loading.value = true
    _generateCode().finally(() => (loading.value = false))
  }
  const loading = ref(false)

  watchDebounced([currentLanguage, dataModelName, apicatSchema], generateCode, { debounce: 200 })
  watch(currentLanguageOptions, generateCode, { deep: true })

  return {
    languages,
    currentLanguage,
    currentLanguageForCodeMirror,
    currentLanguageOptionRender,
    currentLanguageOptions,
    code,
    dataModelName,
    apicatSchema,

    loading,
  }
}
