<script setup lang="tsx">
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import useProjectStore from '@/store/project'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { MAIN_PATH, MAIN_PATH_NAME } from '@/router'
import { apiQuitProject } from '@/api/project/index'

const { t } = useI18n()
const router = useRouter()
const projectStore = useProjectStore()
const { isManager } = storeToRefs(projectStore)

function handleRemove() {
  AsyncMsgBox({
    confirmButtonText: t('app.common.delete'),
    cancelButtonText: t('app.common.cancel'),
    confirmButtonClass: 'red',
    title: t('app.project.setting.delete.btn'),
    content: <div class="break-all" v-html={t('app.project.setting.delete.pop.tip')}></div>,
    onOk: async () => {
      await projectStore.deleteCurrentProject()
      router.push({
        name: MAIN_PATH_NAME,
      })
    },
  })
}

function handleQuit() {
  AsyncMsgBox({
    cancelButtonText: t('app.common.cancel'),
    confirmButtonText: t('app.project.setting.quit.pop.btn'),
    confirmButtonClass: 'red',
    title: t('app.project.setting.quit.btn'),
    content: <div class="break-all" v-html={t('app.project.setting.quit.pop.tip')}></div>,
    onOk: async () => {
      await apiQuitProject(projectStore.projectID!)
      ElMessage.success(t('app.project.setting.quit.success'))
      router.push(MAIN_PATH).then(() => {
        projectStore.clearProject()
      })
    },
  })
}
</script>

<template>
  <div>
    <!-- delete -->
    <div class="mt-3">
      <div v-if="isManager">
        <p>{{ $t('app.project.setting.delete.tip') }}</p>
        <ElButton class="mt-4 red-outline" @click="handleRemove">
          {{ $t('app.project.setting.delete.btn') }}
        </ElButton>
      </div>

      <div v-else>
        <p>{{ $t('app.project.setting.quit.tip') }}</p>
        <ElButton class="mt-4 red-outline" @click="handleQuit">
          {{ $t('app.project.setting.quit.btn') }}
        </ElButton>
      </div>
    </div>

    <!-- quit -->
  </div>
</template>

<style scoped>
h2,
h3 {
  font-weight: 700;
}

.content {
  margin-top: 15px;
  padding: 17px;
  border: 1px solid red;
  border-radius: 5px;

  p {
    margin-top: 7px;
  }
}

.row {
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
  /* justify-content: flex-start; */
  flex-grow: 1;
}

.right {
  justify-content: flex-end;
  /* flex-grow: 1; */
}
</style>
