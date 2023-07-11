<template>
  <el-form ref="projectFormRef" :model="form" :rules="rules" label-position="top" class="max-w-sm py-2 pl-2" @submit.prevent="handleSubmit(projectFormRef)">
    <el-form-item v-if="bgColorRef && iconRef" :label="$t('app.project.form.cover')" v-show="bgColorRef && iconRef">
      <div class="w-full text-white rounded h-128px flex-center" :style="{ backgroundColor: bgColorRef }">
        <Iconfont :icon="iconRef" :size="55" />
      </div>
    </el-form-item>

    <el-form-item v-if="bgColorRef && iconRef && isManager">
      <el-form-item :label="$t('app.project.form.coverColor')" class="flex-1 mr-10px">
        <AcSelect v-model="bgColorRef" class="w-full" :options="projectCoverBgColorsOptions">
          <template #default="{ selected }">
            <div class="flex-center wh-full">
              <span class="inline-flex rounded w-40% h-15px" :style="{ backgroundColor: selected }"></span>
            </div>
          </template>

          <template #option="{ option }">
            <span :style="{ backgroundColor: option.value }" class="flex h-15px"></span>
          </template>
        </AcSelect>
      </el-form-item>

      <el-form-item v-if="isManager" :label="$t('app.project.form.coverIcon')" class="flex-1 mr-10px">
        <AcSelect v-model="iconRef" class="w-full" :options="projectCoverIcons">
          <template #default="{ selected }">
            <div class="flex-center wh-full">
              <Iconfont :icon="selected" />
            </div>
          </template>

          <template #option="{ option }">
            <div class="flex-center wh-full">
              <Iconfont :icon="option.value" />
            </div>
          </template>
        </AcSelect>
      </el-form-item>
    </el-form-item>

    <el-form-item :label="$t('app.project.form.title')" prop="title" class="hide_required">
      <el-input :disabled="!isManager" v-model="form.title" :placeholder="$t('app.project.form.title')" clearable maxlength="255" />
    </el-form-item>

    <el-form-item :label="$t('app.project.form.visibility')">
      <el-radio-group v-model="form.visibility" :disabled="!isManager">
        <el-radio-button :label="ProjectVisibilityEnum.PRIVATE">
          <el-icon class="mr-1"><ac-icon-ep-lock /></el-icon>{{ $t('app.common.private') }}
        </el-radio-button>
        <el-radio-button :label="ProjectVisibilityEnum.PUBLIC">
          <el-icon class="mr-1"><ac-icon-ep-view /></el-icon>{{ $t('app.common.public') }}
        </el-radio-button>
      </el-radio-group>
    </el-form-item>

    <el-form-item :label="$t('app.project.form.desc')">
      <el-input
        :disabled="!isManager"
        v-model="form.description"
        :placeholder="$t('app.project.form.desc')"
        type="textarea"
        :autosize="{ minRows: 4, maxRows: 4 }"
        maxlength="255"
      />
    </el-form-item>

    <el-button v-if="isManager" type="primary" @click="handleSubmit(projectFormRef)" :loading="isLoading">{{ $t('app.common.save') }}</el-button>
  </el-form>
  <el-button v-if="isManager" class="absolute bottom-42px right-30px" type="danger" link @click="handleRemove">{{ $t('app.project.setting.deleteProject') }}</el-button>
</template>
<script setup lang="tsx">
import { deleleProject, updateProjectBaseInfo } from '@/api/project'
import useProjectStore from '@/store/project'
import { ProjectInfo } from '@/typings/project'
import { FormInstance } from 'element-plus'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { useI18n } from 'vue-i18n'
import { useProjectCover } from '../logic/useProjectCover'
import { storeToRefs } from 'pinia'
import { ProjectVisibilityEnum } from '@/commons'

const { t } = useI18n()
const router = useRouter()
const projectStore = useProjectStore()
const { projectDetailInfo, setCurrentProjectInfo } = projectStore
const { isManager } = storeToRefs(projectStore)

const [isLoading, updateProjectBaseInfoApi] = updateProjectBaseInfo()
const [isDeleteLoading, deleleProjectApi] = deleleProject()

const form: ProjectInfo = reactive({
  id: projectDetailInfo!.id,
  title: projectDetailInfo!.title,
  cover: projectDetailInfo!.cover,
  description: projectDetailInfo!.description,
  visibility: projectDetailInfo!.visibility,
})

const { projectCoverBgColorsOptions, projectCoverIcons, bgColorRef, iconRef } = useProjectCover(form)

const rules = {
  title: [
    { required: true, message: t('app.project.rules.title'), trigger: 'blur' },
    { min: 2, message: t('app.project.rules.titleMinLength'), trigger: 'blur' },
  ],
}

const projectFormRef = shallowRef()

const handleSubmit = async (formIns: FormInstance) => {
  try {
    const valid = await formIns.validate()
    if (valid) {
      const info = toRaw(form)
      await updateProjectBaseInfoApi(info)
      setCurrentProjectInfo(info)
    }
  } catch (error) {}
}

const handleRemove = () => {
  AsyncMsgBox({
    title: t('app.common.deleteTip'),
    content: <div class="break-all" v-html={t('app.project.setting.deleteProjectTip')}></div>,
    onOk: async () => {
      await deleleProjectApi(projectDetailInfo!.id)
      router.replace('/main')
    },
  })
}
</script>
