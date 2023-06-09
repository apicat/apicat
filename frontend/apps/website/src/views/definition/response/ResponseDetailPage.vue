<template>
  <div class="ac-header-operate" v-if="hasDocument && response">
    <div class="ac-header-operate__main">
      <p class="ac-header-operate__title">{{ response.name }}</p>
    </div>

    <div class="ac-header-operate__btns" v-if="!isReader">
      <el-button type="primary" @click="() => goResponseEditPage()">{{ $t('app.common.edit') }}</el-button>
    </div>
  </div>

  <Result v-show="!hasDocument && !isLoading">
    <template #icon>
      <img class="h-auto w-260px mb-26px" src="@/assets/images/icon-empty.png" alt="" />
    </template>
  </Result>

  <div :class="[ns.b(), { 'h-20vh': !response && hasDocument }]" v-loading="isLoading">
    <template v-if="response">
      <h4>{{ response.description }}</h4>
      <DefinitionResponseRaw :response="response" :definition-schemas="definitionSchemas" />
    </template>
  </div>
</template>
<script setup lang="ts">
import DefinitionResponseRaw from '@/components/DefinitionResponse/DefinitionResponseRaw.vue'
import { useNamespace } from '@/hooks'
import { useGoPage } from '@/hooks/useGoPage'
import { useDefinitionResponseLogic } from './logic'
import uesProjectStore from '@/store/project'
import { storeToRefs } from 'pinia'

const ns = useNamespace('document')
const { goResponseEditPage } = useGoPage()
const projectStore = uesProjectStore()

const { isReader } = storeToRefs(projectStore)
const { hasDocument, isLoading, response, definitionSchemas } = useDefinitionResponseLogic()
</script>
