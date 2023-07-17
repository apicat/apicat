<template>
  <el-dialog v-model="dialogVisible" title="分享文档" :width="540" append-to-body align-center class="fullscree">
    <div v-loading="isLoadingForShareStatus" class="mb-3">
      <div v-if="shareInfo.visibility === CollectionVisibilityEnum.PUBLIC">
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

      <div v-if="shareInfo.visibility === CollectionVisibilityEnum.PRIVATE">
        <div class="flex items-center px-6">
          <div class="flex-1">
            <h4>开启分享</h4>
            <div class="ivu-list-item-meta-description">开启分享后，获得链接的人可以访问接口内容。</div>
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
          <el-form-item label="文档链接">
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
import { getCollectionShareDetail, resetSecretToCollection, switchCollectionShareStatus } from '@/api/shareCollection'
import { CollectionVisibilityEnum } from '@/commons'
import { useModal } from '@/hooks'
import useApi from '@/hooks/useApi'
import useClipboard from '@/hooks/useClipboard'
import { getDocumentPrivateShareLink, getDocumentPublicShareLink } from '@/router/share'
import { isEmpty } from 'lodash-es'

export type CollectionShareDetailParams = {
  project_id: string
  collection_id: string
  collection_public_id?: string
}

type CollectionShareInfo = {
  link: string
  secret_key: string
  visibility: CollectionVisibilityEnum
}

// 当前分享的文档信息
let currentShareDocParams: CollectionShareDetailParams = { project_id: '', collection_id: '' }

// 分享信息
const shareInfo: Ref<CollectionShareInfo> = ref({
  link: '',
  secret_key: '',
  visibility: CollectionVisibilityEnum.PRIVATE,
})

const copyTextRef = computed(() => {
  const info = unref(shareInfo)
  return info.visibility === CollectionVisibilityEnum.PUBLIC ? info.link : [`文档链接：${info.link}`, `密码：${info.secret_key}`].join('\n')
})

const elCopyTextRef = computed(() => (unref(shareInfo).visibility === CollectionVisibilityEnum.PUBLIC ? '复制链接' : '复制链接和密码'))

const { dialogVisible, showModel } = useModal()
const { handleCopy, isCopied, elCopiedText } = useClipboard(copyTextRef)
const [isLoadingForShareStatus, getCollectionShareDetailApi] = useApi(getCollectionShareDetail)
const [isLoadingForResetSecret, resetSecretToCollectionApi] = useApi(resetSecretToCollection)
const [isLoadingForSwitchShareStatus, switchCollectionShareStatusApi] = useApi(switchCollectionShareStatus)

const isShareForSwitchStatus = ref(false)

const fetchCollectionShareDetail = async (params: CollectionShareDetailParams) => {
  const { visibility, collection_public_id, secret_key } = await getCollectionShareDetailApi(params)
  isShareForSwitchStatus.value = !isEmpty(secret_key)
  currentShareDocParams.collection_public_id = collection_public_id
  shareInfo.value = {
    link: getDocumentShareLink(visibility),
    secret_key,
    visibility,
  }
}

const onResetPasswordBtnClick = async () => {
  const { secret_key } = await resetSecretToCollectionApi({ ...currentShareDocParams })
  shareInfo.value.secret_key = secret_key
}

const onShareStatusSwitch = async (share: boolean) => {
  try {
    const { secret_key, collection_public_id } = await switchCollectionShareStatusApi({ ...currentShareDocParams, share: share ? 'open' : 'close' })
    currentShareDocParams.collection_public_id = collection_public_id
    shareInfo.value.link = getDocumentShareLink(shareInfo.value.visibility)
    shareInfo.value.secret_key = secret_key
  } catch (error) {
    isShareForSwitchStatus.value = !isShareForSwitchStatus.value
  }
}

const getDocumentShareLink = (visibility: CollectionVisibilityEnum) => {
  const { project_id, collection_id, collection_public_id = '' } = currentShareDocParams
  return visibility === CollectionVisibilityEnum.PUBLIC ? getDocumentPublicShareLink(project_id, collection_id) : getDocumentPrivateShareLink(collection_public_id)
}

const show = (params: CollectionShareDetailParams) => {
  currentShareDocParams = params
  showModel()
  fetchCollectionShareDetail(params)
}

defineExpose({
  show,
})
</script>
