<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import useApi from '@/hooks/useApi'
import useClipboard from '@/hooks/useClipboard'
import { Authority, Visibility } from '@/commons/constant'
import {
  apiChangeProjectShareStatus,
  apiGetProjectShareInfo,
  apiGetProjectShareStatus,
  apiResetProjectShareKey,
} from '@/api/project/share'
import { PROJECT_DETAIL_PATH_NAME } from '@/router'

const { t } = useI18n()
const router = useRouter()

const projectID = ref<string>('')
const shareStatus = ref<boolean>(false)
const shareInfo = ref<ShareAPI.ResponseProjectShareInfo>({
  secretKey: '',
  permission: Authority.None,
  visibility: Visibility.Private,
})
const link = computed(() => {
  return (
    location.origin
    + router.resolve({
      name: PROJECT_DETAIL_PATH_NAME,
      params: {
        project_id: projectID.value,
      },
    }).fullPath
  )
})

// 起步加载 ------------------------------------------------------
watch(projectID, init)
const visible = ref<boolean>(false)
const loading = ref<boolean>(false)
const error = ref<string>('')
async function checkStatus() {
  const res = await apiGetProjectShareStatus(projectID.value)
  shareStatus.value = res.hasShare || false
}
async function getInfo() {
  const res = await apiGetProjectShareInfo(projectID.value)
  shareInfo.value = res
}
async function init() {
  if (!projectID.value)
    return
  loading.value = true
  visible.value = true
  try {
    await Promise.all([checkStatus(), getInfo()])
  }
  catch (e) {
    if (e && (e as any).message)
      error.value = (e as any).message
    else error.value = 'UNKNOWN ERROR'
  }
  finally {
    loading.value = false
  }
}
function clear() {
  projectID.value = ''
}
// 起步加载 ------------------------------------------------------

// 操作 ----------------------------------------------------------
const [statusLoading, _changeStatus] = useApi(apiChangeProjectShareStatus)
async function changeStatus(status: any) {
  shareStatus.value = !status
  const res = await _changeStatus(projectID.value, status)
  shareInfo.value.secretKey = res!.secretKey
  shareStatus.value = status
}
const [keyLoading, _resetKey] = useApi(apiResetProjectShareKey)
async function resetKey() {
  const res = await _resetKey(projectID.value)
  shareInfo.value.secretKey = res!.secretKey
}

const copyTextRef = computed(() => {
  const info = unref(shareInfo)
  return info.visibility === Visibility.Public
    ? link.value
    : [`${t('app.project.share.copylink')}${link.value}`, `${t('app.project.share.copypass')}${info.secretKey}`].join(
        '\n',
      )
})
const elCopyText = computed(() =>
  unref(shareInfo).visibility === Visibility.Public ? t('app.project.share.cl') : t('app.project.share.clap'),
)
const { handleCopy, elCopyTextRef } = useClipboard(copyTextRef, elCopyText, t('app.project.share.copied'))
// 操作 ----------------------------------------------------------

function show(id: string) {
  projectID.value = id
}
defineExpose({
  show,
})
</script>

<template>
  <el-dialog
    v-model="visible"
    :title="$t('app.project.share.title')"
    :width="540"
    append-to-body
    align-center
    destroy-on-close
    class="fullscree"
    @closed="clear"
  >
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
              {{ $t('app.project.share.tip') }}
            </div>
          </div>
          <el-switch
            v-model="shareStatus"
            class="ml-20px"
            :loading="statusLoading"
            inline-prompt
            @change="changeStatus"
          />
        </div>

        <ElCollapseTransition>
          <div v-if="shareStatus" class="text-center">
            <el-divider />

            <!-- <el-form label-position="left" class="px-6"> -->
            <!--   <el-form-item :label="$t('app.project.share.link')"> -->
            <!--     <el-input v-model="link" readonly> </el-input> -->
            <!--   </el-form-item> -->
            <!---->
            <!--   <el-form-item :label="$t('app.project.share.password')"> -->
            <!--     <el-form-item prop="date" class="mr-1"> -->
            <!--       <el-input v-model="shareInfo.secretKey" readonly /> -->
            <!--     </el-form-item> -->
            <!--     <el-button :loading="keyLoading" @click="resetKey"> -->
            <!--       重置密码 -->
            <!--     </el-button> -->
            <!--   </el-form-item> -->
            <!-- </el-form> -->
            <div class="grid-container px-6">
              <div class="box1">
                {{ $t('app.project.share.link') }}
              </div>
              <div class="box2">
                <el-input v-model="link" readonly />
              </div>
              <div class="box3">
                {{ $t('app.project.share.password') }}
              </div>
              <div class="box4">
                <el-input v-model="shareInfo.secretKey" style="width: auto; max-width: 150px" readonly />
                <el-button class="ml-3" :loading="keyLoading" @click="resetKey">
                  {{ $t('app.project.share.reset') }}
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
