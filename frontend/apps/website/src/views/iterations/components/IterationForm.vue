<template>
  <div v-loading="isLoading" :class="ns.b()">
    <h3 :class="ns.e('title')">{{ isEditMode ? '编辑迭代' : '创建迭代' }}</h3>
    <el-form :model="iterationInfo" ref="iterationFormRef" :rules="iterationRules" label-position="top">
      <el-form-item label="迭代名称" prop="title">
        <el-input v-model="iterationInfo.title" placeholder="请输入迭代名称" />
      </el-form-item>
      <el-form-item label="所属项目" prop="project_id">
        <el-select class="w-full" v-model="iterationInfo.project_id" placeholder="请选择所属项目">
          <el-option v-for="item in projects" :key="item.id" :label="item.title" :value="item.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="迭代描述" prop="description">
        <el-input type="textarea" :rows="4" v-model="iterationInfo.description" placeholder="请输入迭代描述" />
      </el-form-item>
      <el-form-item>
        <div class="flex-1 text-right">
          <el-button @click="handleCancel">{{ $t('app.common.cancel') }}</el-button>
          <el-button type="primary" :loading="isLoadingForSubmit" @click="handleSubmit(iterationFormRef)">{{ $t('app.common.confirm') }}</el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>
<script setup lang="ts">
import { createIteration, getIterationDetail, updateIteration } from '@/api/iteration'
import { useNamespace } from '@/hooks'
import useApi from '@/hooks/useApi'
import { EmptyStruct, Iteration, ProjectInfo } from '@/typings'
import { FormInstance } from 'element-plus'

const props = withDefaults(defineProps<{ id: number | null; projects: ProjectInfo[] }>(), { id: null })
const emits = defineEmits(['success', 'cancel'])

const { id } = toRefs(props)
const ns = useNamespace('iteration-detail')
const iterationFormRef = shallowRef()
const [isLoading, getIterationDetailApi] = useApi(getIterationDetail)
const iterationInfo = ref<EmptyStruct<Iteration>>({})

const iterationRules = {
  title: [{ required: true, message: '请输入迭代名称' }],
  project_id: [{ required: true, message: '请选择所属项目' }],
  description: [{ message: '请输入迭代描述' }],
  collection_ids: [{ required: true, message: '请选择本次迭代所涉及的 API' }],
}

const isEditMode = computed(() => id.value !== null)
const [isLoadingForSubmit, createOrUpdateIterationApi] = useApi(isEditMode.value ? updateIteration : createIteration)
const resetIterationInfo = () => {
  iterationInfo.value = {}
}

const handleSubmit = async (formIns: FormInstance) => {
  try {
    await formIns.validate()
    console.log(toRaw(unref(iterationInfo)))
    await createOrUpdateIterationApi(toRaw(unref(iterationInfo)))
    resetIterationInfo()
    emits('success')
  } catch (error) {
    //
  }
}

const handleCancel = () => emits('cancel')

/**
 * get detail
 */
watch(id, async () => {
  if (!id) {
    resetIterationInfo()
    return
  }

  try {
    iterationInfo.value = await getIterationDetailApi({ iteration_public_id: unref(id) })
  } catch (error) {
    resetIterationInfo()
  }
})
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;
@include b(iteration-detail) {
  @include e(title) {
    @apply text-18px  text-gray-title mb-30px;
  }
}
</style>
