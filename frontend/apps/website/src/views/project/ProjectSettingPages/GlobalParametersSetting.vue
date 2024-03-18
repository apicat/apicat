<script setup lang="ts">
import { SimpleParamTable, ToggleHeading } from '@apicat/components'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import useApi from '@/hooks/useApi'
import { useGlobalParameters } from '@/store/globalParameter'
import { useParams } from '@/hooks/useParams'
import useProjectStore from '@/store/project'
import { usePopover } from '@/hooks/usePopover'

const { t } = useI18n()
const { projectID } = useParams()
const { isReader } = storeToRefs(useProjectStore())
const globalParametersStore = useGlobalParameters()
const { headers, queries, cookies } = storeToRefs(globalParametersStore)
const [isLoading, getGlobalParameterListApi] = useApi(globalParametersStore.getGlobalParameterList)
let currentDeleteParam: any = null

const isUnRef = ref(false)

const { isShow, popoverRefEl, showPopover, hidePopover } = usePopover()

async function handleCreateParameter(data: any, key: ProjectAPI.GlobalParameterType) {
  data.in = key
  const parameter = await globalParametersStore.addGlobalParameter(projectID, data)
  data.id = parameter.id
}

async function handleUpdateParameter(data: any) {
  await globalParametersStore.updateGlobalParameter(projectID, data)
}

function onClickDeleteIcon(e: PointerEvent, data: any) {
  isUnRef.value = false
  currentDeleteParam = data
  showPopover(e.target as HTMLElement)
}

async function handleSortParameter(data: any, key: ProjectAPI.GlobalParameterType) {
  data.in = key
  await globalParametersStore.sortGlobalParameter(projectID, data)
}

async function handleConfirmDelete() {
  if (currentDeleteParam) {
    hidePopover()
    await globalParametersStore.deleteGlobalParameter(projectID, currentDeleteParam, isUnRef.value)
  }
}

onMounted(async () => await getGlobalParameterListApi(projectID))
</script>

<template>
  <div v-loading="isLoading">
    <ToggleHeading title="Header" type="card">
      <SimpleParamTable
        :datas="headers"
        :readonly="isReader"
        @create="(data) => handleCreateParameter(data, 'header')"
        @sort="(data) => handleSortParameter(data, 'header')"
        @update="handleUpdateParameter"
      >
        <template #operate="{ data }">
          <el-icon class="cursor-pointer" @click="onClickDeleteIcon($event, data)">
            <ac-icon-ep-delete />
          </el-icon>
        </template>
      </SimpleParamTable>
    </ToggleHeading>

    <ToggleHeading title="Cookie" type="card">
      <SimpleParamTable
        :datas="cookies"
        :readonly="isReader"
        @create="(data) => handleCreateParameter(data, 'cookie')"
        @sort="(data) => handleSortParameter(data, 'cookie')"
        @update="handleUpdateParameter"
      >
        <template #operate="{ data }">
          <el-icon class="cursor-pointer" @click="onClickDeleteIcon($event, data)">
            <ac-icon-ep-delete />
          </el-icon>
        </template>
      </SimpleParamTable>
    </ToggleHeading>

    <ToggleHeading title="Query" type="card">
      <SimpleParamTable
        :datas="queries"
        :readonly="isReader"
        @create="(data) => handleCreateParameter(data, 'query')"
        @sort="(data) => handleSortParameter(data, 'query')"
        @update="handleUpdateParameter"
      >
        <template #operate="{ data }">
          <el-icon class="cursor-pointer" @click="onClickDeleteIcon($event, data)">
            <ac-icon-ep-delete />
          </el-icon>
        </template>
      </SimpleParamTable>
    </ToggleHeading>
  </div>

  <el-popover :visible="isShow" width="auto" :virtual-ref="popoverRefEl" virtual-triggering :show-arrow="false">
    <div class="ignore-popper">
      <p class="flex-y-center">
        <el-icon class="mr-5px" color="rgb(255, 153, 0)">
          <ac-icon-ep:question-filled />
        </el-icon>{{ t('app.project.setting.globalParam.tips.delete') }}
      </p>
      <el-checkbox v-model="isUnRef" size="small" style="font-weight: normal">
        <p class="break-words">
          {{ t('app.project.setting.globalParam.tips.unref') }}
        </p>
      </el-checkbox>
      <div class="flex justify-end">
        <el-button size="small" text @click="hidePopover">
          {{ $t('app.common.cancel') }}
        </el-button>
        <el-button size="small" type="primary" @click="handleConfirmDelete">
          {{ $t('app.common.delete') }}
        </el-button>
      </div>
    </div>
  </el-popover>
</template>
