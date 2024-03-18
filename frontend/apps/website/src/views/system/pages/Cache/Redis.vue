<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import { type UseCollapse } from '@/components/collapse/useCollapse'
import IconSvg from '@/components/IconSvg.vue'
import useApi from '@/hooks/useApi'
import { apiUpdateCacheRedis } from '@/api/system'
import { SysCache, isUrlRule, notNullRule } from '@/commons'

const { t } = useI18n()
const tBase = 'app.system.cache.redis'
const props = defineProps<{
  collapse: UseCollapse
  name: SysCache
  config: Partial<SystemAPI.CacheRedis>
  currentUse?: SysCache
}>()
const emit = defineEmits(['update:currentUse'])

const formRef = ref<FormInstance>()
const rules: FormRules<typeof props.config> = {
  host: [...isUrlRule(t(`${tBase}.rules.host`), true, false), ...notNullRule(t(`${tBase}.rules.host`))],
  database: notNullRule(t(`${tBase}.rules.db`)),
}
const [submitting, update] = useApi(apiUpdateCacheRedis)
function submit() {
  formRef.value!.validate((valid) => {
    if (valid) update(props.config as SystemAPI.CacheRedis)
    emit('update:currentUse', props.name)
  })
}
</script>

<template>
  <CollapseCardItem :name="name" :collapse-ctx="collapse">
    <template #title>
      <div class="row-lr">
        <div class="left mr-8px">
          <IconSvg name="ac-redis" width="24" />
        </div>
        <div class="right font-bold">{{ $t(`${tBase}.title`) }}</div>
      </div>
    </template>
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="props.config" @submit.prevent="submit">
      <!-- host -->
      <ElFormItem prop="host" :label="$t(`${tBase}.host`)">
        <ElInput maxlength="255" v-model="props.config.host" />
      </ElFormItem>

      <!-- pw -->
      <ElFormItem prop="password" :label="$t(`${tBase}.pw`)">
        <ElInput maxlength="255" v-model="props.config.password" />
      </ElFormItem>

      <!-- db -->
      <ElFormItem prop="database" :label="$t(`${tBase}.db`)">
        <ElInput maxlength="255" v-model="props.config.database" />
      </ElFormItem>
    </ElForm>

    <el-button :loading="submitting" type="primary" @click="submit">
      {{ $t('app.common.update') }}
    </el-button>
  </CollapseCardItem>
</template>

<style scoped></style>
