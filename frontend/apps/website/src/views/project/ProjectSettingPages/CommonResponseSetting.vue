<template>
  <div v-loading="isLoading">
    <el-button link type="primary" @click="handleAddParam">
      <el-icon><ac-icon-ep-plus /></el-icon>添加
    </el-button>

    <div v-for="(param, index) in responseParamList" class="mt-15px" :key="param.id">
      <ToggleHeading
        :title="`${param.detail?.description ?? param.description}(${param.detail?.code ?? param.code})`"
        type="card"
        :expand="param.expand"
        @on-expand="(isExpand:boolean)=>handleExpand(isExpand,param)"
        v-loading="param.isLoading"
      >
        <template #extra>
          <el-popconfirm width="auto" :title="$t('app.table.deleteResponseConfirm')" @confirm="handleDeleteParam(param, index)">
            <template #reference>
              <el-icon class="cursor-pointer"><ac-icon-ep-delete /></el-icon>
            </template>
          </el-popconfirm>
        </template>

        <div v-if="param.detail">
          <ResponseForm v-model="param.detail" :definitions="definitions" class="mt-10px" :is-common-response="true" />
          <el-button class="mt-20px" type="primary" @click="handleSubmit(param)">{{ $t('app.common.save') }}</el-button>
        </div>
      </ToggleHeading>
    </div>

    <el-empty v-if="!responseParamList.length" :image-size="200" />
  </div>
</template>
<script setup lang="ts">
import { useProjectId } from '@/hooks/useProjectId'
import { useResponseParamDetail } from '../logic/useResponseParamDetail'
import { useResponseparamList } from '../logic/useResponseparamList'
import useDefinitionStore from '@/store/definition'
import { storeToRefs } from 'pinia'

const project_id = useProjectId()
const definitionStore = useDefinitionStore()
const { definitions } = storeToRefs(definitionStore)

const { isLoading, responseParamList, handleAddParam, handleDeleteParam } = useResponseparamList({ id: project_id })
const { handleExpand, handleSubmit } = useResponseParamDetail({ id: project_id })
</script>
