<template>
  <ToggleHeading :title="$t('app.response.title')" v-if="isShow">
    <el-tabs v-model="editableTabsValue">
      <el-tab-pane v-for="item in model" :key="item.id" :name="item.id">
        <template #label>
          <el-space>
            <span>{{ item.name || item.description }}</span>
            <AcTag :style="getResponseStatusCodeBgColor(item.code)">{{ item.code }}</AcTag>
          </el-space>
        </template>
        <ResponseParamPaneRaw :param="item" :definitions="definitions" />
      </el-tab-pane>
    </el-tabs>
  </ToggleHeading>
</template>

<script setup lang="ts">
import { getResponseStatusCodeBgColor } from '@/commons'
import ResponseParamPaneRaw from './ResponseParamPaneRaw.vue'
import { DefinitionSchema } from './APIEditor/types'
import { APICatCommonResponse, HttpDocument } from '@/typings'
import { HTTP_RESPONSE_NODE_KEY, useNodeAttrs } from '@/hooks/useNodeAttrs'
import { APICatResponse } from './ResponseForm.vue'
import { uuid } from '@apicat/shared'
import useDefinitionResponseStore from '@/store/definitionResponse'
import { storeToRefs } from 'pinia'

const props = defineProps<{ doc: HttpDocument; definitions: DefinitionSchema[] }>()
const responseNode = useNodeAttrs(props, HTTP_RESPONSE_NODE_KEY, 'doc')
const definitionResponseStore = useDefinitionResponseStore()

const editableTabsValue = ref()
const { responses } = storeToRefs(definitionResponseStore)

const model = computed(() => {
  const list = responseNode.value.list.map((item: APICatResponse & APICatCommonResponse) => {
    let newItem = { ...item, id: item.id || uuid() }

    // common response
    if (newItem.$ref) {
      const responseId = parseInt(newItem.$ref.split('/').pop() as string, 10)
      const responseDetail = responses.value.find((item) => item.id === responseId)
      newItem = { ...newItem, ...responseDetail, id: newItem.id }
    }
    return newItem
  })

  editableTabsValue.value = (list[0] as APICatResponse).id
  return list
})

const isShow = computed(() => model.value.length > 0)
</script>
