<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import type { UseCollapse } from '@/components/collapse/useCollapse'
import { apiUpdateCacheLocal } from '@/api/system'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import IconSvg from '@/components/IconSvg.vue'
import useApi from '@/hooks/useApi'
import { useI18n } from 'vue-i18n'
import { SysCache } from '@/commons'

const { t } = useI18n()
const tBase = 'app.system.cache.local'
const props = defineProps<{
  collapse: UseCollapse
  name: SysCache
  config: Partial<SystemAPI.CacheMemory>
  currentUse?: SysCache
}>()
const emit = defineEmits(['update:currentUse'])
const form = ref({
  checked: false,
})
watch(
  () => props.currentUse,
  (v) => {
    form.value.checked = v === props.name
  },
  {
    immediate: true,
  },
)
const formRef = ref<FormInstance>()
const rules: FormRules<typeof form.value> = {
  checked: [
    {
      validator: (_, __, c) => {
        if (!form.value.checked) return c(t(`${tBase}.rules.checked`))
        else return c()
      },
      trigger: 'blur',
    },
  ],
}

const [submitting, update] = useApi(apiUpdateCacheLocal)
function submit() {
  formRef.value?.validate().then(() => {
    update()
    emit('update:currentUse', props.name)
  })
}
</script>

<template>
  <CollapseCardItem :name="name" :collapse-ctx="collapse">
    <template #title>
      <div class="row-lr">
        <div class="left mr-8px">
          <IconSvg name="ac-storage-card-one" width="24" />
        </div>
        <div class="right font-bold">{{ $t(`${tBase}.title`) }}</div>
      </div>
    </template>

    <ElForm ref="formRef" label-position="top" :rules="rules" :model="form" @submit.prevent="submit">
      <p class="font-bold">
        {{ $t(`${tBase}.smallTitle`) }}
      </p>

      <!-- checked -->
      <ElFormItem prop="checked">
        <el-checkbox v-model="form.checked" :label="$t(`${tBase}.tip`)" size="large" />
      </ElFormItem>

      <el-button class="mt-10px" :loading="submitting" type="primary" @click="submit">
        {{ $t('app.common.update') }}
      </el-button>
    </ElForm>
  </CollapseCardItem>
</template>

<style scoped></style>
