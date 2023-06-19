<template>
  <el-dialog v-model="dialogVisible" append-to-body :close-on-click-modal="false" class="fullscree hide-header" destroy-on-close center width="70%">
    <div class="flex overflow-hidden rounded h-600px">
      <div class="flex flex-col py-5 w-240px bg-gray-lighter">
        <div class="px-5 pb-5 text-lg font-medium">生成模型代码</div>
        <div class="flex-1 px-5 overflow-y-auto">
          <ElSelect v-model="currentLang" class="w-full">
            <ElOption v-for="item in languages" :key="item.name" :value="item.label">{{ item.label }}</ElOption>
          </ElSelect>
          <ElForm label-width="auto" label-position="top" :model="currentLanguageOptions">
            <template v-if="currentLanguageOptionRender" v-for="item in currentLanguageOptionRender.options">
              <el-form-item :label="item.description">
                <ElInput v-if="item.isStringOption" v-model="currentLanguageOptions[item.name]" />
                <ElSwitch v-if="item.isBooleanOption" v-model="currentLanguageOptions[item.name]" />
                <ElSelect v-if="item.isEnumOption" v-model="currentLanguageOptions[item.name]">
                  <ElOption v-for="value in item.legalValues" :key="value" :value="value">{{ value }}</ElOption>
                </ElSelect>
              </el-form-item>
            </template>
          </ElForm>
        </div>
      </div>
      <div class="flex flex-col flex-1 px-5 pt-16px">
        <div class="text-base font-medium pb-10px">C# 代码</div>
        <el-button class="text-right mb-10px">复制代码</el-button>
        <div class="flex-1 pb-5 overflow-scroll">
          <CodeEditor style="max-height: fit-content" class="h-full" :model-value="code" />
        </div>
      </div>
    </div>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import { quicktype, InputData, JSONSchemaInput, FetchingJSONSchemaStore, JavaScriptPropTypesTargetLanguage, TypeScriptTargetLanguage } from 'quicktype-core'
import jsonschemaData from './jsonschema'
import CodeEditor from '../APIEditor/CodeEditor.vue'
import { useGenerateCode } from './useGenerateCode'

const { dialogVisible, showModel } = useModal()
const { code, languages, currentLang, currentLanguageOptionRender, currentLanguageOptions } = useGenerateCode()

async function quicktypeJSONSchema(targetLanguage: string, typeName: string, jsonSchemaString: string) {
  const schemaInput = new JSONSchemaInput(new FetchingJSONSchemaStore())

  await schemaInput.addSource({ name: typeName, schema: jsonSchemaString })
  const tsLang = new TypeScriptTargetLanguage()
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
