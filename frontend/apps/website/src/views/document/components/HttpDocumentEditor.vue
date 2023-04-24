<template>
  <div :class="ns.b()" class="is-edit">
    <input class="ac-document__title" type="text" v-model="httpDoc.title" maxlength="255" ref="title" placeholder="请输入接口标题" />

    <div class="ac-editor mt-10px">
      <RequestMethodEditor class="mb-10px" v-model="httpDoc" />
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
import { HTTP_RESPONSE_NODE_KEY } from './createHttpDocument'
const ns = useNamespace('document')
const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)

const props = defineProps<{ modelValue: HttpDocument; definitions?: Definition[] }>()
const emit = defineEmits(['update:modelValue'])
const httpDoc = useVModel(props, 'modelValue', emit)

const httpResponseList = computed({
  get: () => {
    const response = httpDoc.value.content.find((node) => node.type === HTTP_RESPONSE_NODE_KEY)
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
