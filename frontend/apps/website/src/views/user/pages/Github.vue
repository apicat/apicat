<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { storeToRefs } from 'pinia'
import { apiDisconnectOAuth } from '@/api/user'
import { OAuthPlatform } from '@/commons/constant'
import { useUserStore } from '@/store/user'
import { createOAuthConnectCallbackURL } from '@/api/sign/oAuth'
import useApi from '@/hooks/useApi'
import { useAppStore } from '@/store/app'

const userStore = useUserStore()
const { oAuthURLConfig } = storeToRefs(useAppStore())
const { userInfo } = storeToRefs(userStore)

const [disconnectLoading, disconnectApi] = useApi(apiDisconnectOAuth)
async function connect() {
  location.href = oAuthURLConfig.value.github(createOAuthConnectCallbackURL(OAuthPlatform.GITHUB))
}

async function disconnect() {
  try {
    await disconnectApi(OAuthPlatform.GITHUB as unknown as SignAPI.OAuthPlatform)
    await userStore.getUserInfo()
  }
  catch (error) {
    //
  }
}
</script>

<template>
  <div class="flex flex-col items-center justify-center mx-auto px-36px">
    <!-- position -->
    <div class="items-start text-start w-40vw">
      <!-- content -->
      <div class="bg-white w-450px">
        <h1 class="text-30px">
          {{ $t('app.user.github.title') }}
        </h1>
        <div class="mt-40px">
          <div class="row">
            <div class="left mr-8px">
              <Icon icon="uil:github" width="24" />
            </div>
            <div class="right">
              GitHub
            </div>
          </div>
          <p>{{ $t('app.user.github.tip') }}</p>
        </div>

        <div v-if="!userInfo.github" class="mt-3">
          <!-- connect -->
          <el-button class="w-full" type="primary" size="large" @click="connect">
            {{ $t('app.user.github.conn') }}
          </el-button>
        </div>
        <div v-else class="mt-3">
          <!-- dis -->
          <el-button class="w-full" type="info" :loading="disconnectLoading" size="large" @click="disconnect">
            {{ $t('app.user.github.disconn') }}
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.row {
  margin-top: 1em;
  margin-bottom: 1em;
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.left,
.right {
  display: flex;
  align-items: center;
}

.left {
  justify-content: flex-start;
}

.right {
  flex-grow: 1;
}

.content {
  margin-top: 40px;
}

/* el-upload */
:deep(.content .el-upload) {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

/* el-image */
.content .block {
  padding: 30px 0;
  text-align: center;
  border-right: solid 1px var(--el-border-color);
  display: inline-block;
  width: 49%;
  box-sizing: border-box;
  vertical-align: top;
}

.content .demonstration {
  display: block;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 20px;
}

.content .el-image {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

.content .image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-secondary);
  font-size: 30px;
}
</style>
