<template>
    <div class="ac-document is-detail" v-loading="isLoading">
        <div v-show="hasDocument && document.id">
            <h1 class="ac-document__title">
                {{ document.title }}
            </h1>

            <div v-if="document.content" class="ProseMirror readonly" v-html="document.content" />

            <div class="ac-document__operate" v-show="!isLoading">
                <div class="ac-document__operate-inner text-right">
                    <el-button type="primary" @click="onEditBtnClick"> 编辑 </el-button>
                </div>
            </div>
        </div>

        <div v-if="!hasDocument">
            <Result :styles="{ width: '260px', height: 'auto', 'margin-bottom': '26px' }">
                <template #icon>
                    <img src="@/assets/image/icon-empty.png" alt="" />
                </template>
                <template #title>
                    <div style="width: 470px; display: block; margin: auto">
                        您当前尚未创建文档，请从左侧目录栏点击添加，开始在线维护 API 文档。您还可以将本地项目
                        <a class="text-blue-600" href="javascript:void(0);" @click="onImportBtnCLick">导入</a>
                    </div>
                </template>
            </Result>
        </div>

        <ac-backtop :bottom="100" :right="100" />

        <div v-html="zoomTemplate" />
    </div>
</template>
<script>
    import { getDocumentDetail, API_DOCUMENT_IMPORT_ACTION_MAPPING } from '@/api/document'
    import mediumZoom from 'medium-zoom'
    import tippy from 'tippy.js'
    import { toggleClass, getAttr, hasClass, showOrHide, DOCUMENT_TYPES } from '@ac/shared'
    import { useHighlight } from '@/hooks/useHighlight'
    import { inject } from 'vue'
    import { hideLoading } from '@/hooks/useLoading'

    function expand(pid, isExpand) {
        document.querySelectorAll('[data-pid="' + pid + '"]').forEach(function (el) {
            let arrow = el.querySelector('i.editor-arrow-right')
            if (arrow && !hasClass(arrow, 'expand')) {
                toggleClass(arrow, 'expand')
            }
            let id = getAttr(el, 'data-id')
            el.style.display = isExpand ? null : 'none'
            id && expand(id, isExpand)
        })
    }

    export default {
        watch: {
            '$route.params.node_id': function () {
                this.getDocumentDetail()
            },
        },

        setup() {
            const { initHighlight } = useHighlight()
            const documentImportModal = inject('documentImportModal')

            return {
                initHighlight,
                documentImportModal,
            }
        },

        data() {
            return {
                zoomTemplate: `<template id="template-zoom-image">
                            <div class="zoom-image-wrapper">
                                <div class="zoom-image-container" data-zoom-container></div>
                            </div>
                          </template>`,
                zoomImageOption: {
                    template: '#template-zoom-image',
                    container: '[data-zoom-container]',
                },
                hasDocument: true,
                isLoading: true,
                DOCUMENT_TYPES,
                document: {},
                project_id: null,
            }
        },

        methods: {
            onEditBtnClick() {
                this.$router.push({ name: 'document.api.edit', params: { project_id: this.$route.params.project_id, node_id: this.$route.params.node_id } })
            },

            initTableToggle() {
                document.querySelectorAll('.ac-param-table .editor-arrow-right').forEach(function (el) {
                    el.onclick = function () {
                        expand(getAttr(this, 'data-id'), !hasClass(this, 'expand'))
                        toggleClass(this, 'expand')
                    }
                })

                document.querySelectorAll('div.collapse-title .response_body_title').forEach(function (el) {
                    el.onclick = function () {
                        let h3 = this.parentElement
                        let parent = h3.parentElement
                        let isShow = hasClass(parent, 'close')
                        showOrHide(h3.nextElementSibling, isShow)
                        showOrHide(parent.nextElementSibling, isShow)
                        toggleClass(parent, 'close')
                    }
                })

                document.querySelectorAll('h3.collapse-title >span').forEach(function (el) {
                    el.onclick = function () {
                        let parent = this.parentElement
                        let isShow = hasClass(parent, 'close')
                        showOrHide(parent.nextElementSibling, isShow)
                        toggleClass(parent, 'close')
                    }
                })
            },

            initMediumZoom() {
                mediumZoom('.ProseMirror .image-view img', this.zoomImageOption)
            },

            initTippy() {
                tippy('[data-tippy-content]', { theme: 'light', appendTo: document.querySelector('.ProseMirror') })
            },

            initCodeBlockToClipboard() {
                document.querySelectorAll('.code-block button').forEach((el) => {
                    el.setAttribute('data-text', el.parentElement.querySelector('code').innerText)
                })
            },

            // 初始化静态文档交互
            initStaticDocInteractive() {
                this.$nextTick(() => {
                    this.initTableToggle()
                    this.initTippy()
                    this.initMediumZoom()
                    this.initCodeBlockToClipboard()
                    this.initHighlight(document.querySelectorAll('pre code'))
                })
            },

            onImportBtnCLick() {
                this.documentImportModal.show(
                    {
                        project_id: this.$route.params.project_id,
                    },
                    API_DOCUMENT_IMPORT_ACTION_MAPPING
                )
            },

            getDocumentDetail() {
                const doc_id = parseInt(this.$route.params.node_id, 10)

                if (isNaN(doc_id)) {
                    this.isLoading = false
                    this.hasDocument = false
                    return
                }

                this.doc_id = doc_id
                this.project_id = parseInt(this.$route.params.project_id, 10)

                this.isLoading = true
                this.hasDocument = true
                getDocumentDetail(this.project_id, this.doc_id, 'html')
                    .then((res) => {
                        this.document = this.transferDoc(res.data || {})
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
                // doc.content = JSON.parse(doc.content || "{}")
                return doc
            },
        },

        mounted() {
            this.getDocumentDetail()
        },

        unmounted() {
            this.document = {}
            this.isLoading = true
            this.hasDocument = true
        },
    }
</script>
