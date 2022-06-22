<template>
    <el-dialog v-model="isShow" width="fit-content" :close-on-click-modal="false" append-to-body :title="title">
        <div class="ac-project-export" v-loading="isLoading" element-loading-text="导入中，请勿关闭窗口！">
            <div class="ac-project-export-item" :key="item.type" v-for="item in exportList">
                <div class="ac-project-export-inner" @click="onItemClick(item)">
                    <img :src="item.icon" :alt="item.text" />
                    <p class="mt-1.5">
                        {{ item.text }}
                    </p>
                </div>
            </div>
        </div>
    </el-dialog>
</template>

<script>
    import { downloadFile } from '@natosoft/shared'
    import { IMPORT_EXPORT_STATE } from '@/common/constant'
    import delay from 'delay'

    export default {
        data() {
            return {
                isShow: false,
                isLoading: false,
                isPoll: false,
                title: '导出项目',
                params: {},
                exportList: [],
                jobParams: {
                    id: null,
                },
                currentConfig: null,
                timer: null,
            }
        },
        watch: {
            isShow: function () {
                if (!this.isShow) {
                    this.clearTimer()
                    this.isLoading = false
                }
            },
        },
        methods: {
            hide() {
                this.isShow = false
            },

            show(exportParams, exportList) {
                this.clearTimer()

                this.isShow = true
                this.exportList = exportList || []
                this.params = exportParams || {}
            },

            onItemClick(item) {
                this.currentConfig = item || null
                this.getImportJobId(this.buildParams(item))
            },

            getImportJobId(params) {
                this.jobParams.id = null

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
                this.isLoading = false
                this.isPoll = false
            },

            async exec() {
                if (!this.currentConfig || !this.currentConfig.getJobResult) {
                    this.clearTimer()
                    return Promise.resolve()
                }

                this.isLoading = true
                try {
                    const { data } = await this.currentConfig.getJobResult(this.params.project_id, this.jobParams.id)

                    let msg = data.description || ''

                    if (data.result === IMPORT_EXPORT_STATE.FAIL && this.isShow) {
                        this.isLoading = false
                        this.$Message.error(msg || '导出失败！')
                        this.clearTimer()
                        return
                    }

                    if (data.result === IMPORT_EXPORT_STATE.WAIT && this.isShow) {
                        return
                    }

                    if (data.result === IMPORT_EXPORT_STATE.FINISH && this.isShow) {
                        this.$Message.success(msg || '导出完成！')
                        this.isLoading = false
                        this.clearTimer()
                        this.hide()
                        data.url && downloadFile(data.url)
                    }
                } catch (e) {
                    this.clearTimer()
                }
            },

            async getExecResult() {
                while (this.isPoll) {
                    await this.exec(this.params)
                    await delay(2000)
                }
            },

            buildParams(item) {
                this.params.type = item.type
                return this.params
            },
        },
        unmounted() {
            this.clearTimer()
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
        }

        &-inner {
            width: 110px;
            height: 110px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            cursor: pointer;

            &:hover {
                background: rgba(70, 146, 255, 0.15);
                border-radius: 4px;
                border: 1px solid var(--primary-color);
            }
        }
    }
</style>
