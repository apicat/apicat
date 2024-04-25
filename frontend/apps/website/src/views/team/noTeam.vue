<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import useApi from '@/hooks/useApi'
import { apiCreateTeam } from '@/api/team'
import { useTeamStore } from '@/store/team'
import { TEAM_NAME } from '@/router'

const { t } = useI18n()
const store = useTeamStore()
const form = ref<TeamAPI.RequestCreateTeam>({
  name: '',
})
const formRef = ref<FormInstance>()
const rules: FormRules<typeof form.value> = {
  name: [
    {
      required: true,
      type: 'string',
      message: t('app.team.create.rule'),
      trigger: 'blur',
    },
  ],
}
const router = useRouter()
const [isLoading, createTeam] = useApi(apiCreateTeam)
async function submit() {
  try {
    await formRef.value?.validate()
    await createTeam(form.value as UserAPI.RequestGeneral)
    await store.init()
    router.push({
      name: TEAM_NAME,
    })
  }
  catch (e) {
    //
  }
}
</script>

<template>
  <div class="flex flex-col items-center justify-center h-full">
    <p class="text-40px">
      👬
    </p>

    <h1 class="text-gray-title font-500">
      {{ $t(`app.team.create.${store.hasTeam ? 'title1' : 'title0'}`) }}
    </h1>

    <ElForm ref="formRef" :rules="rules" :model="form" class="content w-400px" label-position="top">
      <div style="margin-top: 20px">
        <!-- name -->
        <ElFormItem class="w-full" prop="name" :label="$t('app.team.create.name')">
          <ElInput v-model="form.name" maxlength="255" @keydown.enter.prevent="submit" />
        </ElFormItem>
      </div>

      <!-- submit -->
      <ElButton :loading="isLoading" class="w-full" type="primary" @click="submit">
        {{ $t('app.team.create.submit') }}
      </ElButton>
    </ElForm>
  </div>
</template>

<style scoped>
h1 {
  font-size: x-large;
  font-weight: bold;
}
</style>
