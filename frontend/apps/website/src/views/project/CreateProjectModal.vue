<template>
  <el-dialog v-model="dialogVisible" :title="$t('app.project.createModal.title')" :width="400" align-center :close-on-click-modal="false">
    <!-- 内容 -->
    <el-form label-position="top" label-width="100px" :model="form" :rules="rules" ref="projectFormRef" @submit.prevent="handleSubmit(projectFormRef)">
      <el-form-item :label="$t('app.project.form.title')" prop="title">
        <el-input v-model="form.title" :placeholder="$t('app.project.form.title')" clearable />
      </el-form-item>
      <div class="my-40px">
        <el-divider>
          <span class="font-400">
            {{ $t('app.project.createModal.dividerLine') }}
          </span>
        </el-divider>
      </div>
      <div :class="ns.b()">
        <div :class="[ns.e('items'), { [ns.is('active')]: selectedProjectType === 'blank' }]" @click="handleSelectedProjectType('blank')">
          <ac-icon-uil-file-blank class="text-30px" />
          <div :class="ns.e('text')">
            {{ $t('app.project.createModal.blackProject') }}
          </div>
        </div>

        <FileUploaderWrapper
          ref="fileUploaderWrapper"
          accept=".json"
          @change="handleFileSelect"
          v-slot="{ fileName }"
          :class="[ns.e('items'), { [ns.is('active')]: selectedProjectType === 'import' }]"
        >
          <div class="flex flex-col items-center w-full" @click="handleSelectedProjectType('import')">
            <ac-icon-lucide-file-text class="text-30px" />
            <div :class="ns.e('text')" class="w-full">
              <p v-if="!fileName">{{ $t('app.project.createModal.importProject') }}</p>
              <p v-else class="truncate">{{ fileName }}</p>
              <p class="text-gray-400 text-12px">
                {{ $t('app.project.createModal.importProjectTip') }}
              </p>
            </div>
          </div>
        </FileUploaderWrapper>
      </div>
    </el-form>
    <!-- 底部按钮 -->
    <div slot="footer" class="text-right mt-20px">
      <el-button @click="dialogVisible = false">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(projectFormRef)">
        {{ $t('app.common.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { useNamespace, useModal } from '@/hooks'
import { useI18n } from 'vue-i18n'
import type { FormInstance, FormRules } from 'element-plus'
import { createProject } from '@/api/project'
import { getProjectDetailPath } from '@/router'
import { ProjectInfo } from '@/typings'
import uesProjectStore from '@/store/project'

const ns = useNamespace('project-types')
const { t } = useI18n()
const router = useRouter()
const projectStore = uesProjectStore()

const projectFormRef = ref<FormInstance>()
const fileUploaderWrapper = ref()
const { dialogVisible, showModel } = useModal(projectFormRef as any)
const [isLoading, createProjectApi] = createProject()

const selectedProjectType = ref('blank')

const form = reactive({
  title: '',
  data: '',
})

watch(dialogVisible, () => {
  if (!dialogVisible.value) {
    fileUploaderWrapper.value.clear()
    form.data = ''
  }
})

const rules = reactive<FormRules>({
  title: [{ required: true, message: t('app.project.rules.title'), trigger: 'change' }],
})

const handleSelectedProjectType = (type: string) => {
  selectedProjectType.value = type

  if (type !== 'import') {
    form.data = ''
  }
}

const handleFileSelect = (file: File) => {
  if (!file) {
    form.data = ''
    return
  }

  const reader = new FileReader()
  reader.onloadend = () => {
    form.data = reader.result as string
  }
  reader.readAsDataURL(file)
}

const handleSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  try {
    const valid = await formEl.validate()
    if (valid) {
      const project: ProjectInfo = await createProjectApi(toRaw(form))
      projectStore.setCurrentProjectInfo(project)
      router.push(getProjectDetailPath(project.id))
    }
  } catch (error) {
    //
  }
}

defineExpose({
  show: showModel,
})
</script>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;

@include b(project-types) {
  @apply flex justify-between gap-30px;

  @include e(items) {
    @apply flex flex-col items-center flex-1 border border-solid rounded cursor-pointer p-20px border-gray-300 hover:border-gray-45;

    @include when('active') {
      @apply border-blue-primary text-blue-primary;
    }
  }

  @include e(text) {
    @apply h-30px mt-20px;
    line-height: 20px;
  }
}
</style>
