<template>
  <el-dialog v-model="dialogVisible" append-to-body :close-on-click-modal="false" class="fullscree hide-header" destroy-on-close center width="70%">
    <CodeEditor :modelValue="code" readonly lang="json" />
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import { quicktype, InputData, JSONSchemaInput, FetchingJSONSchemaStore, JavaScriptPropTypesTargetLanguage, TypeScriptTargetLanguage } from 'quicktype-core'
import jsonschemaData from './jsonschema'
import CodeEditor from '../APIEditor/CodeEditor.vue'
import { useI18n } from 'vue-i18n'
import { getAllCodeGenerateSupportedLanguages } from './constant'

const { t } = useI18n()
const { dialogVisible, showModel } = useModal()
const langs = getAllCodeGenerateSupportedLanguages()

const code = ref('')

async function quicktypeJSONSchema(targetLanguage: string, typeName: string, jsonSchemaString: string) {
  const schemaInput = new JSONSchemaInput(new FetchingJSONSchemaStore())
  console.log(langs)

  await schemaInput.addSource({ name: typeName, schema: jsonSchemaString })
  const tsLang = new TypeScriptTargetLanguage()
  // const jsPropLang = new JavaScriptPropTypesTargetLanguage()
  const inputData = new InputData()
  inputData.addInput(schemaInput)

  return await quicktype({
    inputData,
    lang: tsLang,
    rendererOptions: {
      'just-types': true,
    },
  })
}

defineExpose({
  show: async () => {
    showModel()
    const { lines } = await quicktypeJSONSchema('typescript', 'User', JSON.stringify(jsonschemaData))
    code.value = lines.join('\n')
  },
})
</script>
