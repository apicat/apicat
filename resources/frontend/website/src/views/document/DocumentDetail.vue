<template>
    <div class="ac-document" v-loading="isLoading" element-loading-background="#fff">
        <div v-show="hasDocument && document.id">
            <h1 class="ac-document__title" ref="title">{{ document.title }}</h1>
            <p class="ac-document__desc">
                <el-tooltip effect="dark" :content="document.last_updated_by + ' 最后编辑'" placement="bottom">
                    <span><i class="iconfont iconIconPopoverUser"></i>{{ document.last_updated_by }}</span>
                </el-tooltip>
                <el-tooltip effect="dark" :content="'更新于 ' + document.updated_time" placement="bottom">
                    <span><i class="iconfont icontime"></i>{{ document.updated_time }}</span>
                </el-tooltip>
            </p>
            <div v-if="document.content" class="ProseMirror readonly" v-html="document.content" />
        </div>

        <div v-if="!hasDocument">
            <Result :styles="{ width: '260px', height: 'auto', 'margin-bottom': '26px' }">
                <template #icon>
                    <img src="@/assets/image/icon-empty.png" alt="" />
                </template>
                <template #title>
                    <div style="width: 470px; display: block; margin: auto" v-if="!isGuest">
                        您当前尚未创建文档，请从左侧目录栏点击添加，开始在线维护 API 文档。您还可以将本地项目
                        <a class="text-blue-600" href="javascript:void(0);" @click="onImportBtnCLick">导入</a>
                    </div>
                    <div style="width: 470px; display: block; margin: auto" v-else>您当前尚未创建文档</div>
                </template>
            </Result>
        </div>

        <ac-backtop :bottom="100" :right="100" />

        <div v-html="zoomTemplate" />
    </div>
</template>
<script>
    import { getDocumentDetail, API_DOCUMENT_IMPORT_ACTION_MAPPING } from '@/api/document'
    import { DOCUMENT_TYPES } from '@/common/constant'
    import { useHighlight } from '@/hooks/useHighlight'
    import { inject, ref, watch } from 'vue'
    import { hideLoading } from '@/hooks/useLoading'
    import { useElementBounding } from '@vueuse/core'
    import emitter, { IS_SHOW_DOCUMENT_TITLE } from '@/common/emitter'
    import { storeToRefs } from 'pinia'
    import { useProjectStore } from '@/stores/project'
    import useIdPublicParam, { generateProjectOrIterateParams } from '@/hooks/useIdPublicParam'
    import { useDocumentDetailInteractive } from '@/hooks/useDocumentDetailInteractive'

    export default {
        watch: {
            '$route.params.node_id': function () {
                this.getDocumentDetail()
            },
        },

        setup() {
            const title = ref(null)
            const { top } = useElementBounding(title)
            const projectStore = useProjectStore()
            const { isGuest, projectInfo } = storeToRefs(projectStore)

            watch(top, () => emitter.emit(IS_SHOW_DOCUMENT_TITLE, top.value < 25 ? true : false), {
                immediate: true,
            })

            const { initHighlight } = useHighlight()
            const documentImportModal = inject('documentImportModal')
            const setDocumentTitle = inject('setDocumentTitle')
            const publicParams = useIdPublicParam()

            return {
                pid: projectInfo.value.id,
                projectInfo,
                isGuest,
                title,
                initHighlight,
                documentImportModal,
                setDocumentTitle,
                publicParams,
            }
        },

        data() {
            return {
                hasDocument: true,
                isLoading: true,
                DOCUMENT_TYPES,
                document: {},
                project_id: null,
            }
        },

        methods: {
            // 初始化静态文档交互
            initStaticDocInteractive() {
                useDocumentDetailInteractive('.document-detail')
            },

            onImportBtnCLick() {
                this.documentImportModal.show(
                    {
                        ...generateProjectOrIterateParams(this.publicParams),
                    },
                    API_DOCUMENT_IMPORT_ACTION_MAPPING
                )
            },

            getDocumentDetail() {
                const doc_id = parseInt(this.$route.params.node_id, 10)

                if (isNaN(doc_id)) {
                    hideLoading()
                    this.isLoading = false
                    this.hasDocument = false
                    return
                }

                this.doc_id = doc_id
                this.isLoading = true
                this.hasDocument = true

                const data = { project_id: this.projectInfo.id, doc_id: this.doc_id }

                getDocumentDetail(data, 'html')
                    .then((res) => {
                        this.document = this.transferDoc(res.data || {})
                        this.setDocumentTitle(this.document.title)
                        this.initStaticDocInteractive()
                    })
                    .catch((e) => {
                        //
                    })
                    .finally(() => {
                        this.$nextTick(() => {
                            this.isLoading = false
                            hideLoading()
                        })
                    })
            },

            transferDoc(doc) {
                return doc
            },
        },

        mounted() {
            emitter.emit(IS_SHOW_DOCUMENT_TITLE, false)
            this.getDocumentDetail()
        },

        unmounted() {
            this.document = {}
            this.isLoading = true
            this.hasDocument = true
        },
    }
</script>
