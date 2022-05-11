<template>
    <el-dialog v-model="isShow" width="fit-content" custom-class="vertical-center-modal" :close-on-click-modal="false" :title="title" append-to-body>
        <div class="ac-project-export" v-loading="isLoading" element-loading-text="导入中，请勿关闭窗口！">
            <div :class="projectItemClass(item)" :key="item.type" v-for="item in importList">
                <div class="ac-project-export-inner">
                    <img :src="item.icon" :alt="item.text" />
                    <p class="mt-1.5">
                        {{ item.text }}
                    </p>

                    <input type="file" @change="onFileChange($event, item)" :accept="item.accept" />
                </div>
            </div>
        </div>
    </el-dialog>
</template>
<script>
    import { IMPORT_EXPORT_STATE } from '@ac/shared'
    import { uploader } from '@/api/uploader'
    import delay from 'delay'

    export default {
        props: {
            maxSize: {
                type: Number,
                default: 1,
            },
            multiple: {
                type: Boolean,
                default: false,
            },
        },
        data() {
            return {
                isShow: false,
                isLoading: false,
                isPoll: false,
                title: '导入文档',
                params: {},
                jobParams: {
                    id: null,
                },
                importList: [],
                currentConfig: null,
                timer: null,
            }
        },
        methods: {
            projectItemClass(item) {
                return [
                    'ac-project-export-item',
                    'ac-project-export-file',
                    {
                        disabled: item.disabled,
                    },
                ]
            },
            hide() {
                this.isShow = false
            },

            show(exportParams, exportList) {
                this.clearTimer()

                this.isShow = true
                this.importList = exportList || []
                this.params = exportParams || {}
            },

            onFileChange(e, config) {
                this.currentConfig = config || {}
                this.validFile(e) && this.upload()
            },

            validFile(e) {
                const files = e.target.files
                if (!files) {
                    return false
                }

                let postFiles = Array.prototype.slice.call(files)
                if (!this.multiple) postFiles = postFiles.slice(0, 1)
                if (postFiles.length === 0) return false

                let file = postFiles[0]

                e.target.value = null

                let { maxSize } = this.currentConfig
                maxSize = maxSize ? maxSize : this.maxSize
                if (file.size > maxSize * 1024 * 1024) {
                    this.$Message.error(`文件大小不能超过${maxSize}MB`)
                    return false
                }

                this.currentConfig.file = file
                return true
            },

            upload() {
                if (!this.currentConfig || !this.currentConfig.file) {
                    return this.$Message.error('未选择文件！')
                }
                this.isLoading = true

                uploader()
                    .options()
                    .send(this.currentConfig.file)
                    .end((error, res) => {
                        if (error) {
                            this.isLoading = false
                            this.$Message.error(error instanceof Error ? error.message : error)
                            return
                        }

                        this.currentConfig.fileUrl = res.data

                        this.getImportJobId(this.buildParams())
                    })
            },

            getImportJobId(params) {
                if (!this.currentConfig || !this.currentConfig.action) {
                    this.isLoading = false
                    return Promise.resolve()
                }

                // 获取job id
                this.currentConfig
                    .action(params)
                    .then((res) => {
                        this.jobParams.id = res.data || null
                        // 启动轮训，获取导入结果
                        this.isPoll = true
                        this.getExecResult()
                    })
                    .catch(() => {
                        this.isLoading = false
                    })
            },

            clearTimer() {
                this.currentConfig = null
                this.isLoading = false
                this.isPoll = false
            },

            async exec() {
                if (!this.currentConfig || !this.currentConfig.getJobResult || !this.jobParams.id) {
                    this.clearTimer()
                    return Promise.resolve()
                }

                try {
                    const { data, msg } = await this.currentConfig.getJobResult(this.params.project_id, this.jobParams.id)

                    if (data.result === IMPORT_EXPORT_STATE.FAIL) {
                        this.$Message.error(msg || '导入失败！')
                        this.clearTimer()
                        return
                    }

                    if (data.result === IMPORT_EXPORT_STATE.WAIT) {
                        return
                    }

                    if (data.result === IMPORT_EXPORT_STATE.FINISH) {
                        this.$Message.success(msg || '导入完成！')
                        this.clearTimer()
                        this.hide()
                        location.reload()
                    }
                } catch (e) {
                    this.clearTimer()
                }
            },

            async getExecResult() {
                while (this.isPoll) {
                    await this.exec()
                    await delay(2000)
                }
            },

            buildParams() {
                if (!this.currentConfig) {
                    return {}
                }

                return {
                    ...this.params,
                    filename: this.currentConfig.file ? this.currentConfig.file.name : '未命名',
                    file: this.currentConfig.fileUrl,
                    type: this.currentConfig.type,
                }
            },
        },
    }
</script>
<style lang="scss" scoped>
    @use '@/assets/stylesheet/mixins/mixins' as *;

    @include b(project-export) {
        display: flex;
        flex-wrap: wrap;
        justify-content: center;
        margin: auto;

        &-item {
            width: 128px;
            height: 128px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;

            img {
                width: 48px;
                height: 48px;
            }

            &.disabled .ac-project-export-inner:hover {
                cursor: not-allowed;
                background: #fafafa;
                border: 1px solid #fafafa;
            }
        }

        &-inner {
            width: 110px;
            height: 110px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            cursor: pointer;
            position: relative;
            overflow: hidden;

            &:hover {
                background: rgba(70, 146, 255, 0.15);
                border-radius: 4px;
                border: 1px solid var(--primary-color);
            }
        }

        &-tip {
            position: absolute;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.4);
            color: rgb(255, 255, 255);
            align-items: center;
            justify-content: center;
            display: none;
        }

        &-file {
            input[type='file'] {
                position: absolute;
                width: 100%;
                height: 100%;
                font-size: 100px;
                opacity: 0;
            }
        }
    }
</style>
