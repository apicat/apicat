<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { delay } from '@apicat/shared'
import useClipboard from '@/hooks/useClipboard'
import { apiGetInviteToken, apiResetInviteToken } from '@/api/team'
import { useTeamStore } from '@/store/team'
import useApi from '@/hooks/useApi'
import Page30pLayout from '@/layouts/Page30pLayout.vue'
import { getTeamJoinPath } from '@/router/team'

const store = useTeamStore()
const { t } = useI18n()
const [loading, getInviteToken] = useApi(apiGetInviteToken)
const token = ref('')
const link = computed(() => token.value ? getTeamJoinPath(token.value) : '')

const { handleCopy, elCopyTextRef } = useClipboard(link, t('app.team.invite.copy'), t('app.project.collection.share.copied'))

async function resetToken() {
  loading.value = true
  try {
    const res = await apiResetInviteToken(store.currentTeam!.id!)
    token.value = res.invitationToken
  }
  catch (error) {
    //
  }
  finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const { invitationToken } = await getInviteToken(store.currentTeam!.id!) as TeamAPI.ResponseInvite
    token.value = invitationToken
  }
  catch (error) {
    token.value = ''
  }
})
</script>

<template>
  <Page30pLayout>
    <div class="w-full max-w-600px" style="background-color: white">
      <h1>{{ $t('app.team.invite.title') }}</h1>
      <div v-loading="loading" element-loading-background="rgba(122, 122, 122, 0)" class="mt-40px">
        <ElInput v-model="link" class="h-40px " readonly>
          <template #append>
            <ElButton
              style="height: 40px;"
              type="primary"
              @click="handleCopy"
            >
              {{ elCopyTextRef }}
            </ElButton>
          </template>
        </ElInput>
      </div>

      <div class="w-full mt-5 text-gray-helper text-14px">
        <p class="mb-4px">
          {{ $t('app.team.invite.tip1') }}
        </p>

        <p>
          <span> {{ $t('app.team.invite.tip2.text') }} </span>
          <ElButton link type="primary" @click="resetToken">
            {{ $t('app.team.invite.tip2.link') }}
          </ElButton>
        </p>
      </div>
    </div>
  </Page30pLayout>
</template>

<style scoped>
:deep(.el-button.is-link) {
  padding: 0;
  vertical-align: unset;
}
</style>
