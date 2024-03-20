<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import type { UseCollapse } from '@/components/collapse/useCollapse'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import IconSvg from '@/components/IconSvg.vue'
import { isUrlRule, notNullRule } from '@/commons'
import { apiUpdateEmailSMTP } from '@/api/system'
import useApi from '@/hooks/useApi'

const props = defineProps<{ collapse: UseCollapse; name: string; config: Partial<SystemAPI.EmailSMTP> }>()
const { t } = useI18n()
const tBase = 'app.system.email.smtp'
const formRef = ref<FormInstance>()
const rules: FormRules<typeof props.config> = {
  host: [...isUrlRule(t(`${tBase}.rules.host`), true, false), ...notNullRule(t(`${tBase}.rules.host`))],
  address: notNullRule(t(`${tBase}.rules.address`)),
  password: notNullRule(t(`${tBase}.rules.pw`)),
}

const [submitting, update] = useApi(apiUpdateEmailSMTP)
function submit() {
  formRef.value!.validate((valid) => {
    if (valid)
      update(props.config as SystemAPI.EmailSMTP)
  })
}
</script>

<template>
  <CollapseCardItem :name="name" :collapse-ctx="collapse">
    <template #title>
      <div class="row-lr">
        <div class="left mr-8px">
          <IconSvg name="ac-mail" width="24" />
        </div>
        <div class="right font-bold">
          {{ $t(`${tBase}.title`) }}
        </div>
      </div>
    </template>
    <ElForm ref="formRef" label-position="top" :rules="rules" :model="props.config" @submit.prevent="submit">
      <!-- host -->
      <ElFormItem prop="host" :label="$t(`${tBase}.host`)">
        <ElInput v-model="props.config.host" maxlength="255" />
      </ElFormItem>

      <!-- user -->
      <ElFormItem prop="user" :label="$t(`${tBase}.user`)">
        <ElInput v-model="props.config.user" maxlength="255" />
      </ElFormItem>

      <!-- address -->
      <ElFormItem prop="address" :label="$t(`${tBase}.address`)">
        <ElInput v-model="props.config.address" maxlength="255" />
      </ElFormItem>

      <!-- pw -->
      <ElFormItem prop="password" :label="$t(`${tBase}.pw`)">
        <ElInput v-model="props.config.password" maxlength="255" />
      </ElFormItem>
    </ElForm>

    <el-button :loading="submitting" type="primary" @click="submit">
      {{ $t('app.common.update') }}
    </el-button>
  </CollapseCardItem>
</template>

<style scoped></style>
