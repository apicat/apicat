<template>
  <el-dialog v-model="dialogVisible" append-to-body :close-on-click-modal="false" align-center destroy-on-close title="导出文档" width="600px">
    <div class="mt-10px"></div>
    <ProjectExportPage :project_id="projectId" :doc_id="docId" export-type="collection" />
  </el-dialog>
</template>
<script setup lang="ts">
import ProjectExportPage from '@/views/project/ProjectSettingPages/ProjectExportPage.vue'
import { useModal } from '@/hooks'

const { dialogVisible, showModel } = useModal()
const projectId = ref()
const docId = ref()

const show = async (project_id?: string, doc_id?: string) => {
  projectId.value = project_id
  docId.value = doc_id
  showModel()
}

defineExpose({
  show,
})
</script>
