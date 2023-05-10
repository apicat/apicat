<template>
  <div :class="ns.b()" class="is-edit">
    <input class="ac-document__title" type="text" v-model="httpDoc.title" maxlength="255" ref="title" :placeholder="$t('app.interface.form.title')" />

    <div class="ac-editor mt-10px">
      <RequestMethodEditor class="mb-10px" v-model="httpDoc" @url-change="onUrlChange" />
      <RequestParamEditor class="mb-10px" v-model="httpDoc" :definitions="definitions" />
      <ResponseEditor v-model:data="httpResponseList" :definitions="definitions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { HttpDocument } from '@/typings'
import RequestParamEditor from '@/components/RequestParamEditor'
import { useNamespace } from '@/hooks/useNamespace'
import ResponseEditor from '@/components/ResponseEditor.vue'
import { Definition } from '@/components/APIEditor/types'
import useDefinitionStore from '@/store/definition'
import { storeToRefs } from 'pinia'
import { HTTP_RESPONSE_NODE_KEY, HTTP_REQUEST_NODE_KEY } from './createHttpDocument'
const ns = useNamespace('document')
const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)

const props = defineProps<{ modelValue: HttpDocument; definitions?: Definition[] }>()
const emit = defineEmits(['update:modelValue'])
const httpDoc = useVModel(props, 'modelValue', emit)

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

const httpResponseList = computed({
  get: () => {
    const response = httpDoc.value.content.find((node) => node.type === HTTP_RESPONSE_NODE_KEY)
    console.log(response)
    return response?.attrs?.list || []
  },
  set: (val) => {
    let response = httpDoc.value.content.find((node) => node.type === HTTP_RESPONSE_NODE_KEY)
    if (!response) {
      response = { attrs: { list: [] } }
    }
    response.attrs.list = val
  },
})
</script>
