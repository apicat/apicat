<template>
    <el-dialog v-model="isShow" :width="540" :close-on-click-modal="false" title="分享项目" append-to-body>
        <div v-if="project && project.visibility === 1" class="-m-4">
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

        <div v-if="project && project.visibility === 0" class="-m-4">
            <div class="px-6 flex items-center" :class="{ 'py-3': !isShare, 'pt-3': isShare }">
                <div class="flex-1">
                    <h4>开启分享</h4>
                    <div class="ivu-list-item-meta-description">开启分享后，获得链接的人可以访问项目内容。</div>
                </div>
                <el-switch :loading="isLoading" v-model="isShare" @change="onShareStatusSwitch" inline-prompt active-text="开" inactive-text="关" />
            </div>

            <el-divider v-if="isShare" />

            <el-form v-if="isShare" label-position="top" class="px-6">
                <el-form-item label="项目链接">
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

        <textarea ref="copyTextEl" style="position: fixed; left: -9999px"></textarea>
    </el-dialog>
</template>

<script>
    import { share, resetSecretkey } from '@/api/project'
    import { toRefs, reactive, ref } from 'vue'

    export default {
        setup() {
            const state = reactive({
                project: null,
                dom: null,
                isShow: false,
                isLoading: false,
                isLoadingResetPwd: false,
                shareUrl: '',
                isShare: false,
                password: '',
                copyText: '',
                copyTextEl: ref(),
            })

            return {
                ...toRefs(state),
            }
        },
        methods: {
            onShareStatusSwitch(status) {
                if (this.project !== null) {
                    this.isLoading = true
                    share(this.project.id, Number(status))
                        .then(({ msg, data }) => {
                            data = data || {}
                            this.$Message.success(msg || '操作成功')
                            this.updateLinkAndPassword(data.link, data.secret_key)
                            if (status) {
                                this.project.secret_key = data.secret_key
                            } else {
                                delete this.project.secret_key
                            }
                        })
                        .catch((e) => {
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
                    this.copyText = this.project.visibility === 1 ? '复制链接' : '复制链接和密码'
                }, 2000)
            },

            show(project = {}) {
                this.isShow = true
                this.project = project

                this.copyText = project.visibility === 1 ? '复制链接' : '复制链接和密码'

                if (project.visibility === 0 && project.secret_key) {
                    this.isShare = true
                } else {
                    this.isShare = false
                }

                this.updateLinkAndPassword(project.preview_link, project.secret_key)
            },

            updateLinkAndPassword(link, secret_key) {
                this.shareUrl = link || this.shareUrl
                this.password = secret_key || this.password
                this.updateCopyText()
            },

            onResetPasswordBtnClick() {
                this.isLoadingResetPwd = true
                resetSecretkey(this.project.id)
                    .then(({ data }) => {
                        this.project.secret_key = data
                        this.updateLinkAndPassword(null, data)
                    })
                    .finally(() => {
                        this.isLoadingResetPwd = false
                    })
            },

            updateCopyText() {
                var copyStr = [this.project.name ? `《${this.project.name}》` : '', `链接：${this.shareUrl}`]
                this.project.visibility !== 1 && copyStr.push(`密码：${this.password}`)

                this.$nextTick(() => {
                    this.copyTextEl.value = copyStr.join('\n')
                })
            },
        },
    }
</script>
