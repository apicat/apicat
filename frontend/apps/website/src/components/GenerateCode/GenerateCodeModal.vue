<template>
  <el-dialog v-model="dialogVisible" append-to-body :close-on-click-modal="false" class="fullscree hide-header" destroy-on-close align-center width="960px">
    <div class="flex overflow-hidden rounded h-600px">
      <div class="flex flex-col py-5 w-240px">
        <div class="px-5 pb-5 text-lg font-medium">{{ $t('app.common.generateModelCode') }}</div>
        <div class="flex-1 px-5 overflow-y-auto">
          <ElForm label-width="auto" label-position="top">
            <el-form-item label="选择语言">
              <ElSelect v-model="currentLanguage" class="w-full">
                <ElOption v-for="item in languages" :key="item.name" :value="item.label">{{ item.label }}</ElOption>
              </ElSelect>
            </el-form-item>
            <el-form-item :label="$t('app.codeGen.model.name')">
              <ElInput v-model="dataModelName" :placeholder="$t('app.codeGen.rules.name')" />
            </el-form-item>
          </ElForm>

          <ElForm label-width="auto" label-position="top" :model="currentLanguageOptions">
            <template v-if="currentLanguageOptionRender" v-for="item in currentLanguageOptionRender.primaryOptions">
              <el-form-item v-if="!item.isBooleanOption" :label="item.description">
                <ElInput v-if="item.isStringOption" v-model="currentLanguageOptions[item.name]" />
                <ElSelect v-if="item.isEnumOption" v-model="currentLanguageOptions[item.name]" class="w-full">
                  <ElOption v-for="value in item.legalValues" :key="value" :value="value">{{ value }}</ElOption>
                </ElSelect>
              </el-form-item>

              <el-form-item size="small" v-if="item.isBooleanOption" label="" style="margin-bottom: 4px">
                <ElSwitch class="" v-if="item.isBooleanOption" v-model="currentLanguageOptions[item.name]"></ElSwitch>
                <div class="flex-1 leading-none ml-4px pt-2px">{{ item.description }}</div>
              </el-form-item>
            </template>

            <template v-if="currentLanguageOptionRender" v-for="item in currentLanguageOptionRender.secondaryOptions">
              <div>
                <el-form-item v-if="!item.isBooleanOption" :label="item.description">
                  <ElInput v-if="item.isStringOption" v-model="currentLanguageOptions[item.name]" />
                  <ElSelect v-if="item.isEnumOption" v-model="currentLanguageOptions[item.name]" class="w-full">
                    <ElOption v-for="value in item.legalValues" :key="value" :value="value">{{ value }}</ElOption>
                  </ElSelect>
                </el-form-item>

                <el-form-item size="small" v-if="item.isBooleanOption" label="" style="margin-bottom: 4px">
                  <ElSwitch class="" v-if="item.isBooleanOption" v-model="currentLanguageOptions[item.name]"></ElSwitch>
                  <div class="flex-1 leading-none ml-4px pt-2px">{{ item.description }}</div>
                </el-form-item>
              </div>
            </template>
          </ElForm>
        </div>
      </div>
      <div class="flex flex-col flex-1 px-5 overflow-hidden pt-16px">
        <div class="text-base font-medium pb-10px">{{ currentLanguage }} {{ $t('app.common.code') }}</div>
        <div class="flex-1 pb-5 overflow-scroll">
          <CodeEditor readonly style="max-height: fit-content" class="h-full" :model-value="code" :lang="currentLanguageForCodeMirror" />
        </div>
      </div>
    </div>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import CodeEditor from '../APIEditor/CodeEditor.vue'
import { useGenerateCode } from './useGenerateCode'
import { DefinitionSchema } from '../APIEditor/types'

const { dialogVisible, showModel } = useModal()
const [isShowMoreMenu, toggleMoreMenu] = useToggle()
const { code, dataModelName, apicatSchema, languages, currentLanguage, currentLanguageOptionRender, currentLanguageOptions, currentLanguageForCodeMirror } = useGenerateCode()

defineExpose({
  show: async (definitionSchema: DefinitionSchema) => {
    showModel()
    dataModelName.value = definitionSchema.name
    apicatSchema.value = definitionSchema
  },
})
</script>
