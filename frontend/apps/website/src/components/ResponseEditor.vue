<template>
  <div v-if="isShow" class="ac-response-editor">
    <h2 class="relative flex justify-between text-16px font-500">
      {{ $t('app.response.title') }}
      <div class="absolute right-0 z-10 bg-white -bottom-30px">
        <el-button link type="primary" @click="handleAddTab">
          <el-icon><ac-icon-ep-plus /></el-icon>{{ $t('app.common.add') }}
        </el-button>
      </div>
    </h2>
    <el-tabs @tab-remove="handleRemoveTab" editable v-model="editableTabsValue">
      <template v-for="(item, index) in model" :key="item._id + index">
        <el-tab-pane v-if="!item._isCommonResponse" :name="item._id" :disabled="disabled">
          <template #label>
            <div
              class="inline-flex items-center"
              draggable="true"
              @dragstart="onDragStart($event, index)"
              @dragend="onDragEnd"
              @dragover="onDragOver($event, index)"
              @dragleave="onDragLeave($event, index)"
              @drop="onDropHandler($event, index)"
            >
              <span class="mr-4px">{{ item.name || '&nbsp' }}</span>
              <AcTag :style="getResponseStatusCodeBgColor(item.code)">{{ item.code }}</AcTag>
            </div>
          </template>
          <ResponseForm v-model="model[index]" :definitions="definitions" />
        </el-tab-pane>
      </template>
      <el-tab-pane name="add-tab" disabled class="ac-response__common">
        <template #label>
          <el-space @click="onShowCommonResponseModal">
            <span>{{ $t('app.publicResponse.title') }}</span>
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
import { ElMessage } from 'element-plus'
import { isEmpty, debounce } from 'lodash-es'
import { useI18n } from 'vue-i18n'

const emits = defineEmits(['update:data'])
const props = defineProps<{ data: Array<any>; definitions?: Definition[] }>()
const { t } = useI18n()

const createResponse = (item?: any) => {
  return {
    _id: uuid(),
    _isCommonResponse: false,
    name: 'Response Name',
    code: 200,
    description: 'success',
    content: createResponseDefaultContent(),
    ...item,
  }
}

const createCommonRefResponse = (name: string) => ({ $ref: `${RefPrefixKeys.CommonResponse.key}${name}`, _id: uuid(), _isCommonResponse: true, _refName: name })

const model: any = ref(
  (props.data || []).map((item: any) => {
    const newItem = { ...item, _id: uuid(), name: item.name || 'Response Name' }
    newItem._isCommonResponse = false
    if (newItem.$ref && newItem.$ref.startsWith(RefPrefixKeys.CommonResponse.key)) {
      newItem._isCommonResponse = true
      newItem._refName = newItem.$ref.match(RefPrefixKeys.CommonResponse.reg)?.[1]
    }
    return newItem
  })
)

const commonResponseCount = computed(() => model.value.filter((item: any) => item._isCommonResponse).length)

const selectCommonResponseModalRef = ref<InstanceType<typeof SelectCommonResponseModal>>()
const editableTabsValue = ref(unref(model).length ? unref(model)[0]._id : null)
const isShow = computed(() => model.value.length > 0)
const disabled = computed(() => model.value.filter((item: any) => !item._isCommonResponse).length <= 1)

const activeLastTab = (_id?: string) => {
  const len = model.value.length
  const res = model.value[len - 1]
  editableTabsValue.value = _id || res._id
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
  const names = model.value.filter((item: any) => item._isCommonResponse).map((item: any) => parseInt(item._refName, 10))
  selectCommonResponseModalRef.value?.show(names)
}

const handleCommonResponseSelectFinish = (selectedIds: string[]) => {
  model.value = model.value.filter((item: any) => !item._isCommonResponse)
  selectedIds.forEach((id) => {
    model.value.push(createCommonRefResponse(id))
  })
}

const { onDragStart, onDragOver, onDragLeave, onDragEnd, onDropHandler } = useDragAndDrop({
  onDrop: (dragIndex: number, dropIndex: number, offset: number) => {
    const dropItem = model.value[dropIndex]
    const dragItemArr = model.value.splice(dragIndex, 1)
    if (dragItemArr.length) {
      let i = model.value.indexOf(dropItem)
      if (offset < 0) i += 1
      model.value.splice(i < 0 ? 0 : i, 0, dragItemArr[0])
    }
  },
})

const validResponseName = () => {
  let len = model.value.length
  for (let i = 0; i < len; i++) {
    const item = model.value[i]
    if (isEmpty(item.name) && !item._isCommonResponse) {
      ElMessage.error(t('app.response.rules.name'))
      // activeLastTab(model.value[i]._id)
      return false
    }
  }
  return true
}

// v-model
watch(
  model,
  debounce(() => {
    if (validResponseName()) {
      const normalResponse: any = []
      const commonResponse: any = []

      model.value.forEach(({ _id, _isCommonResponse, _refName, ...other }: any) => {
        if (_isCommonResponse) {
          commonResponse.push(toRaw(other))
        } else {
          normalResponse.push(toRaw(other))
        }
      })
      emits('update:data', [...normalResponse, ...commonResponse])
    }
  }, 300),
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
