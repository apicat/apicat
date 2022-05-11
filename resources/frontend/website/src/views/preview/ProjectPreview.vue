<template>
    <main class="ac-preview" :class="{ 'open-sidebar': isShowCatalog }">
        <aside class="sidebar">
            <div class="sidebar-body">
                <DocumentCatalog :project="projectInfo" :list="catalogList" />
            </div>
            <div class="sidebar-dragbar" @click="isShowCatalog = !isShowCatalog"></div>
            <div class="sidebar-pin-wrapper">
                <div class="pin" @click="isShowCatalog = !isShowCatalog">
                    <div class="pin-inner">
                        <i class="iconfont iconback"></i>
                    </div>
                </div>
            </div>
        </aside>

        <!-- 文档详情 -->
        <div class="ac-preview-content" v-loading="isLoading">
            <AcEditorDocument :doc="documentInfo" v-if="catalogList.length" />
            <div class="text-center pt-16" v-if="!catalogList.length">
                <div class="result">
                    <div class="result-img">
                        <img src="@/assets/image/img-empty.png" alt="" />
                    </div>
                    <p class="mt-6 result-title">这是一篇荒芜之地，还未创建过一篇文档...</p>
                </div>
            </div>
        </div>
    </main>
</template>

<script>
    import DocumentCatalog from './components/DocumentCatalog.vue'
    import AcEditorDocument from './components/AcEditorDocument.vue'
    import { getApiDocumentDetail, getProjectCatalog } from '@/api/preview'
    import { Storage } from '@ac/shared'
    import { mapState } from 'pinia'
    import { usePreviewStore } from '@/stores/preview'
    import { showLoading, hideLoading } from '@/hooks/useLoading'

    export default {
        name: 'ProjectPreview',
        components: {
            DocumentCatalog,
            AcEditorDocument,
        },

        data() {
            return {
                isLoading: false,
                isShowCatalog: true,
                navigation: {},
                documentInfo: {},
                catalogList: [],
                project_id: this.$route.params.project_id,
                token: Storage.get(Storage.KEYS.SECRET_PROJECT_TOKEN + this.$route.params.project_id || '', true),
            }
        },
        computed: {
            ...mapState(usePreviewStore, ['projectInfo']),
        },
        watch: {
            $route: function () {
                this.getDocumentDetail()
            },
        },

        methods: {
            setTitle(title) {
                document.title = title || 'ApiCat 文档预览'
            },

            async getProjectCatalog() {
                try {
                    const { data } = await getProjectCatalog(this.token, this.project_id)
                    this.catalogList = data || []
                } catch (error) {
                    //
                }
            },

            syncTreeTitle(id, title) {
                let el = document.getElementById('tree_' + id)
                if (el && el.innerText !== title) {
                    el.innerText = title
                }
            },

            getDocumentDetail() {
                if (!this.$route.params.node_id || !this.catalogList.length) {
                    return
                }

                this.isLoading = true

                getApiDocumentDetail(this.token, this.project_id, this.$route.params.node_id)
                    .then((res) => {
                        this.documentInfo = res.data || {}
                        this.setTitle(this.documentInfo.title)
                        this.syncTreeTitle(this.$route.params.node_id, this.documentInfo.title)
                    })
                    .catch(() => {
                        //
                    })
                    .finally(() => {
                        this.isLoading = false
                    })
            },
        },

        async mounted() {
            showLoading()
            await this.getProjectCatalog()
            hideLoading()
            this.getDocumentDetail()
        },
    }
</script>
