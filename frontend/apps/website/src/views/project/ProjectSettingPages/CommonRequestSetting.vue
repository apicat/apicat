<template>
  <div v-loading="isLoading">
    <ToggleHeading title="Header" type="card" :expand="false">
      <SimpleParameterEditor
        v-model="data.header"
        :draggable="false"
        :on-create="(raw) => onCreateParams(raw, 'header')"
        :on-change="(raw) => onUpdateParams(raw, 'header')"
        :on-delete="(raw) => onDeleteParams(raw, 'header')"
      />
    </ToggleHeading>

    <ToggleHeading title="Cookie" type="card" :expand="false">
      <SimpleParameterEditor v-model="data.cookie" :draggable="false" :on-create="(raw) => onCreateParams(raw, 'cookie')" />
    </ToggleHeading>

    <ToggleHeading title="Query" type="card" :expand="false">
      <SimpleParameterEditor v-model="data.query" :draggable="false" :on-create="(raw) => onCreateParams(raw, 'query')" />
    </ToggleHeading>
  </div>
</template>
<script setup lang="ts">
import { getGlobalParamList, createGlobalParamerter, updateGlobalParamerter, deleteGlobalParamerter } from '@/api/param'

import SimpleParameterEditor from '@/components/APIEditor/SimpleEditor.vue'
import { useProjectId } from '@/hooks/useProjectId'

const project_id = useProjectId()

const [isLoading, getCommonParamListApi] = getGlobalParamList()

const data: any = ref({ header: [], cookie: [], query: [] })

const onCreateParams = async (raw: any, type: string) => await createGlobalParamerter({ project_id, in: type, ...raw })
const onUpdateParams = async (raw: any, type: string) => await updateGlobalParamerter({ project_id, in: type, ...raw })
const onDeleteParams = async (raw: any, type: string) => await deleteGlobalParamerter({ project_id, in: type, ...raw })

onMounted(async () => {
  data.value = await getCommonParamListApi({ project_id })
})
</script>
