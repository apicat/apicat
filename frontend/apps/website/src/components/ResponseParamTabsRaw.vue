<template>
  <ToggleHeading title="响应参数" v-if="isShow">
    <el-tabs v-model="editableTabsValue">
      <el-tab-pane v-for="item in model.list" :key="item.id" :name="item.id">
        <template #label>
          <el-space>
            <span>{{ item.description }}</span>
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
import { Definition } from './APIEditor/types'
import { HttpDocument } from '@/typings'
import { HTTP_RESPONSE_NODE_KEY, useNodeAttrs } from '@/hooks/useNodeAttrs'
import { APICatResponse } from './ResponseForm.vue'
import { uuid } from '@apicat/shared'

const props = defineProps<{ doc: HttpDocument; definitions: Definition[] }>()
const response = useNodeAttrs(props, HTTP_RESPONSE_NODE_KEY, 'doc')

const model = ref<{ list: APICatResponse[] }>({ list: [] })
const isShow = computed(() => model.value.list.length > 0)
const editableTabsValue = ref()

watch(response, () => {
  model.value.list = response.value.list.map((item: APICatResponse) => ({ ...item, id: item.id || uuid() }))
  editableTabsValue.value = (model.value.list[0] as APICatResponse).id
})
</script>
