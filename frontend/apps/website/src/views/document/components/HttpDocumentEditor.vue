<template>
  <div :class="ns.b()" class="is-edit">
    <input class="ac-document__title" type="text" v-model="httpDoc.title" maxlength="255" ref="title" :placeholder="$t('app.interface.form.title')" />

    <div class="ac-editor mt-10px">
      <RequestMethodEditor class="mb-10px" v-model="httpDoc" @url-change="onUrlChange" />
      <RequestParamEditor class="mb-10px" v-model="httpDoc" :definitions="definitions" />
      <ResponseEditor v-model:data="httpResponseList" :definitions="definitions" :definition-responses="responses" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { APICatCommonResponse, HttpDocument } from '@/typings'
import RequestParamEditor from '@/components/RequestParamEditor'
import { useNamespace } from '@/hooks/useNamespace'
import ResponseEditor from '@/components/ResponseEditor.vue'
import useDefinitionStore from '@/store/definition'
import { storeToRefs } from 'pinia'
import { HTTP_RESPONSE_NODE_KEY, HTTP_REQUEST_NODE_KEY } from './createHttpDocument'
import useDefinitionResponseStore from '@/store/definitionResponse'
import { ElMessage } from 'element-plus'
import { isEmpty } from 'lodash-es'
import { useI18n } from 'vue-i18n'

const ns = useNamespace('document')
const { t } = useI18n()
const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)

const definitionResponseStore = useDefinitionResponseStore()
const { responses } = storeToRefs(definitionResponseStore)

const props = defineProps<{ modelValue: HttpDocument }>()
const emit = defineEmits(['update:modelValue'])
const httpDoc = useVModel(props, 'modelValue', emit, { passive: false })

const onUrlChange = (paths: string[]) => {
  const request = httpDoc.value.content.find((node) => node.type === HTTP_REQUEST_NODE_KEY)
  request.attrs.parameters.path = paths.map((name: string) => {
    const param = request.attrs.parameters.path.find((param: any) => param.name === name)
    return {
      name,
      required: true,
      description: '',
      schema: {
        type: 'string',
        examples: '',
        description: '',
      },
      ...param,
    }
  })
}

const httpResponseList = computed(() => {
  const response = httpDoc.value.content.find((node) => node.type === HTTP_RESPONSE_NODE_KEY)
  return response?.attrs?.list || []
})
</script>
