<script setup lang="ts">
import type { FormInstance } from 'element-plus'
import { storeToRefs } from 'pinia'
import useProjectGroupStore from '@/store/projectGroup'
import { useModal } from '@/hooks'
import { apiChangeProjectGroup } from '@/api/project/index'

interface FormState {
  id: number
  projectId: string
}

const emits = defineEmits(['success'])

const form = reactive<FormState>({
  id: 0,
  projectId: '',
})

const formRef = ref<FormInstance>()
const isLoading = ref(false)
const { dialogVisible, showModel, hideModel } = useModal(formRef as any)
const projectGroupStore = useProjectGroupStore()
const { groupsForOptions } = storeToRefs(projectGroupStore)

async function handleSubmit(formEl: FormInstance | undefined) {
  if (!formEl)
    return

  try {
    const valid = await formEl.validate()
    if (valid) {
      isLoading.value = true
      const { projectId, id } = form
      await apiChangeProjectGroup(projectId, { groupID: id })
      emits('success')
      hideModel()
    }
  }
  catch (error) {
    //
  }
  finally {
    isLoading.value = false
  }
}

async function show(projectInfo: ProjectAPI.ResponseProject) {
  showModel()
  await nextTick()
  formRef.value?.clearValidate()
  Object.assign(form, {
    id: projectInfo.selfMember.groupID || 0,
    projectId: projectInfo.id,
  })
}

defineExpose({
  show,
})
</script>

<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('app.project.groups.grouping')"
    :width="348"
    align-center
    :close-on-click-modal="false"
  >
    <!-- 内容 -->
    <el-form
      ref="formRef"
      label-position="top"
      label-width="100px"
      :model="form"
      @submit.prevent="handleSubmit(formRef)"
    >
      <el-form-item label="" prop="id">
        <el-select v-model="form.id" class="w-full">
          <el-option v-for="item in groupsForOptions" :key="item.id" :label="item.name" :value="item.id!">
            <div class="truncate max-w-260px">
              {{ item.name }}
            </div>
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <!-- 底部按钮 -->
    <div class="text-right -mb-10px">
      <el-button @click="dialogVisible = false">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(formRef)">
        {{ $t('app.common.change') }}
      </el-button>
    </div>
  </el-dialog>
</template>
