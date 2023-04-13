<template>
  <div :class="ns.b()" class="is-edit">
    <input class="ac-document__title" type="text" v-model="httpDoc.title" maxlength="255" ref="title" placeholder="请输入接口标题" />

    <div class="ac-editor mt-10px">
      <RequestMethodEditor class="mb-10px" :model-value="httpDoc" />
      <RequestParamEditor class="mb-10px" v-model="httpDoc" :definitions="definitions" />
      <ResponseEditor v-model="httpDoc" :definitions="definitions" />
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
const ns = useNamespace('document')
const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)

const props = defineProps<{ modelValue: HttpDocument; definitions?: Definition[] }>()
const emit = defineEmits(['update:modelValue'])
const httpDoc = useVModel(props, 'modelValue', emit)
</script>
