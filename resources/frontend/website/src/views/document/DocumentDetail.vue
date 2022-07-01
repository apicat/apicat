<template>
    <div class="ac-document" v-loading="isLoading">
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
    import mediumZoom from 'medium-zoom'
    import tippy from 'tippy.js'
    import { getDocumentDetail, API_DOCUMENT_IMPORT_ACTION_MAPPING } from '@/api/document'
    import { toggleClass, getAttr, hasClass, showOrHide } from '@natosoft/shared'
    import { DOCUMENT_TYPES } from '@/common/constant'
    import { useHighlight } from '@/hooks/useHighlight'
    import { inject, ref, watch } from 'vue'
    import { hideLoading } from '@/hooks/useLoading'
    import { useElementBounding } from '@vueuse/core'
    import emitter, { IS_SHOW_DOCUMENT_TITLE } from '@/common/emitter'
    import { storeToRefs } from 'pinia'
    import { useProjectStore } from '@/stores/project'
    import { DOCUMENT_EDIT_NAME } from '@/router/constant'

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
            const title = ref(null)
            const { top } = useElementBounding(title)
            const projectStore = useProjectStore()
            const { isGuest } = storeToRefs(projectStore)

            watch(top, () => emitter.emit(IS_SHOW_DOCUMENT_TITLE, top.value < 25 ? true : false), {
                immediate: true,
            })

            const { initHighlight } = useHighlight()
            const documentImportModal = inject('documentImportModal')
            const setDocumentTitle = inject('setDocumentTitle')

            return {
                isGuest,
                title,
                initHighlight,
                documentImportModal,
                setDocumentTitle,
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
                this.$router.push({ name: DOCUMENT_EDIT_NAME, params: { project_id: this.$route.params.project_id, node_id: this.$route.params.node_id } })
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
                    hideLoading()
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
