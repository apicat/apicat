<!-- eslint-disable ts/no-non-null-asserted-optional-chain -->
<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/store/user'
import { useTeamStore } from '@/store/team'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { apiDeleteTeam, apiGetTeamActiveMembers, apiTeamSetting, apiTransferOwnership } from '@/api/team'
import useApi from '@/hooks/useApi'
import { Role } from '@/commons/constant'

const { t } = useI18n()
const userStore = useUserStore()
const teamStore = useTeamStore()
const { userInfo } = storeToRefs(userStore)
const { currentTeam, currentID } = storeToRefs(teamStore)

const formRefName = ref<FormInstance>()
const formName = ref<TeamAPI.RequestTeamSetting>({
  name: '',
})
const setting = computed(() => {
  formName.value.name = currentTeam!.value!.name
  return formName.value
})
const rules: FormRules<typeof formName.value> = {
  name: [
    {
      required: true,
      trigger: 'blur',
    },
  ],
}
const [settingLoading, updateTeamSetting] = useApi(apiTeamSetting)
async function submitName() {
  try {
    await formRefName.value?.validate()
    await updateTeamSetting(currentID.value, formName.value)
    teamStore.updateTeam(currentID.value, {
      name: formName.value.name,
      avatar: formName.value.avatar,
    })
  }
  catch (e) {
    //
  }
}
const formTransfer = ref<Partial<TeamAPI.RequestTransferOwnership>>({
  memberID: undefined,
})
const rulesTransfer: FormRules<typeof formTransfer.value> = {
  memberID: [
    {
      required: true,
      trigger: 'blur',
      message: t('app.team.setting.transfer.selectAdmin'),
    },
  ],
}
const formRefOwnership = ref<FormInstance>()
const [_, transferOwnership] = useApi(apiTransferOwnership)
const [getMembersLoading, getTeamMembers] = useApi(apiGetTeamActiveMembers)
async function submitOwnership() {
  try {
    await formRefOwnership.value?.validate()
    AsyncMsgBox({
      confirmButtonClass: 'red',
      confirmButtonText: t('app.team.setting.transfer.btn'),
      cancelButtonText: t('app.common.cancel'),
      title: t('app.team.setting.transfer.title'),
      content: t('app.team.setting.transfer.pop.tip'),
      onOk: async () => {
        if (!formTransfer.value.memberID)
          return
        await transferOwnership(teamStore.currentID, {
          memberID: formTransfer.value.memberID,
        })
        window.location.reload()
      },
    })
  }
  catch (e) {}
}
const memberOptions = ref<{ value: number; label: string }[]>([])

async function getMembers() {
  // TODO: 骚代码
  const get = () =>
    getTeamMembers(
      {
        roles: Role.Admin,
        page: 0,
        pageSize: 10000,
      },
      currentID.value,
    )

  const list: TeamAPI.TeamMember[] = (await get())!.items!
  // filter and push
  for (const index in list) {
    const val = list[index as any]
    if (userInfo.value.id === val.id)
      continue
    memberOptions.value.push({ value: val.id, label: val.user.name })
  }
}
getMembers()

const [__, deleteTeam] = useApi(apiDeleteTeam)
async function submitDeletion() {
  AsyncMsgBox({
    confirmButtonClass: 'red',
    confirmButtonText: t('app.common.delete'),
    cancelButtonText: t('app.common.cancel'),
    title: t('app.team.setting.rm.title'),
    content: t('app.team.setting.rm.pop.tip'),
    onOk: async () => {
      await deleteTeam(teamStore.currentID)
      window.location.reload()
    },
  })
}
</script>

<template>
  <div class="max-w-600px">
    <!-- name -->
    <div class="content">
      <h2 class="mb-3">
        {{ $t('app.team.setting.name.title') }}
      </h2>
      <ElForm ref="formRefName" label-position="top" :rules="rules" :model="formName">
        <ElFormItem prop="name">
          <ElInput v-model="setting.name" maxlength="255" class="h-40px" style="width: 500px" />
        </ElFormItem>
        <ElButton class="mt-3 w-100px" :loading="settingLoading" @click="submitName">
          {{ $t('app.team.setting.name.btn') }}
        </ElButton>
      </ElForm>
    </div>
    <!-- transfer -->
    <div class="content">
      <h2 class="mb-3">
        {{ $t('app.team.setting.transfer.title') }}
      </h2>
      <ElForm ref="formRefOwnership" label-position="top" :rules="rulesTransfer" :model="formTransfer">
        <ElFormItem prop="memberID">
          <ElSelect v-model="formTransfer.memberID" :loading="getMembersLoading" style="width: 500px">
            <el-option v-for="item in memberOptions" :key="item.value" :label="item.label" :value="item.value" />
          </ElSelect>
        </ElFormItem>
        <p class="mt-20px">
          {{ $t('app.team.setting.transfer.tip') }}
        </p>
        <ElButton class="mt-3 red-outline w-100px" type="primary" @click="submitOwnership">
          {{ $t('app.team.setting.transfer.btn') }}
        </ElButton>
      </ElForm>
    </div>
    <!-- delete -->
    <div class="content">
      <h2 class="mb-3">
        {{ $t('app.team.setting.rm.title') }}
      </h2>
      <p class="tip">
        {{ $t('app.team.setting.rm.tip') }}
      </p>
      <ElButton class="mt-3 red-outline w-100px" @click="submitDeletion">
        {{ $t('app.team.setting.rm.btn') }}
      </ElButton>
    </div>
  </div>
</template>

<style scoped>
:deep(.el-form-item) {
  margin-bottom: 0;
}
</style>
