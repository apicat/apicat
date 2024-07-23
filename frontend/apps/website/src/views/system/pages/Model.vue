<script setup lang="ts">
import { apiGetModel } from '@/api/system'
import { SysModel } from '@/commons'
import { useCollapse } from '@/components/collapse/useCollapse'
import Azure from '@/views/system/pages/Model/Azure.vue'
import OpenAI from '@/views/system/pages/Model/OpenAI.vue'

const tBase = 'app.system.model'
const collapse = useCollapse({})

const data = ref<SystemAPI.ModelItem[]>([])
const llmModels = ref<string[]>([])
const embeddingModels = ref<string[]>([])
  
const findDriver = (driver:SysModel) => data.value.find((item) => item.driver === driver)

const openAIConfig = computed<SystemAPI.ModelOpenAI>(() => {
  const config = findDriver(SysModel.OpenAI)
  if(config){
    llmModels.value = config.models?.llm || []
    embeddingModels.value = config.models?.embedding || []
  }
  
  return config ? config.config as SystemAPI.ModelOpenAI : {
    apiKey:'',
    llm:'',
    embedding:''
  }
})

const azureConfig = computed<SystemAPI.ModelAzure>(() =>{
  const config = findDriver(SysModel.Azure)
  return config ? config.config as SystemAPI.ModelAzure : {
    apiKey:'',
    endpoint:'',
    llm:'',
    embedding:''
  }
})

onBeforeMount(async () => {
  const res = await apiGetModel()
  data.value = [
  {
    "driver": "openai",
    "config": {
      "apiKey": "abcdefg",
      "organizationID": "a1b2c3",
      "apiBase": "abcde12345",
      "llm": "gpt-3.5-turbo",
      "embedding": "text-embedding-3-small"
    },
    "models": {
      "llm": ["gpt-4-turbo", "gpt-4o", "gpt-3.5-turbo"],
      "embedding": ["text-embedding-3-small", "text-embedding-3-large"]
    }
  },
  {
    "driver": "azure-openai",
    "config": {
      "apiKey": "abcdefg",
      "endpoint": "https://test.azure-openai.com/wahaha/",
      "llm": "gpt-35-turbo",
      "embedding": "text-embedding-3-small"
    }
  }
] as SystemAPI.ModelItem[]
})

</script>

<template>
  <div class="bg-white w-85%">
    <h1>{{ $t(`${tBase}.title`) }}</h1>

    <div class="flex flex-col">
      <OpenAI
        :config="openAIConfig"
        :name="SysModel.OpenAI"
        class="collapse-box"
        :collapse="collapse"
        :embedding-models="embeddingModels"
        :llm-models="llmModels"
      />
      <Azure
        :config="azureConfig"
        class="collapse-box mt-30px"
        :name="SysModel.Azure"
        :collapse="collapse"
      />
    </div>
  </div>
</template>

<style scoped>
:deep(.el-select .el-input) {
  height: 40px;
}

.row {
  margin-top: 1em;
  margin-bottom: 1em;
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.left,
.right {
  display: flex;
  align-items: center;
}

.left {
  justify-content: flex-start;
}

.right {
  flex-grow: 1;
}
</style>
