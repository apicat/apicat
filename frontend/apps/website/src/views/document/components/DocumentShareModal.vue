<template>
  <el-dialog v-model="dialogVisible" title="分享文档" :width="540" append-to-body align-center class="fullscree hide-header">
    <div v-loading="isLoadingForShareStatus">
      <div v-if="shareInfo.visibility === CollectionVisibilityEnum.PUBLIC">
        <el-form label-position="top" class="px-6 py-3">
          <el-form-item label="">
            <el-input readonly v-model="shareInfo.link">
              <template #append>
                <el-button type="primary" @click="handleCopy">{{ elCopyTextRef }}</el-button>
              </template>
            </el-input>
          </el-form-item>
        </el-form>
      </div>

      <div v-if="shareInfo.visibility === CollectionVisibilityEnum.PRIVATE">
        <div class="flex items-center px-6" :class="{ 'py-3': !isShareForSwitchStatus, 'pt-3': isShareForSwitchStatus }">
          <div class="flex-1">
            <h4>开启分享</h4>
            <div class="ivu-list-item-meta-description">开启分享后，获得链接的人可以访问接口内容。</div>
          </div>
          <el-switch :loading="isLoadingForSwitchShareStatus" v-model="isShareForSwitchStatus" @change="onShareStatusSwitch" inline-prompt active-text="开" inactive-text="关" />
        </div>

        <el-divider v-if="isShareForSwitchStatus" />

        <el-form v-if="isShareForSwitchStatus" label-position="top" class="px-6">
          <el-form-item label="文档链接">
            <el-input readonly v-model="shareInfo.link">
              <template #append>
                <el-button type="primary" @click="handleCopy">{{ elCopyTextRef }}</el-button>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="密码">
            <el-form-item prop="date" class="mr-1">
              <el-input readonly v-model="shareInfo.password" />
            </el-form-item>
            <el-button :loading="isLoadingForResetSecret" @click="onResetPasswordBtnClick">重置密码</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { getCollectionShareDetail, resetSecretToCollection } from '@/api/collection'
import { CollectionVisibilityEnum } from '@/commons'
import { useModal } from '@/hooks'
import useApi from '@/hooks/useApi'
import useClipboard from '@/hooks/useClipboard'

// 分享信息
const shareInfo: Ref<{ link: string; password: string; visibility: CollectionVisibilityEnum }> = ref({
  link: '',
  password: '',
  visibility: CollectionVisibilityEnum.PUBLIC,
})

const copyTextRef = computed(() => {
  const info = unref(shareInfo)
  return info.visibility === CollectionVisibilityEnum.PUBLIC ? info.link : [`文档链接：${info.link}`, `密码：${info.password}`].join('\n')
})

const { dialogVisible, showModel, hideModel } = useModal()
const { handleCopy, elCopyTextRef } = useClipboard(copyTextRef, '复制链接', '复制链接和密码')

const [isLoadingForShareStatus, getCollectionShareDetailApi] = useApi(getCollectionShareDetail)
const [isLoadingForResetSecret, resetSecretToCollectionApi] = useApi(resetSecretToCollection)
const [isLoadingForSwitchShareStatus, resetSecretToCollectionApi2] = useApi(resetSecretToCollection)

const isShareForSwitchStatus = ref(false)

const onResetPasswordBtnClick = async () => {}

const onShareStatusSwitch = async () => {}
</script>
