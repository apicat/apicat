<template>
    <main class="ac-preview is-single">
        <div class="ac-preview-content">
            <AcEditorDocument v-if="documentInfo && documentInfo.id" :doc="documentInfo" />
        </div>
    </main>
</template>

<script>
    import AcEditorDocument from './components/AcEditorDocument.vue'
    import { usePreviewStore } from '@/stores/preview'
    import { storeToRefs } from 'pinia'
    import { inject, onMounted } from 'vue'
    import { useRoute } from 'vue-router'
    import { showLoading, hideLoading } from '@/hooks/useLoading'
    export default {
        name: 'DocumentPreview',
        components: {
            AcEditorDocument,
        },

        setup() {
            const showSearchInput = inject('showSearchInput')
            const previewStore = usePreviewStore()
            const { documentInfo } = storeToRefs(previewStore)
            const { params } = useRoute()

            onMounted(async () => {
                showLoading()
                const info = await previewStore.getDocumentInfo(params.doc_id)
                document.title = info.title
                hideLoading()
            })

            showSearchInput(false)
            return {
                documentInfo,
            }
        },
    }
</script>
