<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import type { FormInstance, FormRules } from 'element-plus'
import { FileUploaderWrapper } from '@apicat/components'
import { useNamespace } from '@apicat/hooks'
import { ElMessage } from 'element-plus'
import { useProjectCover } from '../logic/useProjectCover'
import { useProjectListContext } from '../logic/useProjectListContext'
import swaggerLogo from '@/assets/images/logo-swagger@2x.png'
import openApiLogo from '@/assets/images/logo-openapis.svg'
import apiCatLogo from '@/assets/images/logo-square.svg'
import postmanLogo from '@/assets/images/logo-postman@2x.png'
import useProjectGroupStore from '@/store/projectGroup'
import { Visibility } from '@/commons/constant'
import { createProject } from '@/api/project'
import { useTeamStore } from '@/store/team'
import { useGoPage } from '@/hooks/useGoPage'

const emits = defineEmits<{ (e: 'cancel' | 'create-group'): void }>()
const groupStore = useProjectGroupStore()
const teamStore = useTeamStore()
const { groupsForOptions } = storeToRefs(groupStore)
const groups = groupsForOptions
const { t } = useI18n()
const ns = useNamespace('project-types')
const { goProjectDetailPage } = useGoPage()
const { createOrUpdateProjectGroupRef } = useProjectListContext()

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

const defaultForm = {
  teamID: '',
  groupID: 0,
  title: '',
  cover: '',
  data: '',
  type: '',
  visibility: Visibility.Private,
  description: '',
}
const form = reactive<ProjectAPI.RequestCreateProject>({ ...defaultForm })
const fileValid = ref(false)

const { projectCoverBgColorsOptions, projectCoverIcons, bgColorRef, iconRef } = useProjectCover(form)

const rules = reactive<FormRules>({
  title: [
    {
      required: true,
      message: t('app.project.rules.title'),
      trigger: 'change',
    },
  ],
})

function setFileUploaderWrapper(refInstance: any, type: string) {
  if (!fileUploaderWrappers.has(type))
    fileUploaderWrappers.set(type, refInstance)
}

function handleSelectedProjectType(type: string) {
  if (type !== form.type) {
    fileUploaderWrappers.forEach(fileUploaderWrapper => fileUploaderWrapper.clear())
    form.data = ''
  }

  selectedProjectType.value = type
  form.type = type

  if (type === 'blank') {
    form.data = ''
    form.type = ''
  }
}

