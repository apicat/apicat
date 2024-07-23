<script setup lang="tsx">
import { apiGetDefaultModel, apiGetModel, apiUpdateDefaultModel } from '@/api/system'
import { notNullRule, SysModel } from '@/commons'
import { useCollapse } from '@/components/collapse/useCollapse'
import useApi from '@/hooks/useApi'
import Azure from '@/views/system/pages/Model/Azure.vue'
import OpenAI from '@/views/system/pages/Model/OpenAI.vue'
import IconSvg from '@/components/IconSvg.vue'
import { useI18n } from 'vue-i18n'
import { FormInstance } from 'element-plus'

const tBase = 'app.system.model'
const collapse = useCollapse({})
const { t } = useI18n()

const [isLoading, updateDefaultModel] = useApi(apiUpdateDefaultModel)
const localDefaultConfig = ref({
  llm:'',
  embedding:''
})

const localDefaultConfigRules = {
  // llm: [{message: t(`${tBase}.openai.rules.llmName`), required: true, trigger: 'change'}],
  llm: notNullRule(t(`${tBase}.openai.rules.llmName`)),
  embedding: notNullRule(t(`${tBase}.openai.rules.embedding`)),
}

const localDefaultRequestData = computed(()=>{
  const [driver,model] = localDefaultConfig.value.llm.split('/')
  const [driver2,model2] = localDefaultConfig.value.embedding.split('/')
  return {
    llm:{driver,model},
    embedding:{driver:driver2,model:model2}
  }
})

const formRef = ref<FormInstance>()
const data = ref<SystemAPI.ModelItem[]>([])
const localDefaultConfigData = ref<SystemAPI.ModelDefaultConfig>()
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

const defaultLLMOptions = computed(()=> {
  return localDefaultConfigData.value?.llm.map((item) => {
    const value = `${item.driver}/${item.model}`
    if(item.selected){
      localDefaultConfig.value.llm = value
    }

    return {
      icon:item.driver === SysModel.OpenAI ? <IconSvg name="ac-openai"/> : <IconSvg name="ac-azure" width="24" />,
      value:value
    }
  }) || []
})

const defaultEmbeddingOptions = computed(()=> {
  return localDefaultConfigData.value?.embedding.map((item) => {
    const value = `${item.driver}/${item.model}`
    if(item.selected){
      localDefaultConfig.value.embedding = value
    }
    return {
      icon:item.driver === SysModel.OpenAI ? <IconSvg name="ac-openai"/> : <IconSvg name="ac-azure" width="24" />,
      value:value
    }
  }) || []
})

async function handleSubmit(formIns:FormInstance | undefined){
  try {
    if (!formIns) return
    await formIns.validate()
    await updateDefaultModel(localDefaultRequestData.value)
  } catch (error) {
    //
  }
}

onBeforeMount(async () => {
  try {
      const [config,defaultConfig] = await Promise.all([apiGetModel(),apiGetDefaultModel()])
      data.value = config || []
      localDefaultConfigData.value = defaultConfig || {}
  } catch (error) {
    //
  }
  
})

</script>

<template>
  <div class="bg-white w-85%">
    <h1>{{ $t(`${tBase}.title`) }}</h1>
    <ElForm ref="formRef" label-position="top" :model="localDefaultConfig" :rules="localDefaultConfigRules" @submit.prevent="handleSubmit(formRef)">
      <ElFormItem label="Reasoning model" prop="llm">
        <ElSelect v-model="localDefaultConfig.llm" class="w-full" placeholder="Select model">
          <ElOption v-for="item in defaultLLMOptions" :key="item.value" :label="item.value" :value="item.value" >
            <div class="flex-y-center">
            <component :is="item.icon" />
            <span class="ml-4px">{{ item.value }}</span>
          </div>
          </ElOption>
        </ElSelect>
      </ElFormItem>

      <ElFormItem label="Embedding model" prop="embedding">
        <ElSelect v-model="localDefaultConfig.embedding" class="w-full" placeholder="Select model">
          <ElOption v-for="item in defaultEmbeddingOptions" :key="item.value" :label="item.value" :value="item.value" >
            <div class="flex-y-center">
            <component :is="item.icon" />
            <span class="ml-4px">{{ item.value }}</span>
          </div>
          </ElOption>
        </ElSelect>
      </ElFormItem>

      <el-button :loading="isLoading" class="w-full" type="primary" @click="handleSubmit(formRef)">
      {{ $t('app.common.save') }}
    </el-button>
    </ElForm>

    <div class="my-40px">
      <ElDivider>Model provider</ElDivider>
    </div>

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
