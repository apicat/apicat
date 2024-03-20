<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import type { UseCollapse } from '@/components/collapse/useCollapse'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import Iconfont from '@/components/Iconfont.vue'
import useApi from '@/hooks/useApi'
import { apiUpdateStroageLocal } from '@/api/system'
import { notNullRule } from '@/commons'

const props = defineProps<{ collapse: UseCollapse; name: string; config: Partial<SystemAPI.StorageDisk> }>()
const { t } = useI18n()
const tBase = 'app.system.storage.local'
const formRef = ref<FormInstance>()
const rules: FormRules<typeof props.config> = {
  path: notNullRule(t(`${tBase}.rules.path`)),
}
const [submitting, updateLocal] = useApi(apiUpdateStroageLocal)
function submit() {
  formRef.value!.validate((valid) => {
    if (valid)
      updateLocal(props.config as SystemAPI.StorageDisk)
  })
}
</script>

<template>
  <CollapseCardItem :name="name" :collapse-ctx="collapse">
    <template #title>
      <div class="row-lr">
        <div class="left mr-8px">
          <Iconfont icon="ac-memory-one" :size="24" />
        </div>
        <div class="right font-bold">
          {{ $t(`${tBase}.title`) }}
        </div>
      </div>
    </template>
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="props.config" @submit.prevent="submit">
      <ElFormItem prop="path" :label="$t(`${tBase}.path`)">
        <ElInput v-model="props.config.path" maxlength="255" />
      </ElFormItem>
    </ElForm>

    <el-button :loading="submitting" type="primary" @click="submit">
      {{ $t('app.common.update') }}
    </el-button>
  </CollapseCardItem>
</template>

<style scoped></style>