function handleFileSelect(file: File) {
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

function fileValidate() {
  if (form.type && !form.data) {
    ElMessage.error(t('app.project.createModal.noInputTip'))
    fileValid.value = true
    setTimeout(() => (fileValid.value = false), 500)
    throw new Error('no input')
  }
}

async function handleSubmit(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  try {
    const valid = await formEl.validate()
    fileValidate()
    if (valid) {
      isLoading.value = true
      const res = await createProject(teamStore.currentID, toRaw(form))
      goProjectDetailPage(res.id)
    }
  }
  catch (e) {
    //
  }
  finally {
    isLoading.value = false
  }
}

function setSelectedGroup(group_id: number) {
  form.groupID = group_id
}

watch(
  () => form.groupID,
  (val, oldVal) => {
    if (val === -1) {
      form.groupID = oldVal
      createOrUpdateProjectGroupRef?.value.showWithCallback((group_id: number) => {
        nextTick().then(() => setSelectedGroup(group_id))
      })
    }
  },
)

defineExpose({
  setSelectedGroup,
  reset() {
    projectFormRef.value?.clearValidate()
    projectFormRef.value?.resetFields()
    Object.assign(form, defaultForm)
  },
})
</script>

<template>
  <div class="flex flex-col justify-center mx-auto px-36px">
    <p class="border-b border-solid border-gray-lighter pb-15px text-24px text-gray-title mb-30px">
      {{ $t('app.project.createModal.title') }}
    </p>

    <el-form
      ref="projectFormRef"
      class="m-auto w-776px"
      label-position="top"
      label-width="100px"
      :model="form"
      :rules="rules"
      @submit.prevent="handleSubmit(projectFormRef)"
    >
      <el-form-item v-if="bgColorRef && iconRef" :label="$t('app.project.form.cover')">
        <div class="w-full text-white rounded h-128px flex-center" :style="{ backgroundColor: bgColorRef }">
          <Iconfont :icon="iconRef" :size="55" />
        </div>
      </el-form-item>

      <el-form-item v-if="bgColorRef && iconRef">
        <el-form-item :label="$t('app.project.form.coverColor')" class="flex-1 mr-10px">
          <AcSelect v-model="bgColorRef" class="w-full" :options="projectCoverBgColorsOptions">
            <template #default="{ selected }">
              <div class="flex-center wh-full">
                <span class="inline-flex rounded w-40% h-15px" :style="{ backgroundColor: selected }" />
              </div>
            </template>

            <template #option="{ option }">
              <span :style="{ backgroundColor: option.value }" class="flex h-15px" />
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
          <el-input v-model="form.title" :placeholder="$t('app.project.form.title')" clearable maxlength="255" />
        </el-form-item>

        <el-form-item :label="$t('app.project.form.group')" prop="group_id" class="flex-1 mr-10px">
          <el-select v-model="form.groupID" class="w-full">
            <el-option v-for="item in groups" :key="item.id" :value="item.id!" :label="item.name">
              <div class="truncate max-w-300px">
                {{ item.name }}
              </div>
            </el-option>
            <el-option-group />
            <el-option-group>
              <el-option :label="$t('app.project.groups.createGroup')" :value="-1" />
            </el-option-group>
          </el-select>
        </el-form-item>
      </el-form-item>

      <el-form-item :label="$t('app.project.form.visibility')">
        <el-radio-group v-model="form.visibility">
          <el-radio-button :label="Visibility.Private">
            <el-icon class="mr-1">
              <ac-icon-ep-lock />
            </el-icon>{{ $t('app.common.private') }}
          </el-radio-button>
          <el-radio-button :label="Visibility.Public">
            <el-icon class="mr-1">
              <ac-icon-ep-view />
            </el-icon>{{ $t('app.common.public') }}
          </el-radio-button>
        </el-radio-group>
      </el-form-item>

      <el-form-item :label="$t('app.project.form.desc')">
        <ElInput v-model="form.description" maxlength="255" type="textarea" :autosize="{ minRows: 4 }" />
      </el-form-item>

      <div class="my-40px">
        <el-divider>
          <span class="font-400">
            {{ $t('app.project.createModal.dividerLine') }}
          </span>
        </el-divider>
      </div>
      <div :class="ns.b()">
        <div
          :class="[ns.e('items'), { [ns.is('active')]: selectedProjectType === 'blank' }]"
          class="flex flex-col items-center p-20px"
          @click="handleSelectedProjectType('blank')"
        >
          <ac-icon-uil-file-blank class="text-40px" />
          <div :class="ns.e('text')">
            {{ $t('app.project.createModal.blackProject') }}
          </div>
          <p class="w-full text-center truncate text-gray text-12px" v-html="'&nbsp;'" />
        </div>

        <FileUploaderWrapper
          v-for="item in importTypes"
          :key="item.type"
          :ref="(ref: any) => setFileUploaderWrapper(ref, item.type)"
          accept=".json,.yaml"
          :max-size="20"
          :class="[ns.e('items'), { [ns.is('active')]: selectedProjectType === item.type }]"
          class="transition-all duration-200 ease-in-out"
          :style="
            fileValid && form.type === item.type
              ? { border: '1.5px #f76d6d solid', boxShadow: '0 0 5px rgba(255, 0, 0, 0.3) ' }
              : undefined
          "
          @change="handleFileSelect"
        >
          <template #default="{ fileName }">
            <div
              :key="item.type"
              class="flex-col w-full flex-y-center p-20px"
              @click="handleSelectedProjectType(item.type)"
            >
              <img :src="item.logo">
              <div :class="ns.e('text')" :title="item.name">
                {{ item.name }}
              </div>
              <p
                v-if="form.type === item.type"
                class="w-full text-center truncate text-gray text-12px"
                :title="fileName"
              >
                {{ fileName ? fileName : $t('app.project.createModal.noInput') }}
              </p>
            </div>
          </template>
        </FileUploaderWrapper>
      </div>

      <div class="text-right mt-20px">
        <el-button @click="emits('cancel')">
          {{ $t('app.common.cancel') }}
        </el-button>
        <el-button type="primary" :loading="isLoading" @click="handleSubmit(projectFormRef)">
          {{ $t('app.project.form.create') }}
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;

@include b(project-types) {
  @apply grid grid-cols-5 gap-14px;

  @include e(items) {
    @apply w-full h-full flex border border-solid rounded cursor-pointer border-gray-lighter hover:border-gray-45;

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
