<script setup lang="ts">
import { BadRequestError } from '@/api/error'
import { type InviteTokenTeamData, apiJoinTeam, getJoinTeamInfoByToken } from '@/api/team'
import useApi from '@/hooks/useApi'
import { useGlobalLoading } from '@/hooks/useGlobalLoading'
import { LOGIN_NAME, MAIN_PATH_NAME } from '@/router'
import { useUserStore } from '@/store/user'

const { showGlobalLoading, hideGlobalLoading } = useGlobalLoading()
const userStore = useUserStore()
const route = useRoute()
const router = useRouter()
const token = route.params.token as string
const [loading, joinTeam] = useApi(apiJoinTeam)

const teamData = ref<InviteTokenTeamData>()
const errorMsg = ref()
const [, getTeamDetail] = useApi(getJoinTeamInfoByToken)
showGlobalLoading()

async function join() {
  if (userStore.isLogin) {
    await joinTeam(token)
    router.push({ name: MAIN_PATH_NAME })
    return
  }

  // 未登录
  router.push({
    name: LOGIN_NAME,
    query: {
      invitationToken: token,
    },
  })
}
getTeamDetail(token).then((res) => {
  teamData.value = res
})
  .catch((error) => {
    if (error instanceof BadRequestError)
      errorMsg.value = error.response || {}
  })
  .finally(() => hideGlobalLoading())
</script>

<template>
  <div v-if="teamData" class="flex min-h-screen">
    <main class="p-2 ac-login">
      <div class="shadow-xl ac-login__box">
        <div class="text-center">
          <AcLogo href="/" />
          <h1 class="mt-5 text-gray-title">
            {{ $t('app.team.join.title') }}
          </h1>
          <p class="mt-3" style="line-height: 20px; word-break: break-word">
            <span class="bold mr-2px">{{ teamData?.inviter }}</span><span>{{ $t('app.team.join.tip1') }}</span><span class="bold mr-4px">{{ teamData?.team }}</span><span> {{ $t('app.team.join.tip2') }} </span>
          </p>

          <el-button :loading="loading" class="w-full mt-8 b-btn" type="primary" @click="join">
            {{ $t('app.team.join.btn') }}
          </el-button>
        </div>
      </div>
    </main>
  </div>
  <MessageTemplate v-else v-bind="errorMsg" />
</template>

<style scoped>
h1 {
  font-size: 24px;
  font-weight: bold;
  font-style: normal;
  writing-mode: horizontal-tb;
  color: rgb(16, 16, 16);
  letter-spacing: 0px;
  line-height: 33px;
  text-shadow: none;
}

.bold {
  font-weight: bold;
}

.b-btn {
  height: 40px;
}

.ac-login__box {
  width: 600px;
  min-height: 380px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
