<template>
    <main class="ac-preview is-single" v-loading="isLoading">
        <div class="ac-preview-content">
            <AcEditorDocument :doc="documentInfo" />
        </div>
    </main>
</template>

<script>
    import AcEditorDocument from './components/AcEditorDocument.vue'
    import { getTrashNormalDocumentDetail } from '@/api/preview'
    import { inject } from 'vue'

    export default {
        name: 'TrashDocumentPreview',
        components: {
            AcEditorDocument,
        },

        data() {
            return {
                isLoading: false,
                documentInfo: {},
            }
        },

        setup() {
            const showSearchInput = inject('showSearchInput')
            showSearchInput(false)
        },

        async mounted() {
            let project_id = this.$route.params.project_id || ''
            let document_id = this.$route.params.doc_id || ''
            this.isLoading = true
            getTrashNormalDocumentDetail(project_id, document_id)
                .then((res) => {
                    res.data.title = res.data.title || res.data.name
                    this.documentInfo = res.data || {}
                    document.title = this.documentInfo.title || 'ApiCat 文档预览'
                })
                .catch((e) => {
                    //
                })
                .finally(() => {
                    this.isLoading = false
                })
        },
    }
</script>
