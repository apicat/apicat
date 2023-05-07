<template>
  <div v-loading="isLoading">
    <el-button link type="primary" @click="handleAddParam">
      <el-icon><ac-icon-ep-plus /></el-icon>{{ $t('app.common.add') }}
    </el-button>

    <div v-for="(param, index) in responseParamList" class="mt-15px" :key="param.id">
      <ToggleHeading
        :title="`${param.detail?.name ?? param.name}(${param.detail?.code ?? param.code})`"
        type="card"
        :expand="param.expand"
        @on-expand="(isExpand:boolean)=>handleExpand(isExpand,param)"
        v-loading="param.isLoading"
      >
        <template #extra>
          <!-- <el-popconfirm width="auto" :title="$t('app.response.tips.confirmDelete')" @confirm="handleDeleteParam(param, index)">
            <template #reference> -->
          <el-icon @click="(e) => onClickDeleteIcon(param, index, e)" class="cursor-pointer"><ac-icon-ep-delete /></el-icon>
          <!-- </template>
          </el-popconfirm> -->
        </template>

        <div v-if="param.detail">
          <ResponseForm v-model="param.detail" :definitions="definitions" class="mt-10px" :is-common-response="true" />
          <el-button class="mt-20px" type="primary" @click="handleSubmit(param)">{{ $t('app.common.save') }}</el-button>
        </div>
      </ToggleHeading>
    </div>

    <el-empty v-if="!responseParamList.length" :image-size="200" />
  </div>

  <el-popover :visible="isShow" width="auto" :virtual-ref="popoverRefEl" virtual-triggering>
    <div class="ignore-popper">
      <p class="flex-y-center">
        <el-icon class="mr-2px" :color="'rgb(255, 153, 0)'"><ac-icon-ep:question-filled /></el-icon>确认删除此公共响应吗？
      </p>
      <el-checkbox size="small" style="font-weight: normal" v-model="isUnRef" :true-label="1" :false-label="0">对引用此响应的内容解引用</el-checkbox>
      <div class="flex justify-end">
        <el-button size="small" text @click="hidePopover">{{ $t('app.common.cancel') }}</el-button>
        <el-button size="small" type="primary" :loading="isLoadingForDelete" @click="handelConfirmDelete">{{ $t('app.common.confirm') }}</el-button>
      </div>
    </div>
  </el-popover>
</template>
<script setup lang="ts">
import { useProjectId } from '@/hooks/useProjectId'
import { useResponseParamDetail } from '../logic/useResponseParamDetail'
import { useResponseparamList } from '../logic/useResponseparamList'
import useDefinitionStore from '@/store/definition'
import { storeToRefs } from 'pinia'
import { usePopover } from '@/hooks/usePopover'

const project_id = useProjectId()
const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)
let currentDeleteParam: { param: any; index: number } | null = null
const isUnRef = shallowRef(1)

const { isShow, popoverRefEl, showPopover, hidePopover } = usePopover({
  onHide() {
    currentDeleteParam = null
  },
})

const onClickDeleteIcon = (param: any, index: number, e: PointerEvent) => {
  isUnRef.value = 1
  currentDeleteParam = { param, index }
  showPopover(e.target as HTMLElement)
}

const handelConfirmDelete = async () => {
  if (currentDeleteParam) {
    const { param, index } = currentDeleteParam
    hidePopover()
    await nextTick()
    await handleDeleteParam(param, index, unref(isUnRef))
  }
}

const { isLoading, isLoadingForDelete, responseParamList, handleAddParam, handleDeleteParam } = useResponseparamList({ id: project_id })
const { handleExpand, handleSubmit } = useResponseParamDetail({ id: project_id })
</script>
