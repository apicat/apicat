<template>
    <main class="ac-preview is-single">
        <div class="ac-preview-content">
            <AcEditorDocument :doc="documentInfo" />
        </div>
    </main>
</template>

<script>
    import AcEditorDocument from './components/AcEditorDocument.vue'
    import { usePreviewStore } from '@/stores/preview'
    import { storeToRefs } from 'pinia'
    import { inject } from 'vue'

    export default {
        name: 'DocumentPreview',
        components: {
            AcEditorDocument,
        },

        setup() {
            const showSearchInput = inject('showSearchInput')
            const previewStore = usePreviewStore()
            const { documentInfo } = storeToRefs(previewStore)

            showSearchInput(false)
            document.title = documentInfo.value.title

            return {
                documentInfo,
            }
        },
    }
</script>
