<template>
  <div v-if="isShow" class="ac-response-editor">
    <h2 class="text-16px font-500">响应参数</h2>
    <el-tabs @tab-add="handleAddTab" @tab-remove="handleRemoveTab" editable v-model="editableTabsValue">
      <el-tab-pane v-for="(item, index) in model" :key="item.id + index" :name="item.id" :disabled="disabled">
        <template #label>
          <el-space
            draggable="true"
            @dragstart="onDragStart($event, index)"
            @dragend="onDragEnd"
            @dragover="onDragOver($event, index)"
            @dragleave="onDragLeave($event, index)"
            @drop="onDropHandler($event, index)"
          >
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
import { useDragAndDrop } from '@/hooks/useDragAndDrop'

const props = defineProps<{ modelValue: HttpDocument; definitions?: Definition[] }>()
const nodeAttrs = useNodeAttrs(props, HTTP_RESPONSE_NODE_KEY)

const { onDragStart, onDragOver, onDragLeave, onDragEnd, onDropHandler } = useDragAndDrop({
  onDrop: (dragIndex: number, dropIndex: number) => {
    const dragItem = nodeAttrs.value.list[dragIndex]
    nodeAttrs.value.list.splice(dragIndex, 1)
    nodeAttrs.value.list.splice(dropIndex, 0, dragItem)
  },
})

const model = computed(() => {
  nodeAttrs.value.list = nodeAttrs.value.list.map((item: any) => ({ ...item, id: item.id || uuid() }))
  return nodeAttrs.value.list
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

watch(nodeAttrs, () => {
  editableTabsValue.value = model.value[0].id
})
</script>
<style lang="scss">
.ac-response-editor {
  .el-tabs__item .is-icon-close {
    margin-top: 13px;
    padding-bottom: 1px;
  }
}
</style>
