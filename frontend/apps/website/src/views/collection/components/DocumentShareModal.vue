<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import useApi from '@/hooks/useApi'
import useClipboard from '@/hooks/useClipboard'
import { Visibility } from '@/commons/constant'
import { COLLECTION_SHARE_PATH_NAME, PROJECT_COLLECTION_PATH_NAME } from '@/router'
import { apiChangeDocShareStatus, apiGetDocShareInfo, apiResetDocShareKey } from '@/api/project/share'

const { t } = useI18n()
const router = useRouter()
const visible = ref<boolean>(false)
const loading = ref<boolean>(false)
const error = ref<string>('')

const defaultBase = { projectID: '', collectionID: -1 }

const collectionShareBase = ref({ ...defaultBase })

const shareStatus = ref<boolean>(false)
const shareInfo = ref<ShareAPI.ResponseDocShareInfo>({
  collectionPublicID: '',
  secretKey: '',
  visibility: Visibility.Private,
})

const link = computed(() => {
  // public collection
  if (shareInfo.value.visibility === Visibility.Public) {
    return (
      location.origin +
      router.resolve({
        name: PROJECT_COLLECTION_PATH_NAME,
        params: {
          project_id: collectionShareBase.value.projectID,
          collectionID: collectionShareBase.value.collectionID,
        },
      }).fullPath
    )
  }

  // private collection
  return (
    location.origin +
    router.resolve({
      name: COLLECTION_SHARE_PATH_NAME,
      params: {
        collectionPublicID: shareInfo.value.collectionPublicID,
      },
    }).fullPath
  )
})
const [statusLoading, _changeStatus] = useApi(apiChangeDocShareStatus)

async function getInfo() {
  shareInfo.value = await apiGetDocShareInfo(
    collectionShareBase.value.projectID,
    collectionShareBase.value.collectionID,
  )
  shareStatus.value = !!shareInfo.value.secretKey
}

async function init() {
  if (
    !collectionShareBase.value.projectID &&
    (!collectionShareBase.value.collectionID || collectionShareBase.value.collectionID < 0)
  )
    return
  loading.value = true
  visible.value = true
  try {
    await getInfo()
    loading.value = false
  } catch (e) {
    if (e && (e as any).message) error.value = (e as any).message
    else error.value = 'UNKNOWN ERROR'
  }
}
function clear() {
  collectionShareBase.value = {
    ...defaultBase,
  }
}

async function changeStatus(status: any) {
  shareStatus.value = !status
  const res = await _changeStatus(collectionShareBase.value.projectID, collectionShareBase.value.collectionID, status)
  shareInfo.value = {
    ...shareInfo.value,
    ...res,
  }
  shareStatus.value = status
}
const [keyLoading, _resetKey] = useApi(apiResetDocShareKey)
async function resetKey() {
  const res: any = await _resetKey(collectionShareBase.value.projectID, collectionShareBase.value.collectionID)
  shareInfo.value = {
    ...shareInfo.value,
    ...res,
  }
}

const copyTextRef = computed(() => {
  const info = unref(shareInfo)
  return info.visibility === Visibility.Public
    ? link.value
    : [
        `${t('app.project.collection.share.copylink')}${link.value}`,
        `${t('app.project.collection.share.copypass')}${info.secretKey}`,
      ].join('\n')
})

const elCopyText = computed(() =>
  unref(shareInfo).visibility === Visibility.Public
    ? t('app.project.collection.share.cl')
    : t('app.project.collection.share.clap'),
)
const { handleCopy, elCopyTextRef } = useClipboard(copyTextRef, elCopyText, t('app.project.collection.share.copied'))

function show(projectID: string, collectionID: number) {
  collectionShareBase.value = {
    projectID,
    collectionID,
  }
}

watch(collectionShareBase, init)

defineExpose({
  show,
})
</script>

<template>
  <el-dialog
    v-model="visible"
    :title="$t('app.project.collection.share.title')"
    :width="540"
    append-to-body
    align-center
    destroy-on-close
    class="fullscree"
    @closed="clear">
    <div v-loading="loading" style="padding: 10px" class="mb-3">
      <div v-if="shareInfo.visibility === Visibility.Public">
        <el-form label-position="top" class="px-6">
          <el-form-item label="">
            <el-input v-model="link" readonly>
              <template #append>
                <el-button type="primary" @click="handleCopy">
                  {{ elCopyTextRef }}
                </el-button>
              </template>
            </el-input>
          </el-form-item>
        </el-form>
      </div>

      <div v-else>
        <div class="flex items-center px-6">
          <div class="flex-1">
            <div class="ivu-list-item-meta-description">
              {{ $t('app.project.collection.share.tip') }}
            </div>
          </div>
          <el-switch
            v-model="shareStatus"
            class="ml-20px"
            :loading="statusLoading"
            inline-prompt
            @change="changeStatus" />
        </div>

        <ElCollapseTransition>
          <div v-if="shareStatus" class="text-center">
            <el-divider />
            <div class="px-6 grid-container">
              <div class="box1">
                {{ $t('app.project.collection.share.link') }}
              </div>
              <div class="box2">
                <el-input v-model="link" readonly />
              </div>
              <div class="box3">
                {{ $t('app.project.collection.share.password') }}
              </div>
              <div class="box4">
                <el-input v-model="shareInfo.secretKey" style="width: auto; max-width: 150px" readonly />
                <el-button class="ml-3" :loading="keyLoading" @click="resetKey">
                  {{ $t('app.project.collection.share.reset') }}
                </el-button>
              </div>
            </div>

            <el-button style="width: 90%; margin: 20px 0px 10px 0px; height: 45px" type="primary" @click="handleCopy">
              {{ elCopyTextRef }}
            </el-button>
          </div>
        </ElCollapseTransition>
      </div>
    </div>
  </el-dialog>
</template>

<style scoped>
.grid-container {
  display: grid;
  /* grid-template-columns: repeat(2, 1fr); */
  grid-template-columns: max-content 1fr;
  /* grid-template-rows: repeat(2, 1fr); */
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.box1,
.box2,
.box3,
.box4 {
  display: flex;
}
</style>
