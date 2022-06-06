<template>
    <div class="ac-document is-edit" v-loading="isDocumentLoading" @click="intoEditor">
        <input class="ac-document__title" type="text" maxlength="255" ref="title" v-model="document.title" placeholder="请输入文档标题" />

        <AcEditor v-if="document.content" ref="editor" :document="document.content" :options="editorOptions" @on-change="onDocumentChange" />

        <div class="ac-document__operate">
            <div class="ac-document__operate-inner text-right">
                <el-button :loading="isLoading" type="primary" @click="onSaveBtnClick"> 保存 </el-button>
            </div>
        </div>
    </div>
</template>
<script lang="ts">
    import { defineComponent, defineAsyncComponent, inject } from 'vue'
    import { ElMessage as $Message } from 'element-plus'
    import { updateDoc, getDocumentDetail, getUrlTipList, deleteUrlTip, renameDoc } from '@/api/document'
    import { getApiParamList, addApiParam, deleteApiParam } from '@/api/params'
    import { uploader } from '@/api/uploader'
    import { loadImage } from '@/api/upload'
    import { debounce, isEmpty } from 'lodash-es'
    import { hideLoading } from '@/hooks/useLoading'
    import { useRoute, useRouter } from 'vue-router'

    export default defineComponent({
        components: {
            AcEditor: defineAsyncComponent(() => import('@ac/editor')),
        },

        setup() {
            const updateTreeNode: any = inject('updateTreeNode')
            const $route: any = useRoute()
            const $router: any = useRouter()

            return {
                project_id: $route.params.project_id,
                node_id: parseInt($route.params.node_id as string, 10),
                $route,
                $router,
                updateTreeNode,
            }
        },

        data() {
            return {
                editorOptions: {
                    uploadImage: (file: any) => this.uploadImage(file),
                    getAllCommonParams: () => this.getAllCommonParams(),
                    addCommonParam: (param: any) => this.addCommonParam(param),
                    deleteCommonParam: (param: any) => this.deleteCommonParam(param),
                    getUrlList: () => this.getUrlList(),
                    deleteUrl: (id: any) => this.deleteUrl(id),
                    openNotification: () => this.openNotification(),
                },
                // project_id: this.$route.params.project_id,
                // node_id: parseInt(this.$route.params.node_id as string, 10),
                document: {} as any,
                isLoading: false,
                isDocumentLoading: false,
            }
        },

        watch: {
            '$route.params.node_id': {
                immediate: true,
                handler: function () {
                    this.getDocumentDetail()
                },
            },
            'document.title': function () {
                this.onDocumentTitleChange()
            },
        },

        methods: {
            openNotification() {
                ;(this.$refs['notice'] as any).show()
            },

            intoEditor(e: any) {
                if (e.target.nodeName === 'INPUT') {
                    return
                }

                setTimeout(() => this.$refs.editor && (this.$refs['editor'] as any).editor.focus(), 200)
            },

            uploadImage(file: any) {
                return new Promise((resolve, reject) => {
                    if (!file) {
                        $Message.error('请选择图片')
                        return reject('请选择图片')
                    }

                    if (file.size > 10 * 1024 * 1024) {
                        $Message.error('图片不能超过10MB')
                        return reject('图片不能超过10MB')
                    }

                    uploader()
                        .send(file)
                        .end((error: any, res: any) => {
                            if (!error) {
                                loadImage(res.data).then(({ src }) => resolve(src))
                                return
                            }
                            $Message.error(error || '上传失败，请重试！')
                            reject('上传失败，请重试！')
                        })
                })
            },

            getUrlList() {
                return getUrlTipList(this.project_id)
                    .then((res: any) => res.data)
                    .catch(() => {
                        //
                    })
            },

            deleteUrl(id: any) {
                return deleteUrlTip(this.project_id, id).then(() => {
                    $Message.success('常用URL删除成功')
                })
            },

            addCommonParam(param: any) {
                const newParam = { ...param, project_id: this.project_id }
                delete newParam.sub_params
                delete newParam._id
                return addApiParam(newParam).then((res) => {
                    $Message.success('常用参数添加成功')
                    return res.data
                })
            },

            deleteCommonParam(param: any) {
                return deleteApiParam(this.project_id, param.id).then((res) => {
                    $Message.success('常用参数删除成功')
                    return res.data
                })
            },

            getAllCommonParams() {
                return getApiParamList(this.project_id)
                    .then((res) => res.data)
                    .catch((e) => {
                        //
                    })
            },

            getDocumentDetail() {
                const node_id = parseInt(this.$route.params.node_id as string, 10)

                if (isNaN(node_id)) {
                    hideLoading()
                    return
                }

                this.node_id = node_id
                this.isDocumentLoading = true

                getDocumentDetail(this.project_id, this.node_id)
                    .then((res) => {
                        res.data.project_id = this.project_id
                        res.data.doc_id = res.data.id
                        !res.data.url && (res.data.url = '')
                        res.data.content = JSON.parse(res.data.content || '{}')
                        this.document = res.data
                        this.autoFocus()
                    })
                    .catch((e) => {
                        // this.reactiveNode()
                    })
                    .finally(() => {
                        this.isDocumentLoading = false
                        hideLoading()
                    })
            },

            updateTreeNodeTitle(node: any) {
                node && this.updateTreeNode && this.updateTreeNode(node.doc_id, { title: node.title || '' })
            },

            onSaveBtnClick() {
                this.save()
            },

            save() {
                this.isLoading = true

                updateDoc(this.getDocumentContent())
                    .then((res: any) => {
                        $Message.success(res.message || '保存成功')
                        this.updateTreeNodeTitle(res.data)
                        this.$router.push({ name: 'document.api.detail', params: { project_id: this.project_id, node_id: this.node_id } })
                    })
                    .catch((e) => e)
                    .finally(() => {
                        this.isLoading = false
                    })
            },

            getDocumentContent() {
                if (this.$refs['editor'] && (this.$refs['editor'] as any).editor) {
                    let content = (this.$refs['editor'] as any).editor.getJSON()
                    return { ...this.document, content: JSON.stringify(content) }
                }

                // 默认返回原数据
                return this.document
            },

            onDocumentChange: debounce(function (this: any) {
                updateDoc(this.getDocumentContent()).then((res: any) => {
                    this.updateTreeNodeTitle(res.data)
                })
            }, 500),

            onDocumentTitleChange: debounce(function (this: any) {
                if (isEmpty(this.document.title)) {
                    return
                }
                renameDoc({ project_id: this.project_id, title: this.document.title, doc_id: this.node_id })
                this.updateTreeNodeTitle({ doc_id: this.node_id, title: this.document.title })
            }, 500),

            autoFocus() {
                if (this.$route.query.isNew) {
                    this.$refs.title && (this.$refs.title as any).focus()
                }
            },
        },
    })
</script>
