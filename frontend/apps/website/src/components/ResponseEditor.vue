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
      <template v-for="(item, index) in model" :key="item._id">
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
          <ResponseForm v-if="!item.$ref" v-model="model[index]" :definitions="definitions" />
          <ResponseRefForm v-else v-model:response="model[index]" :definition-responses="definitionResponses" :definition-schemas="definitions" />
        </el-tab-pane>
      </template>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { DefinitionSchema } from './APIEditor/types'
import ResponseForm from './ResponseForm.vue'
import ResponseRefForm from './ResponseRefForm.vue'
import { RefPrefixKeys, getResponseStatusCodeBgColor, markDataWithKey } from '@/commons'
import { createDefaultResponseContent } from '@/views/document/components/createDefaultDefinition'
import { useDragAndDrop } from '@/hooks/useDragAndDrop'
import { useI18n } from 'vue-i18n'
import SelectDefinitionResponse from './DefinitionResponse/SelectDefinitionResponse.vue'
import { DefinitionResponse } from '@/typings'

const props = defineProps<{ data: Array<any>; definitions?: DefinitionSchema[]; definitionResponses?: DefinitionResponse[] }>()

const createResponse = (item?: any) => {
  return {
    name: 'Response Name',
    code: 200,
    description: '',
    content: createDefaultResponseContent(),
    ...item,
  }
}

const { definitionResponses } = toRefs(props)

const model = useVModel(props, 'data', undefined, { passive: true })
const editableTabsValue = ref('')
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
  const res = { code: 200, $ref: `${RefPrefixKeys.DefinitionResponse.key}${response.id}` }
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

watch(
  [model, definitionResponses],
  ([data, responses]: any) => {
    if (!data.length || !responses.length) {
      return
    }

    for (const item of data) {
      if (item.$ref) {
        const id = item.$ref.match(RefPrefixKeys.DefinitionResponse.reg)?.[1]
        const resId = parseInt(id as string, 10)
        const response: any = responses.find((item: any) => item.id === resId)
        response && markDataWithKey(item, 'name', response.name)
      }

      markDataWithKey(item)
    }

    !editableTabsValue.value && activeLastTab(data[0]._id)
  },
  { immediate: true, deep: true }
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
