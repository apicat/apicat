<script setup lang="ts">
import ProjectExportPage from '@/views/project/ProjectSettingPages/ProjectExportPage.vue'
import { useModal } from '@/hooks'

const { dialogVisible, showModel } = useModal()
const projectId = ref()
const docId = ref()

async function show(project_id?: string, doc_id?: string | number) {
  projectId.value = project_id
  docId.value = doc_id
  showModel()
}

defineExpose({
  show,
})
</script>

<template>
  <el-dialog
    v-model="dialogVisible"
    append-to-body
    :close-on-click-modal="false"
    align-center
    destroy-on-close
    :title="$t('app.project.collection.export.title')"
    width="660px"
  >
    <div class="mt-10px" />
    <ProjectExportPage
      :project_id="projectId"
      :doc_id="docId"
      export-type="collection"
    />
  </el-dialog>
</template>
