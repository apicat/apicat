<template>
  <div v-loading="isLoading">
    <ToggleHeading title="Header" type="card" :expand="false">
      <SimpleParameterEditor
        v-model="parameters.header"
        :draggable="false"
        :on-create="(raw) => onCreateParams(raw, 'header')"
        :on-change="(raw) => onUpdateParams(raw, 'header')"
        :on-delete="(raw) => onDeleteParams(raw, 'header')"
      />
    </ToggleHeading>

    <ToggleHeading title="Cookie" type="card" :expand="false">
      <SimpleParameterEditor
        v-model="parameters.cookie"
        :draggable="false"
        :on-create="(raw) => onCreateParams(raw, 'cookie')"
        :on-change="(raw) => onUpdateParams(raw, 'cookie')"
        :on-delete="(raw) => onDeleteParams(raw, 'cookie')"
      />
    </ToggleHeading>

    <ToggleHeading title="Query" type="card" :expand="false">
      <SimpleParameterEditor
        v-model="parameters.query"
        :draggable="false"
        :on-create="(raw) => onCreateParams(raw, 'query')"
        :on-change="(raw) => onUpdateParams(raw, 'query')"
        :on-delete="(raw) => onDeleteParams(raw, 'query')"
      />
    </ToggleHeading>

    <ToggleHeading title="Path" type="card" :expand="false">
      <SimpleParameterEditor
        v-model="parameters.path"
        :draggable="false"
        :on-create="(raw) => onCreateParams(raw, 'path')"
        :on-change="(raw) => onUpdateParams(raw, 'path')"
        :on-delete="(raw) => onDeleteParams(raw, 'path')"
      />
    </ToggleHeading>
  </div>
</template>
<script setup lang="ts">
import SimpleParameterEditor from '@/components/APIEditor/SimpleEditor.vue'
import { useProjectId } from '@/hooks/useProjectId'

import useApi from '@/hooks/useApi'
import uesGlobalParametersStore from '@/store/globalParameters'
import { storeToRefs } from 'pinia'

const project_id = useProjectId()
const globalParametersStore = uesGlobalParametersStore()
const { parameters } = storeToRefs(globalParametersStore)
const [isLoading, getCommonParamListApi] = useApi(globalParametersStore.getGlobalParameters)()

const onCreateParams = async (raw: any, type: string) => await globalParametersStore.addGlobalParameter(project_id, type, raw)
const onUpdateParams = async (raw: any, type: string) => await globalParametersStore.updateGlobalParameter(project_id, type, raw)
const onDeleteParams = async (raw: any, type: string) => await globalParametersStore.deleteGlobalParameter(project_id, type, raw)

onMounted(async () => await getCommonParamListApi(project_id as string))
</script>
