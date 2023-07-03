<template>
  <el-dialog v-model="dialogVisible" title="分享项目" :width="540" append-to-body align-center class="fullscree">
    <div v-loading="isLoadingForShareStatus" class="mb-3">
      <div v-if="shareInfo.visibility === ProjectVisibilityEnum.PUBLIC">
        <el-form label-position="top" class="px-6">
          <el-form-item label="">
            <el-input readonly v-model="shareInfo.link">
              <template #append>
                <el-button type="primary" @click="handleCopy">{{ isCopied ? elCopiedText : elCopyTextRef }}</el-button>
              </template>
            </el-input>
          </el-form-item>
        </el-form>
      </div>

      <div v-if="shareInfo.visibility === ProjectVisibilityEnum.PRIVATE">
        <div class="flex items-center px-6">
          <div class="flex-1">
            <h4>开启分享</h4>
            <div class="ivu-list-item-meta-description">开启分享后，获得链接的人可以访问项目内容。</div>
          </div>
          <el-switch
            :loading="isLoadingForSwitchShareStatus"
            v-model="isShareForSwitchStatus"
            @change="(v:any)=>onShareStatusSwitch(v)"
            inline-prompt
            :active-text="$t('app.common.on')"
            :inactive-text="$t('app.common.off')"
          />
        </div>

        <el-divider v-if="isShareForSwitchStatus" />

        <el-form v-if="isShareForSwitchStatus" label-position="top" class="px-6">
          <el-form-item label="项目链接">
            <el-input readonly v-model="shareInfo.link">
              <template #append>
                <el-button type="primary" @click="handleCopy">{{ isCopied ? elCopiedText : elCopyTextRef }}</el-button>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="密码">
            <el-form-item prop="date" class="mr-1">
              <el-input readonly v-model="shareInfo.secret_key" />
            </el-form-item>
            <el-button :loading="isLoadingForResetSecret" @click="onResetPasswordBtnClick">重置密码</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { getProjectShareDetail, resetSecretToProject, switchProjectShareStatus } from '@/api/project'
import { ProjectVisibilityEnum } from '@/commons'
import { useModal } from '@/hooks'
import useApi from '@/hooks/useApi'
import useClipboard from '@/hooks/useClipboard'
import { getProjectShareLink } from '@/router/share'
import { isEmpty } from 'lodash-es'

export type ProjectShareDetailParams = {
  project_id: string
}

type ProjectShareInfo = {
  link: string
  secret_key: string
  visibility: ProjectVisibilityEnum
}

// 当前分享的文档信息
let currentShareProjectParams: Record<string, any> = {}

// 分享信息
const shareInfo: Ref<ProjectShareInfo> = ref({
  link: '',
  secret_key: '',
  visibility: ProjectVisibilityEnum.PRIVATE,
})

const copyTextRef = computed(() => {
  const info = unref(shareInfo)
  return info.visibility === ProjectVisibilityEnum.PUBLIC ? info.link : [`项目链接：${info.link}`, `密码：${info.secret_key}`].join('\n')
})

const elCopyTextRef = computed(() => (unref(shareInfo).visibility === ProjectVisibilityEnum.PUBLIC ? '复制链接' : '复制链接和密码'))

const { dialogVisible, showModel } = useModal()
const { handleCopy, isCopied, elCopiedText } = useClipboard(copyTextRef)
const [isLoadingForShareStatus, getProjectShareDetailApi] = useApi(getProjectShareDetail)
const [isLoadingForResetSecret, resetSecretToProjectApi] = useApi(resetSecretToProject)
const [isLoadingForSwitchShareStatus, switchProjectShareStatusApi] = useApi(switchProjectShareStatus)

const isShareForSwitchStatus = ref(false)

const fetchProjectShareDetail = async (params: ProjectShareDetailParams) => {
  const { visibility, secret_key } = await getProjectShareDetailApi(params.project_id)
  shareInfo.value = {
    link: getProjectShareLink(params.project_id),
    secret_key,
    visibility,
  }
  isShareForSwitchStatus.value = !isEmpty(secret_key)
}

const onResetPasswordBtnClick = async () => {
  const { secret_key } = await resetSecretToProjectApi({ ...currentShareProjectParams })
  shareInfo.value.secret_key = secret_key
}

const onShareStatusSwitch = async (share: boolean) => {
  try {
    const { secret_key } = await switchProjectShareStatusApi({ ...currentShareProjectParams, share: share ? 'open' : 'close' })
    shareInfo.value.secret_key = secret_key
  } catch (error) {
    isShareForSwitchStatus.value = !isShareForSwitchStatus.value
  }
}

const show = (params: ProjectShareDetailParams) => {
  currentShareProjectParams = params
  showModel()
  fetchProjectShareDetail(params)
}

defineExpose({
  show,
})
</script>
