<template>
  <div v-loading="isLoading">
    <ToggleHeading title="Header" type="card">
      <SimpleParameterEditor
        :readonly="isReader"
        v-model="parameters.header"
        :draggable="false"
        :on-create="(raw) => onCreateParams(raw, 'header')"
        :on-change="(raw) => onUpdateParams(raw, 'header')"
      >
        <template #operate="{ row, index, delHandler }">
          <el-icon @click="(e) => onClickDeleteIcon(row, index, delHandler, e)" class="cursor-pointer"><ac-icon-ep-delete /></el-icon>
        </template>
      </SimpleParameterEditor>
    </ToggleHeading>

    <ToggleHeading title="Cookie" type="card">
      <SimpleParameterEditor
        :readonly="isReader"
        v-model="parameters.cookie"
        :draggable="false"
        :on-create="(raw) => onCreateParams(raw, 'cookie')"
        :on-change="(raw) => onUpdateParams(raw, 'cookie')"
      >
        <template #operate="{ row, index, delHandler }">
          <el-icon @click="(e) => onClickDeleteIcon(row, index, delHandler, e)" class="cursor-pointer"><ac-icon-ep-delete /></el-icon>
        </template>
      </SimpleParameterEditor>
    </ToggleHeading>

    <ToggleHeading title="Query" type="card">
      <SimpleParameterEditor
        :readonly="isReader"
        v-model="parameters.query"
        :draggable="false"
        :on-create="(raw) => onCreateParams(raw, 'query')"
        :on-change="(raw) => onUpdateParams(raw, 'query')"
      >
        <template #operate="{ row, index, delHandler }">
          <el-icon @click="(e) => onClickDeleteIcon(row, index, delHandler, e)" class="cursor-pointer"><ac-icon-ep-delete /></el-icon>
        </template>
      </SimpleParameterEditor>
    </ToggleHeading>

    <ToggleHeading title="Path" type="card">
      <SimpleParameterEditor
        :readonly="isReader"
        v-model="parameters.path"
        :draggable="false"
        :on-create="(raw) => onCreateParams(raw, 'path')"
        :on-change="(raw) => onUpdateParams(raw, 'path')"
      >
        <template #operate="{ row, index, delHandler }">
          <el-icon @click="(e) => onClickDeleteIcon(row, index, delHandler, e)" class="cursor-pointer"><ac-icon-ep-delete /></el-icon>
        </template>
      </SimpleParameterEditor>
    </ToggleHeading>
  </div>

  <el-popover :visible="isShow" width="auto" :virtual-ref="popoverRefEl" virtual-triggering>
    <div class="ignore-popper">
      <p class="flex-y-center">
        <el-icon class="mr-2px" :color="'rgb(255, 153, 0)'"><ac-icon-ep:question-filled /></el-icon>确认删除此全局参数吗？
      </p>
      <el-checkbox size="small" style="font-weight: normal" v-model="isUnRef" :true-label="1" :false-label="0">对引用此参数的内容解引用</el-checkbox>
      <div class="flex justify-end">
        <el-button size="small" text @click="hidePopover">{{ $t('app.common.cancel') }}</el-button>
        <el-button size="small" type="primary" @click="handleConfirmDelete">{{ $t('app.common.confirm') }}</el-button>
      </div>
    </div>
  </el-popover>
</template>
<script setup lang="ts">
import SimpleParameterEditor from '@/components/APIEditor/SimpleEditor.vue'
import useApi from '@/hooks/useApi'
import uesGlobalParametersStore from '@/store/globalParameters'
import { storeToRefs } from 'pinia'
import { usePopover } from '@/hooks/usePopover'
import useProjectStore from '@/store/project'
import { useParams } from '@/hooks/useParams'

const { project_id } = useParams()
const globalParametersStore = uesGlobalParametersStore()
const { isReader } = storeToRefs(useProjectStore())
const { parameters } = storeToRefs(globalParametersStore)
const [isLoading, getCommonParamListApi] = useApi(globalParametersStore.getGlobalParameters)
let currentDeleteParam: { param: any; index: number; delHandler: any } | null = null
const isUnRef = shallowRef(1)

const { isShow, popoverRefEl, showPopover, hidePopover } = usePopover({
  onHide() {
    currentDeleteParam = null
  },
})

const onClickDeleteIcon = (param: any, index: number, delHandler: any, e: PointerEvent) => {
  isUnRef.value = 1
  currentDeleteParam = { param, index, delHandler }
  showPopover(e.target as HTMLElement)
}

const handleConfirmDelete = async () => {
  if (currentDeleteParam) {
    const { param, index, delHandler } = currentDeleteParam
    hidePopover()
    await nextTick()
    await globalParametersStore.deleteGlobalParameter(project_id, param, isUnRef.value)
    delHandler(index)
  }
}

const onCreateParams = async (raw: any, type: string) => await globalParametersStore.addGlobalParameter(project_id, type, raw)
const onUpdateParams = async (raw: any, type: string) => await globalParametersStore.updateGlobalParameter(project_id, type, raw)

onMounted(async () => await getCommonParamListApi(project_id as string))
</script>
