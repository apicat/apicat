<template>
  <el-form ref="projectFormRef" :model="form" :rules="rules" label-position="top" class="max-w-sm py-2 pl-2" @submit.prevent="handleSubmit(projectFormRef)">
    <el-form-item :label="$t('app.project.form.title')" prop="title" class="hide_required">
      <el-input v-model="form.title" :placeholder="$t('app.project.form.title')" clearable maxlength="255" />
    </el-form-item>

    <el-form-item :label="$t('app.project.form.desc')">
      <el-input v-model="form.description" :placeholder="$t('app.project.form.desc')" type="textarea" :autosize="{ minRows: 4, maxRows: 4 }" maxlength="255" />
    </el-form-item>

    <el-button type="primary" @click="handleSubmit(projectFormRef)" :loading="isLoading">{{ $t('app.common.save') }}</el-button>
    <el-button type="danger" @click="handleRemove" :loading="isDeleteLoading">{{ $t('app.common.delete') }}</el-button>
  </el-form>
</template>
<script setup lang="tsx">
import { deleleProject, updateProjectBaseInfo } from '@/api/project'
import uesProjectStore from '@/store/project'
import { ProjectInfo } from '@/typings/project'
import { FormInstance } from 'element-plus'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const { projectDetailInfo, setCurrentProjectInfo } = uesProjectStore()
const [isLoading, updateProjectBaseInfoApi] = updateProjectBaseInfo()
const [isDeleteLoading, deleleProjectApi] = deleleProject()

const form: ProjectInfo = reactive({
  id: projectDetailInfo!.id,
  title: projectDetailInfo!.title,
  description: projectDetailInfo!.description,
})

const rules = {
  title: [
    { required: true, message: '请输入项目名称', trigger: 'blur' },
    { min: 2, message: '项目名称不能少于两个字', trigger: 'blur' },
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
  const tip = t('app.common.confirmDelete', { msg: projectDetailInfo!.title })
  AsyncMsgBox({
    title: t('app.common.deleteTip'),
    content: <div class="break-all">{tip}</div>,
    onOk: async () => {
      await deleleProjectApi(projectDetailInfo!.id)
      router.replace('/home')
    },
  })
}
</script>
