<script setup lang="tsx">
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import { useProjectCover } from '../logic/useProjectCover'
import useProjectStore from '@/store/project'
import { Visibility } from '@/commons/constant'
import useApi from '@/hooks/useApi'
import { useTitle } from '@/hooks/useTitle'

const { t } = useI18n()
const projectStore = useProjectStore()
const { project } = projectStore
const { isManager } = storeToRefs(projectStore)
const browserTitle = useTitle()

const [isLoading, updateProjectBaseInfoApi] = useApi(
  projectStore.updateProjectGeneral,
)

const form = reactive<ProjectAPI.ResponseProject>(
  JSON.parse(JSON.stringify(project || {})),
)

const { projectCoverBgColorsOptions, projectCoverIcons, bgColorRef, iconRef } = useProjectCover(form)

const rules = {
  title: [
    { required: true, message: t('app.project.rules.title'), trigger: 'blur' },
    { min: 2, message: t('app.project.rules.titleMinLength'), trigger: 'blur' },
  ],
}

const projectFormRef = shallowRef()

async function handleSubmit(formIns: FormInstance) {
  await formIns.validate()
  const info = toRaw(form)
  await updateProjectBaseInfoApi(info as any)
  browserTitle.value = info.title
  ElMessage.success(t('app.tips.opeartionSuceess'))
}
</script>

<template>
  <el-form
    ref="projectFormRef"
    :model="form"
    :rules="rules"
    label-position="top"
    class="w-full py-2 pl-2"
    @submit.prevent="handleSubmit(projectFormRef)"
  >
    <el-form-item
      v-if="bgColorRef && iconRef"
      v-show="bgColorRef && iconRef"
      :label="$t('app.project.form.cover')"
    >
      <div
        class="w-full text-white rounded h-128px flex-center"
        :style="{ backgroundColor: bgColorRef }"
      >
        <Iconfont :icon="iconRef" :size="55" />
      </div>
    </el-form-item>

    <el-form-item v-if="bgColorRef && iconRef && isManager">
      <el-form-item
        :label="$t('app.project.form.coverColor')"
        class="flex-1 mr-10px"
      >
        <AcSelect
          v-model="bgColorRef"
          class="w-full"
          :options="projectCoverBgColorsOptions"
        >
          <template #default="{ selected }">
            <div class="flex-center wh-full">
              <span
                class="inline-flex rounded w-40% h-15px"
                :style="{ backgroundColor: selected }"
              />
            </div>
          </template>

          <template #option="{ option }">
            <span
              :style="{ backgroundColor: option.value }"
              class="flex h-15px"
            />
          </template>
        </AcSelect>
      </el-form-item>

      <el-form-item
        v-if="isManager"
        :label="$t('app.project.form.coverIcon')"
        class="flex-1 mr-10px"
      >
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

    <el-form-item
      :label="$t('app.project.form.title')"
      prop="title"
      class="hide_required"
    >
      <el-input
        v-model="form.title"
        :disabled="!isManager"
        :placeholder="$t('app.project.form.title')"
        clearable
        maxlength="255"
      />
    </el-form-item>

    <el-form-item :label="$t('app.project.form.visibility')">
      <el-radio-group v-model="form.visibility" :disabled="!isManager">
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
      <el-input
        v-model="form.description"
        :disabled="!isManager"
        :placeholder="$t('app.project.form.desc')"
        type="textarea"
        :autosize="{ minRows: 4, maxRows: 4 }"
        maxlength="255"
      />
    </el-form-item>

    <el-button
      v-if="isManager"
      type="primary"
      :loading="isLoading"
      @click="handleSubmit(projectFormRef)"
    >
      {{ $t('app.project.setting.basic.update') }}
    </el-button>
  </el-form>
</template>

<style lang="scss" scoped>
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
