<template>
  <div v-if="isShow">
    <h2 class="text-16px font-500">响应参数</h2>
    <el-tabs @tab-add="handleAddTab" @tab-remove="handleRemoveTab" editable v-model="editableTabsValue">
      <el-tab-pane v-for="(item, index) in model" :key="item.id" :name="item.id" :disabled="disabled">
        <template #label>
          <el-space draggable="true" @dragstart="dragStartHandler($event, index)" @dragend="dragEndHandler">
            <span>{{ item.description }}</span>
            <AcTag :style="getResponseStatusCodeBgColor(item.code)">{{ item.code }}</AcTag>
          </el-space>
        </template>
        <ResponseForm v-model="model[index]" :definitions="definitions" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { HttpDocument } from '@/typings'
import { Definition } from './APIEditor/types'
import ResponseForm from './ResponseForm.vue'
import { getResponseStatusCodeBgColor } from '@/commons'
import { useNodeAttrs, HTTP_RESPONSE_NODE_KEY } from '@/hooks/useNodeAttrs'
import { uuid } from '@apicat/shared'
import { createResponseDefaultContent } from '@/views/document/components/createHttpDocument'

const props = defineProps<{ modelValue: HttpDocument; definitions?: Definition[] }>()
const nodeAttrs = useNodeAttrs(props, HTTP_RESPONSE_NODE_KEY)

const model = computed({
  get: () => {
    nodeAttrs.value.list = nodeAttrs.value.list.map((item: any) => ({ ...item, id: item.id || uuid() }))
    return nodeAttrs.value.list
  },
  set: (val: any) => {
    nodeAttrs.value.list = val
  },
})

const editableTabsValue = ref()

const isShow = computed(() => model.value.length > 0)

const disabled = computed(() => model.value.length <= 1)

const createResponse = () => {
  return {
    id: uuid(),
    code: 200,
    description: 'success',
    content: createResponseDefaultContent(),
  }
}

const activeLastTab = () => {
  const len = model.value.length
  const res = model.value[len - 1]
  editableTabsValue.value = res.id
}

const handleAddTab = () => {
  model.value.push(createResponse())
  activeLastTab()
}

const handleRemoveTab = (id: any) => {
  const index = model.value.findIndex((item: any) => item.id === id)
  nodeAttrs.value.list.splice(index, 1)
  if (id === editableTabsValue.value) {
    activeLastTab()
  }
}

const dragDataKey = 'application/apicat-response-tab'
const dragStartHandler = (e: DragEvent, index: number) => {
  e.dataTransfer!.dropEffect = 'move'
  console.log(e)
  const nodeEle = (e.target as Element).parentElement
  nodeEle && nodeEle.classList.add('dragging')
  e.dataTransfer?.setDragImage(nodeEle as HTMLElement, 0, 0)
  e.dataTransfer?.setData(dragDataKey, index + '')
}

const dragEndHandler = (ev: DragEvent) => {
  const nodeEle = (ev.target as Element).parentElement
  nodeEle && nodeEle.classList.remove('dragging')
}

watch(nodeAttrs, () => {
  editableTabsValue.value = model.value[0].id
})
</script>
