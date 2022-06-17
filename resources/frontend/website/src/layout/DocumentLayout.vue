<template>
    <DocumentVerticalHeader />

    <main class="ac-doc-layout">
        <div class="ac-doc-layout--left">
            <DirectoryTree ref="directoryTree" />
        </div>
        <div class="ac-doc-layout--right">
            <router-view />
        </div>
    </main>

    <ProjectExportModal ref="projectExportModal" />
    <DocumentImportModal ref="documentImportModal" />
</template>
<script setup lang="ts">
    import DocumentVerticalHeader from '../views/document/components/DocumentVerticalHeader.vue'
    import DirectoryTree from '../views/document/components/DirectoryTree.vue'
    import DocumentImportModal from '../views/document/components/DocumentImportModal.vue'
    import { ref, provide } from 'vue'

    const projectExportModal = ref()
    const documentImportModal = ref()
    const directoryTree = ref()

    provide('projectExportModal', projectExportModal)
    provide('documentImportModal', documentImportModal)
    provide('updateTreeNode', (id: any, node: any) => directoryTree.value && directoryTree.value.updateTreeNode(id, node))
</script>
