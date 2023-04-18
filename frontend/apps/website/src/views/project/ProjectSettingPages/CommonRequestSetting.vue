<template>
  <div v-loading="isLoading">
    <ToggleHeading title="Header" type="card" v-loading="savingHeader" :expand="false">
      <SimpleParameterEditor v-model="data.header" />
      <el-button class="mt-10px" type="primary" :loading="savingHeader" @click="handleSubmit(saveHeaderParamerterApi, 'header')">{{ $t('app.common.save') }}</el-button>
    </ToggleHeading>

    <ToggleHeading title="Cookie" type="card" v-loading="savingCookie" :expand="false">
      <SimpleParameterEditor v-model="data.cookie" />
      <el-button class="mt-10px" type="primary" :loading="savingCookie" @click="handleSubmit(saveCookieParamerterApi, 'cookie')">{{ $t('app.common.save') }}</el-button>
    </ToggleHeading>

    <ToggleHeading title="Query" type="card" v-loading="savingQuery" :expand="false">
      <SimpleParameterEditor v-model="data.query" />
      <el-button class="mt-10px" type="primary" :loading="savingQuery" @click="handleSubmit(saveQueryParamerterApi, 'query')">{{ $t('app.common.save') }}</el-button>
    </ToggleHeading>
  </div>
</template>
<script setup lang="ts">
import { getCommonParamList, saveHeaderParamerter, saveQueryParamerter, saveCookieParamerter } from '@/api/param'

import SimpleParameterEditor from '@/components/APIEditor/SimpleEditor.vue'
import { useProjectId } from '@/hooks/useProjectId'

const project_id = useProjectId()

const [isLoading, getCommonParamListApi] = getCommonParamList()
const [savingHeader, saveHeaderParamerterApi] = saveHeaderParamerter()
const [savingQuery, saveQueryParamerterApi] = saveQueryParamerter()
const [savingCookie, saveCookieParamerterApi] = saveCookieParamerter()

const data: any = ref({})

onMounted(async () => {
  data.value = await getCommonParamListApi({ project_id })
})

const handleSubmit = async (api: any, key: string) => await api({ project_id, content: toRaw(unref(data)[key]) })
</script>
