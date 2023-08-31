<template>
  <el-form class="m-auto px-36px w-776px" label-position="top" label-width="100px" :model="form" :rules="rules" ref="projectFormRef" @submit.prevent="handleSubmit(projectFormRef)">
    <el-form-item v-if="bgColorRef && iconRef" :label="$t('app.project.form.cover')" v-show="bgColorRef && iconRef">
      <div class="w-full text-white rounded h-128px flex-center" :style="{ backgroundColor: bgColorRef }">
        <Iconfont :icon="iconRef" :size="55" />
      </div>
    </el-form-item>

    <el-form-item v-if="bgColorRef && iconRef">
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

      <el-form-item :label="$t('app.project.form.coverIcon')" class="flex-1 mr-10px">
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

    <el-form-item>
      <el-form-item :label="$t('app.project.form.title')" prop="title" class="flex-1 mr-10px">
        <el-input v-model="form.title" :placeholder="$t('app.project.form.title')" clearable />
      </el-form-item>

      <el-form-item label="项目分组" prop="group_id" class="flex-1 mr-10px">
        <el-select v-model="form.group_id" class="w-full">
          <el-option v-for="item in groups" :key="item.id" :label="item.name" :value="item.id!" />
        </el-select>
      </el-form-item>
    </el-form-item>

    <el-form-item :label="$t('app.project.form.visibility')">
      <el-radio-group v-model="form.visibility">
        <el-radio-button :label="ProjectVisibilityEnum.PRIVATE">
          <el-icon class="mr-1"><ac-icon-ep-lock /></el-icon>{{ $t('app.common.private') }}
        </el-radio-button>
        <el-radio-button :label="ProjectVisibilityEnum.PUBLIC">
          <el-icon class="mr-1"><ac-icon-ep-view /></el-icon>{{ $t('app.common.public') }}
        </el-radio-button>
      </el-radio-group>
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
        <ac-icon-uil-file-blank class="text-40px" />
        <div :class="ns.e('text')">
          {{ $t('app.project.createModal.blackProject') }}
        </div>
      </div>

      <FileUploaderWrapper
        :class="[ns.e('items'), { [ns.is('active')]: selectedProjectType === item.type }]"
        :ref="(ref:any)=>setFileUploaderWrapper(ref, item.type)"
        v-for="item in importTypes"
        accept=".json,.yaml"
        @change="handleFileSelect"
        v-slot="{ fileName }"
      >
        <div :key="item.type" class="flex flex-col w-full flex-y-center" @click="handleSelectedProjectType(item.type)">
          <img :src="item.logo" />
          <div :class="ns.e('text')" :title="!!fileName ? fileName : item.name">
            {{ !!fileName ? fileName : item.name }}
          </div>
        </div>
      </FileUploaderWrapper>
    </div>

    <div class="text-right mt-20px">
      <el-button @click="emits('cancel')">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(projectFormRef)">
        {{ $t('app.common.confirm') }}
      </el-button>
    </div>
  </el-form>
</template>

<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { useI18n } from 'vue-i18n'
import type { FormInstance, FormRules } from 'element-plus'
import { createProject } from '@/api/project'
import { getProjectDetailPath } from '@/router'
import { ProjectGroup, ProjectInfo } from '@/typings'
import { useProjectCover } from '../logic/useProjectCover'
import { ProjectVisibilityEnum } from '@/commons'
import swaggerLogo from '@/assets/images/logo-swagger@2x.png'
import openApiLogo from '@/assets/images/logo-openapis.svg'
import apiCatLogo from '@/assets/images/logo-square.svg'
import postmanLogo from '@/assets/images/logo-postman@2x.png'

const emits = defineEmits<{ (e: 'cancel'): void }>()
const props = withDefaults(defineProps<{ group_id: number; groups: ProjectGroup[] }>(), {
  group_id: 0,
  groups: () => [],
})

const { t } = useI18n()
const ns = useNamespace('project-types')
const router = useRouter()

const projectFormRef = ref<FormInstance>()
const fileUploaderWrappers = new Map()
const isLoading = ref(false)

const selectedProjectType = ref('blank')

const importTypes = [
  { type: 'apicat', name: 'ApiCat', logo: apiCatLogo },
  { type: 'openapi', name: 'OpenAPI', logo: openApiLogo },
  { type: 'swagger', name: 'Swagger', logo: swaggerLogo },
  { type: 'postman', name: 'Postman', logo: postmanLogo },
]

const setFileUploaderWrapper = (refInstance: any, type: string) => {
  if (!fileUploaderWrappers.has(type)) {
    fileUploaderWrappers.set(type, refInstance)
  }
}

const form = reactive({
  group_id: props.group_id,
  title: '',
  cover: '',
  data: '',
  data_type: '',
  visibility: ProjectVisibilityEnum.PRIVATE,
})

const { projectCoverBgColorsOptions, projectCoverIcons, bgColorRef, iconRef } = useProjectCover(form)

const rules = reactive<FormRules>({
  title: [{ required: true, message: t('app.project.rules.title'), trigger: 'change' }],
})

const handleSelectedProjectType = (type: string) => {
  fileUploaderWrappers.forEach((fileUploaderWrapper) => {
    fileUploaderWrapper.clear()
  })

  selectedProjectType.value = type
  form.data_type = type

  if (type === 'blank') {
    form.data = ''
    form.data_type = ''
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
      isLoading.value = true
      const project: ProjectInfo = await createProject(toRaw(form))
      await router.push(getProjectDetailPath(project.id))
    }
  } catch (error) {
    isLoading.value = false
  }
}
</script>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;

@include b(project-types) {
  @apply flex justify-between gap-14px;

  @include e(items) {
    @apply w-full flex flex-col items-center flex-1 border border-solid rounded cursor-pointer p-20px border-gray-lighter hover:border-gray-45;

    img {
      width: 48px;
    }

    @include when('active') {
      @apply border-blue-primary text-blue-primary;
    }
  }

  @include e(text) {
    @apply h-30px mt-20px truncate w-full text-center;
    line-height: 20px;
  }
}
</style>
