<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { type FormInstance, type FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import CollapseCardItem from '@/components/collapse/CollapseCardItem.vue'
import { useCollapse } from '@/components/collapse/useCollapse'
import useApi from '@/hooks/useApi'
import useAppStore from '@/store/app'
import { apiGetOAuth, apiUpdateOAuth } from '@/api/system'
import { notNullRule } from '@/commons'

const { t } = useI18n()
const appStore = useAppStore()
const collapse = useCollapse({ defaults: ['github'] })
const [isLoading, updateOAuth] = useApi(apiUpdateOAuth)

const tBase = 'app.system.oauth'
const form = ref<SystemAPI.OAuthData>({
  clientID: '',
  clientSecret: '',
})

const formRef = ref<FormInstance>()
const rules = reactive<FormRules<SystemAPI.OAuthData>>({
  clientID: notNullRule(t(`${tBase}.rules.clientID`)),
  clientSecret: notNullRule(t(`${tBase}.rules.clientSecret`)),
})

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    await updateOAuth(form.value)
    appStore.updatGithuClienId(form.value.clientID)
  }
  catch (error) {
    //
  }
}

apiGetOAuth().then((v) => {
  form.value = v
})
</script>

<template>
  <div class="bg-white w-85%">
    <h1 class="text-30px">
      {{ $t('app.system.oauth.title') }}
    </h1>
    <div class="mt-40px">
      <CollapseCardItem name="github" :collapse-ctx="collapse">
        <template #title>
          <div class="row-lr">
            <div class="left mr-8px">
              <Icon icon="uil:github" width="24" />
            </div>
            <div class="font-bold right">
              GitHub
            </div>
          </div>
        </template>
        <ElForm
          ref="formRef"
          label-position="top"
          :rules="rules"
          :model="form"
          @submit.prevent="handleSubmit"
          @keyup.enter="handleSubmit"
        >
          <!-- id -->
          <ElFormItem prop="clientID" :label="$t('app.system.oauth.github.id')">
            <ElInput v-model="form.clientID" maxlength="255" />
          </ElFormItem>

          <!-- secret -->
          <ElFormItem prop="clientSecret" :label="$t('app.system.oauth.github.secret')">
            <ElInput v-model="form.clientSecret" maxlength="255" />
          </ElFormItem>
        </ElForm>

        <el-button type="primary" :loading="isLoading" @click="handleSubmit">
          {{ $t('app.common.save') }}
        </el-button>
      </CollapseCardItem>
    </div>
  </div>
</template>

<style scoped>
:deep(.el-button) {
  height: 40px;
}

.content {
  margin-top: 40px;
}
</style>
