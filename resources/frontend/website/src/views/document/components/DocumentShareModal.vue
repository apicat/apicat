<template>
    <el-dialog v-model="isShow" title="分享文档" class="vertical-center-modal" :width="540" append-to-body>
        <div v-loading="isShowDetail" class="-m-4">
            <div v-if="document && document.visibility === 1">
                <el-form label-position="top" class="px-6 py-3">
                    <el-form-item label="">
                        <el-input readonly v-model="shareUrl">
                            <template #append>
                                <el-button type="primary" @click="copy">{{ copyText }}</el-button>
                            </template>
                        </el-input>
                    </el-form-item>
                </el-form>
            </div>

            <div v-if="document && document.visibility === 0">
                <div class="flex items-center px-6" :class="{ 'py-3': !isShare, 'pt-3': isShare }">
                    <div class="flex-1">
                        <h4>开启分享</h4>
                        <div class="ivu-list-item-meta-description">开启分享后，获得链接的人可以访问项目内容。</div>
                    </div>
                    <el-switch :loading="isLoading" v-model="isShare" @change="onShareStatusSwitch" inline-prompt active-text="开" inactive-text="关" />
                </div>

                <el-divider v-if="isShare" />

                <el-form v-if="isShare" label-position="top" class="px-6">
                    <el-form-item label="文档链接">
                        <el-input readonly v-model="shareUrl">
                            <template #append>
                                <el-button type="primary" @click="copy">{{ copyText }}</el-button>
                            </template>
                        </el-input>
                    </el-form-item>

                    <el-form-item label="密码">
                        <el-form-item prop="date" class="mr-1">
                            <el-input readonly v-model="password" />
                        </el-form-item>
                        <el-button :loading="isLoadingResetPwd" @click="onResetPasswordBtnClick">重置密码</el-button>
                    </el-form-item>
                </el-form>
            </div>
        </div>
        <textarea ref="copyTextEl" style="position: fixed; left: -9999px" />
    </el-dialog>
</template>

<script>
    import { shareDoc, shareDetailDoc, resetDocShareSecretkey, generateDocumentDetailPath } from '@/api/document'
    import { DOCUMENT_VISIBLE_TYPES } from '@/common/constant'
    import { toRefs, reactive, ref } from 'vue'
    import { useProjectStore } from '@/stores/project'
    import { storeToRefs } from 'pinia'
    import { generatePreviewDocumentPath } from '@/api/preview'

    export default {
        watch: {
            shareData() {
                if (this.shareData.docId !== -1) {
                    this.getShareDetail()
                }
            },
        },
        setup() {
            const projectStore = useProjectStore()
            const { projectInfo } = storeToRefs(projectStore)

            const state = reactive({
                projectInfo,
                document: null,
                isShow: false,
                isShowDetail: false,
                isLoading: false,
                isLoadingResetPwd: false,
                shareUrl: '',
                isShare: false,
                password: '',
                copyText: '',
                copyTextEl: ref(),

                shareData: {
                    docId: -1,
                },
            })

            return {
                ...toRefs(state),
            }
        },

        methods: {
            onShareStatusSwitch(status) {
                if (this.document !== null) {
                    this.isLoading = true
                    shareDoc({ project_id: this.projectInfo.id, doc_id: this.shareData.docId, share: status })
                        .then(({ data }) => {
                            if (status) {
                                this.updateLinkAndPassword(generatePreviewDocumentPath(this.shareData.docId), data.secret_key)
                                this.document.secret_key = data.secret_key
                            } else {
                                this.document.secret_key = ''
                                this.isShare = false
                            }
                        })
                        .catch((e) => {
                            console.log(e)
                            this.isShare = !this.isShare
                        })
                        .finally(() => {
                            this.isLoading = false
                        })
                }
            },

            copy() {
                this.copyTextEl.select()
                if (!document.execCommand('copy')) return
                this.copyText = '复制成功'
                setTimeout(() => {
                    this.copyText = this.document.visibility === DOCUMENT_VISIBLE_TYPES.PUBLIC ? '复制链接' : '复制链接和密码'
                }, 2000)
            },

            show(shareData) {
                this.shareData = shareData || { docId: -1 }
                this.isShow = true
            },

            setDocument(document = {}) {
                this.document = document
                const isPublic = document.visibility === DOCUMENT_VISIBLE_TYPES.PUBLIC
                this.copyText = isPublic ? '复制链接' : '复制链接和密码'
                this.isShare = !isPublic && !!document.secret_key

                document.link = isPublic
                    ? generateDocumentDetailPath(this.projectInfo.id, this.shareData.docId)
                    : this.isShare
                    ? generatePreviewDocumentPath(this.shareData.docId)
                    : ''

                this.updateLinkAndPassword(document.link, document.secret_key)

                return this
            },

            updateLinkAndPassword(link, secret_key) {
                this.shareUrl = link || this.shareUrl
                this.password = secret_key || this.password
                this.updateCopyText()
            },

            onResetPasswordBtnClick() {
                this.isLoadingResetPwd = true
                resetDocShareSecretkey({ project_id: this.projectInfo.id, doc_id: this.shareData.docId })
                    .then(({ data }) => this.updateLinkAndPassword(null, data))
                    .finally(() => {
                        this.isLoadingResetPwd = false
                    })
            },

            updateCopyText() {
                this.copyTextEl.value =
                    this.document.visibility === DOCUMENT_VISIBLE_TYPES.PUBLIC ? this.shareUrl : [`链接：${this.shareUrl}`, `密码：${this.password}`].join('\n')
            },

            getShareDetail() {
                this.isShowDetail = true
                shareDetailDoc({ project_id: this.projectInfo.id, doc_id: this.shareData.docId })
                    .then(({ data }) => {
                        this.setDocument(data)
                    })
                    .finally(() => {
                        this.isShowDetail = false
                    })
            },
        },
    }
</script>
