<template>
    <main class="ac-doc-layout">
        <DocumentOperateHeader />
        <ProjectInfoHeader />
        <div class="ac-doc-layout__left">
            <DirectoryTree ref="directoryTree" />
        </div>
        <div class="ac-doc-layout__right scroll-content">
            <router-view />
        </div>
    </main>

    <ProjectExportModal ref="projectExportModal" />
    <DocumentImportModal ref="documentImportModal" />
    <DocumentShareModal ref="documentShareModal" />
</template>
<script setup lang="ts">
    import DocumentShareModal from '../views/document/components/DocumentShareModal.vue'
    import ProjectInfoHeader from '../views/document/components/ProjectInfoHeader.vue'
    import DocumentOperateHeader from '../views/document/components/DocumentOperateHeader.vue'
    import DirectoryTree from '../views/document/components/DirectoryTree.vue'
    import DocumentImportModal from '../views/document/components/DocumentImportModal.vue'
    import { ref, provide } from 'vue'

    const projectExportModal = ref()
    const documentImportModal = ref()
    const documentShareModal = ref()
    const directoryTree = ref()

    provide('documentShareModal', documentShareModal)
    provide('projectExportModal', projectExportModal)
    provide('documentImportModal', documentImportModal)
    provide('updateTreeNode', (id: any, node: any) => directoryTree.value && directoryTree.value.updateTreeNode(id, node))
</script>
