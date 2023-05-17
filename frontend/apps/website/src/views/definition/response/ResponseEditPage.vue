<template>
  <div class="ac-header-operate">
    <div class="ac-header-operate__main">
      <p class="flex-y-center">
        <el-icon :size="18" class="mt-1px mr-4px"><ac-icon-ic-sharp-cloud-queue /></el-icon>
        {{ isSaving ? $t('app.common.saving') : $t('app.common.savedCloud') }}
      </p>
    </div>
    <div class="ac-header-operate__btns">
      <el-button type="primary" @click="() => goResponseDetailPage()">{{ $t('app.common.preview') }}</el-button>
    </div>
  </div>

  <div :class="[ns.b(), { 'h-50vh': !response }]" v-loading="isLoading">
    <template v-if="response">
      <input class="ac-document__title" type="text" v-input-limit v-model="response.name" maxlength="255" :placeholder="$t('app.definitionResponse.form.title')" />
      <input class="ac-document__desc" type="text" v-model="response.description" maxlength="255" :placeholder="$t('app.definitionResponse.form.desc')" />
      <DefinitionResponseForm v-model:response="response" :definition-schemas="definitionSchemas" />
    </template>
  </div>
</template>
<script setup lang="ts">
import { useGoPage } from '@/hooks/useGoPage'
import { useNamespace } from '@/hooks'
import { useEditDefinitionResponseLogic } from './logic'
import DefinitionResponseForm from '@/components/DefinitionResponse/DefinitionResponseForm.vue'

const ns = useNamespace('document')
const { goResponseDetailPage } = useGoPage()
const { response, isSaving, isLoading, definitionSchemas } = useEditDefinitionResponseLogic()
</script>
