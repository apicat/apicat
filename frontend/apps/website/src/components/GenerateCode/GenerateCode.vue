<script setup lang="ts">
import { CodeMirror } from '@apicat/components'
import { languages as langs } from '@codemirror/language-data'
import { useGenerateCode } from './useGenerateCode'

const props = defineProps({
  schema: {
    type: Object as PropType<Definition.Schema>,
    required: true,
  },
})

const {
  code,
  dataModelName,
  apicatSchema,
  languages,
  currentLanguage,
  currentLanguageOptionRender,
  currentLanguageOptions,
  currentLanguageForCodeMirror,

  loading,
} = useGenerateCode()

const isShowModelName = computed(() => {
  return !['JSON', 'JSONSchema'].includes(currentLanguage.value)
})

watch(
  () => props.schema,
  (val) => {
    if (val) {
      dataModelName.value = val.name
      apicatSchema.value = val
    }
  },
  { immediate: true },
)

defineExpose({
  dispose() {
    loading.value = true
    code.value = ''
  },
})
</script>

<template>
  <div class="flex overflow-hidden rounded">
    <div class="flex-1 h-auto overflow-hidden max-h-664px py-16px">
      <CodeMirror
        v-loading="loading"
        background="transparent"
        readonly
        class="h-full scroll-content"
        :config-extensions="{ isShowLineNumber: false }"
        :model-value="code"
        :lang="currentLanguageForCodeMirror"
        :languages="langs" />
    </div>
    <div class="flex flex-col py-16px pl-16px w-240px">
      <div class="flex-1 overflow-y-auto">
        <ElForm label-width="auto" label-position="top">
          <el-form-item :label="$t('app.codeGen.tips.chooseLanguage')">
            <ElSelect v-model="currentLanguage" class="w-full">
              <ElOption v-for="item in languages" :key="item.name" :value="item.label">
                {{ item.label }}
              </ElOption>
            </ElSelect>
          </el-form-item>
          <el-form-item v-show="isShowModelName" :label="$t('app.codeGen.model.name')">
            <ElInput v-model="dataModelName" :placeholder="$t('app.codeGen.rules.name')" />
          </el-form-item>
        </ElForm>

        <ElForm label-width="auto" label-position="top" :model="currentLanguageOptions">
          <template v-if="currentLanguageOptionRender">
            <template v-for="item in currentLanguageOptionRender.primaryOptions" :key="item.name">
              <el-form-item v-if="!item.isBooleanOption" :label="item.description">
                <ElInput v-if="item.isStringOption" v-model="currentLanguageOptions[item.name]" />
                <ElSelect v-if="item.isEnumOption" v-model="currentLanguageOptions[item.name]" class="w-full">
                  <ElOption v-for="value in item.legalValues" :key="value" :value="value">
                    {{ value }}
                  </ElOption>
                </ElSelect>
              </el-form-item>

              <el-form-item v-if="item.isBooleanOption" size="small" label="" style="margin-bottom: 4px">
                <el-checkbox v-if="item.isBooleanOption" v-model="currentLanguageOptions[item.name]" />
                <div
                  class="flex-1 leading-none cursor-pointer ml-4px pt-2px"
                  @click="currentLanguageOptions[item.name] = !currentLanguageOptions[item.name]">
                  {{ item.description }}
                </div>
              </el-form-item>
            </template>
          </template>
          <template v-if="currentLanguageOptionRender">
            <template v-for="item in currentLanguageOptionRender.secondaryOptions" :key="item.name">
              <div>
                <el-form-item v-if="!item.isBooleanOption" :label="item.description">
                  <ElInput v-if="item.isStringOption" v-model="currentLanguageOptions[item.name]" />
                  <ElSelect v-if="item.isEnumOption" v-model="currentLanguageOptions[item.name]" class="w-full">
                    <ElOption v-for="value in item.legalValues" :key="value" :value="value">
                      {{ value }}
                    </ElOption>
                  </ElSelect>
                </el-form-item>

                <el-form-item v-if="item.isBooleanOption" size="small" label="" style="margin-bottom: 4px">
                  <el-checkbox v-if="item.isBooleanOption" v-model="currentLanguageOptions[item.name]" />
                  <div
                    class="flex-1 leading-none cursor-pointer ml-4px pt-2px"
                    @click="currentLanguageOptions[item.name] = !currentLanguageOptions[item.name]">
                    {{ item.description }}
                  </div>
                </el-form-item>
              </div>
            </template>
          </template>
        </ElForm>
      </div>
    </div>
  </div>
</template>
