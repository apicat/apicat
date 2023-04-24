<template>
  <el-dialog v-model="dialogVisible" center append-to-body :close-on-click-modal="false" :close-on-press-escape="false" destroy-on-close title="公共响应" width="40%">
    <template #header>
      <div class="flex-y-center">
        <el-icon class="mr-5px"><ac-icon-bi:robot /></el-icon>AI生成接口
      </div>
    </template>
    <el-input size="large" v-model="promptText" placeholder="请输入您想生成的接口名称" clearable />
    <el-button :loading="isLoading" type="primary" class="mt-20px" @click="onCreateClick">生成</el-button>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import { createCollectionByAI } from '@/api/collection'
import useApi from '@/hooks/useApi'
import { useParams } from '@/hooks/useParams'
const emits = defineEmits(['ok'])
const { dialogVisible, showModel, hideModel } = useModal()
const [isLoading, createCollectionByAIApi] = useApi(createCollectionByAI)()
const promptText = ref('')
const { project_id } = useParams()

const onCreateClick = async () => {
  try {
    const data = await createCollectionByAIApi({ project_id, title: promptText.value })
    emits('ok', data.id)
    hideModel()
  } catch (error) {
    //
  }
}

defineExpose({
  show: showModel,
})
</script>
