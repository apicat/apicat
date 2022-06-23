<template>
    <main :class="layoutClass" v-if="projectInfo">
        <GuestProjectInfoHeader v-if="isGuest" />
        <ProjectInfoHeader v-else />

        <template v-if="hasDocument">
            <DocumentOperateHeader :title="title" v-if="isManager || isDeveloper" />
        </template>

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
    import GuestProjectInfoHeader from '../views/document/components/GuestProjectInfoHeader.vue'
    import DocumentOperateHeader from '../views/document/components/DocumentOperateHeader.vue'
    import DirectoryTree from '../views/document/components/DirectoryTree.vue'
    import DocumentImportModal from '../views/document/components/DocumentImportModal.vue'
    import { ref, provide, computed } from 'vue'
    import { useRouter } from 'vue-router'
    import { storeToRefs } from 'pinia'
    import { useProjectStore } from '@/stores/project'

    const { currentRoute } = useRouter()
    const projectStore = useProjectStore()
    const { isGuest, isReader, isManager, isDeveloper, projectInfo } = storeToRefs(projectStore)

    const projectExportModal = ref()
    const documentImportModal = ref()
    const documentShareModal = ref()
    const directoryTree = ref()
    const title = ref('')
    const layoutClass = computed(() => ['ac-doc-layout', { readonly: isGuest.value || isReader.value }])
    const hasDocument = computed(() => !isNaN(parseInt(currentRoute.value.params.node_id as string, 10)))

    provide('documentShareModal', documentShareModal)
    provide('projectExportModal', projectExportModal)
    provide('documentImportModal', documentImportModal)
    provide('updateTreeNode', (id: any, node: any) => directoryTree.value && directoryTree.value.updateTreeNode(id, node))
    provide('setDocumentTitle', (t: string) => {
        title.value = t
    })
</script>
