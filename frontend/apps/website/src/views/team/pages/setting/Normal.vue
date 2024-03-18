<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import useApi from '@/hooks/useApi'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { apiQuitTeam } from '@/api/team'
import { useTeamStore } from '@/store/team'

const { t } = useI18n()
const [__, quit] = useApi(apiQuitTeam)
const teamStore = useTeamStore()
async function submitDeletion() {
  AsyncMsgBox({
    confirmButtonClass: 'red',
    confirmButtonText: t('app.team.setting.quit.btn'),
    cancelButtonText: t('app.common.cancel'),
    title: t('app.team.setting.quit.title'),
    content: t('app.team.setting.quit.pop.tip'),
    onOk: async () => {
      await quit(teamStore.currentTeam!.id)
      window.location.reload()
    },
  })
}
</script>

<template>
  <!-- quit -->
  <div class="content">
    <h2 class="mb-3">
      {{ $t('app.team.setting.quit.title') }}
    </h2>
    <p class="tip">
      {{ $t('app.team.setting.quit.tip') }}
    </p>
    <ElButton class="mt-3 red-outline w-100px" @click="submitDeletion">
      {{ $t('app.team.setting.quit.btn') }}
    </ElButton>
  </div>
</template>

<style scoped></style>
