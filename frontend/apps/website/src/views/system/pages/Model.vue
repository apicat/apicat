<script setup lang="ts">
import { apiGetModel } from '@/api/system'
import { SysModel } from '@/commons'
import { useCollapse } from '@/components/collapse/useCollapse'
import Azure from '@/views/system/pages/Model/Azure.vue'
import OpenAI from '@/views/system/pages/Model/OpenAI.vue'

const tBase = 'app.system.model'
const collapse = useCollapse({})

interface A {
  [SysModel.OpenAI]: SystemAPI.ModelOpenAI
  [SysModel.Azure]: SystemAPI.ModelAzure
}
const data = ref<A>({
  [SysModel.OpenAI]: {
    apiKey: '',
    organizationID: '',
    apiBase: '',
    llmName: '',
  },
  [SysModel.Azure]: {
    apiKey: '',
    endpoint: '',
    llmName: '',
  },
})
apiGetModel().then((res) => {
  for (let i = 0; i < res.length; i++) {
    const v = res[i]
    data.value[v.driver as keyof A] = v.config as any
    if (v.use)
      collapse.ctx.open(v.driver)
  }
})
</script>

<template>
  <div class="bg-white w-85%">
    <h1>{{ $t(`${tBase}.title`) }}</h1>

    <div class="mt-40px flex flex-col">
      <OpenAI
        v-model:config="data[SysModel.OpenAI]"
        :name="SysModel.OpenAI"
        class="collapse-box"
        :collapse="collapse"
      />
      <Azure
        v-model:config="data[SysModel.Azure]"
        class="collapse-box mt-30px"
        :name="SysModel.Azure"
        :collapse="collapse"
      />
    </div>
  </div>
</template>

<style scoped>
h1 {
  font-size: 30px;
}

:deep(.el-select .el-input) {
  height: 40px;
}

:deep(.el-button) {
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
  /* flex-grow: 1; */
}

.right {
  /* justify-content: flex-end; */
  flex-grow: 1;
}

.content {
  margin-top: 40px;
}
</style>
