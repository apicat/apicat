<template>
  <div v-if="isShow" class="ac-response-editor">
    <h2 class="relative flex justify-between text-16px font-500">
      响应参数
      <div class="absolute right-0 z-10 bg-white -bottom-30px">
        <el-button link type="primary" @click="handleAddTab">
          <el-icon><ac-icon-ep-plus /></el-icon>添加
        </el-button>
      </div>
    </h2>
    <el-tabs @tab-remove="handleRemoveTab" editable v-model="editableTabsValue">
      <el-tab-pane v-for="(item, index) in responseList" :key="item._id + index" :name="item._id" :disabled="disabled">
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
      <el-tab-pane name="new-tab" disabled class="ac-response__common">
        <template #label>
          <el-space @click="onShowCommonResponseModal">
            <span>公共响应</span>
            <span class="inline-block leading-none bg-gray-200 rounded px-4px py-2px">{{ commonResponseCount }}</span>
          </el-space>
        </template>
      </el-tab-pane>
    </el-tabs>
  </div>

  <SelectCommonResponseModal ref="selectCommonResponseModalRef" @ok="handleCommonResponseSelectFinish" />
</template>

<script setup lang="ts">
import { Definition } from './APIEditor/types'
import ResponseForm from './ResponseForm.vue'
import { RefPrefixKeys, getResponseStatusCodeBgColor } from '@/commons'
import { uuid } from '@apicat/shared'
import { createResponseDefaultContent } from '@/views/document/components/createHttpDocument'
import { useDragAndDrop } from '@/hooks/useDragAndDrop'
import SelectCommonResponseModal from '@/views/document/components/SelectCommonResponseModal.vue'

const emits = defineEmits(['update:data'])
const props = defineProps<{ data: Array<any>; definitions?: Definition[] }>()

const createResponse = (item?: any) => {
  return {
    _id: uuid(),
    code: 200,
    description: 'success',
    content: createResponseDefaultContent(),
    ...item,
  }
}

const createCommonRefResponse = (name: string) => ({ $ref: `${RefPrefixKeys.CommonResponse.key}${name}`, _id: uuid(), _isCommonResponse: true, _refName: name })

const model = ref(
  (props.data || []).map((item: any) => {
    const newItem = { ...item, _id: uuid() }

    if (newItem.$ref && newItem.$ref.startsWith(RefPrefixKeys.CommonResponse.key)) {
      newItem._isCommonResponse = true
      newItem._refName = newItem.$ref.match(RefPrefixKeys.CommonResponse.reg)?.[1]
    }
    return newItem
  })
)

const responseList = computed(() => model.value.filter((item) => !item._isCommonResponse))
const commonResponseCount = computed(() => model.value.filter((item) => item._isCommonResponse).length)

const selectCommonResponseModalRef = ref<InstanceType<typeof SelectCommonResponseModal>>()
const editableTabsValue = ref(unref(model).length ? unref(model)[0]._id : null)
const isShow = computed(() => model.value.length > 0)
const disabled = computed(() => model.value.length <= 1)

const activeLastTab = () => {
  const len = model.value.length
  const res = model.value[len - 1]
  editableTabsValue.value = res._id
}

const handleAddTab = () => {
  model.value.push(createResponse())
  activeLastTab()
}

const handleRemoveTab = (_id: any) => {
  const index = model.value.findIndex((item: any) => item._id === _id)
  model.value.splice(index, 1)
  if (_id === editableTabsValue.value) {
    activeLastTab()
  }
}

const onShowCommonResponseModal = () => {
  const names = model.value.filter((item) => item._isCommonResponse).map((item) => item._refName)
  selectCommonResponseModalRef.value?.show(names)
}

const handleCommonResponseSelectFinish = (selectedNames: string[]) => {
  model.value = model.value.filter((item) => !item._isCommonResponse)
  selectedNames.forEach((name) => {
    model.value.push(createCommonRefResponse(name))
  })
}

const { onDragStart, onDragOver, onDragLeave, onDragEnd, onDropHandler } = useDragAndDrop({
  onDrop: (dragIndex: number, dropIndex: number) => {
    const dragItem = model.value[dragIndex]
    model.value.splice(dragIndex, 1)
    model.value.splice(dropIndex, 0, dragItem)
  },
})

// v-model
watch(
  model,
  () => {
    emits(
      'update:data',
      model.value.map(({ _id, _isCommonResponse, _refName, ...other }) => toRaw(other))
    )
  },
  { deep: true }
)
</script>
<style lang="scss">
.ac-response-editor {
  .el-tabs__item .is-icon-close {
    margin-top: 13px;
    padding-bottom: 1px;
  }

  .el-tabs__new-tab {
    width: 40px;
  }

  .el-tabs--top .el-tabs__item.is-top:last-child {
    color: inherit;
    cursor: pointer;
  }
}
</style>
