<template>
  <div v-if="isShow" class="ac-response-editor">
    <h2 class="relative flex justify-between text-16px font-500">
      {{ $t('app.response.title') }}
      <div class="absolute right-0 z-10 bg-white -bottom-30px">
        <el-popover :width="250" trigger="hover" class="">
          <template #reference>
            <el-button link type="primary" @click="handleAddTab">
              <el-icon><ac-icon-ep-plus /></el-icon>{{ $t('app.common.add') }}
            </el-button>
          </template>
          <div class="clear-popover-space">
            <p class="border-b cursor-pointer border-gray-lighter px-10px h-44px flex-y-center hover:bg-gray-100">
              <el-icon class="mr-4px"><ac-icon-ep-plus /></el-icon>
              新建响应
            </p>
            <SelectDefinitionResponse :responses="definitionResponses" @select="handleAddRefResponse" />
          </div>
        </el-popover>
      </div>
    </h2>

    <el-tabs @tab-remove="handleRemoveTab" editable v-model="editableTabsValue">
      <template v-for="(item, index) in localModel" :key="item._id + index">
        <el-tab-pane :name="item._id" :disabled="disabled">
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
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { DefinitionSchema } from './APIEditor/types'
import ResponseForm from './ResponseForm.vue'
import { RefPrefixKeys, getResponseStatusCodeBgColor, markDataWithKey } from '@/commons'
import { uuid } from '@apicat/shared'
import { createDefaultResponseContent } from '@/views/document/components/createDefaultDefinition'
import { useDragAndDrop } from '@/hooks/useDragAndDrop'
import { ElMessage } from 'element-plus'
import { isEmpty, debounce, cloneDeep } from 'lodash-es'
import { useI18n } from 'vue-i18n'
import SelectDefinitionResponse from './DefinitionResponse/SelectDefinitionResponse.vue'
import { DefinitionResponse } from '@/typings'

const emits = defineEmits(['update:data'])
const props = defineProps<{ data: Array<any>; definitions?: DefinitionSchema[]; definitionResponses?: DefinitionResponse[] }>()
const { t } = useI18n()

const createResponse = (item?: any) => {
  return {
    name: 'Response Name',
    code: 200,
    description: '',
    content: createDefaultResponseContent(),
    ...item,
  }
}

const { data, definitionResponses } = toRefs(props)

const model: any = ref([])

watch(
  data,
  () => {
    model.value = cloneDeep(data.value)
  },
  {
    immediate: true,
  }
)

const localModel = computed(() => {
  return model.value.map((item: any) => {
    if (item.$ref) {
      const id = item.$ref.match(RefPrefixKeys.DefinitionResponse.reg)?.[1]
      const resId = parseInt(id, 10)
      const response = definitionResponses?.value?.find((response: DefinitionResponse) => response.id === resId)
      if (response) {
        item = { ...item, ...cloneDeep(response), id: undefined }
      }
    }
    markDataWithKey(item)
    return item
  })
})

const editableTabsValue = ref(unref(model).length ? unref(model)[0]._id : null)
const isShow = computed(() => model.value.length > 0)
const disabled = computed(() => model.value.filter((item: any) => !item._isCommonResponse).length <= 1)

const activeLastTab = (_id?: string) => {
  const len = model.value.length
  const res = model.value[len - 1]
  editableTabsValue.value = _id || res._id
}

const handleAddTab = () => {
  const res = createResponse()
  markDataWithKey(res)
  model.value.push(res)
  activeLastTab()
}

const handleAddRefResponse = (response: DefinitionResponse) => {
  const res = { ...response, code: 200, $ref: `${RefPrefixKeys.DefinitionResponse.key}${response.id}` }
  markDataWithKey(res)
  model.value.push(res)
  activeLastTab()
}

const handleRemoveTab = (_id: any) => {
  const index = model.value.findIndex((item: any) => item._id === _id)
  model.value.splice(index, 1)
  if (_id === editableTabsValue.value) {
    activeLastTab()
  }
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
      // emits('update:data', [...model.value])
    }
  }, 300),
  { deep: true }
)
</script>
<style lang="scss">
.ac-response-editor {
  .el-tabs__item .is-icon-close {
    vertical-align: -1px;
    display: inline-block;
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
