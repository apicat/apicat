<template>
    <el-form ref="updateApiUrlAttrs" :model="attrs" :rules="rules" size="small" :show-message="false" class="update-attr-layer update-api-url-attrs">
        <el-row type="flex" justify="space-between">
            <el-col :span="12">
                <el-form-item prop="url">
                    <el-input ref="url" placeholder="protocol://hostname[:port]" v-model="attrs.url" clearable />
                </el-form-item>
            </el-col>

            <el-col :span="11">
                <el-form-item>
                    <el-input placeholder="path" v-model="attrs.path" clearable />
                </el-form-item>
            </el-col>
        </el-row>

        <el-form-item class="url-list scroll-content" v-show="urls.length">
            <div class="url-list-item" :class="urlItemClass(idx)" v-for="(item, idx) in urls" :key="item.id" @click="onUrlItemClick(item)">
                <div class="url-list-content" :title="item.value">{{ item.value }}</div>
                <el-icon @click.stop="onDeleteUrlBtnClick($event, item.id, idx)"><delete /></el-icon>
            </div>
        </el-form-item>
    </el-form>
</template>

<script>
    import { ElCol, ElForm, ElFormItem, ElInput, ElRow, ElIcon } from 'element-plus'
    import { Delete } from '@element-plus/icons-vue'
    import setNode from './mixins/setNode'
    import { $emit } from '@natosoft/shared'

    const protocolReg = /(\w+):\/\/([^/:]+)(:\d*)?/

    export default {
        name: 'UpdateApiUrlAttrs',
        mixins: [setNode],

        components: {
            ElCol,
            ElForm,
            ElFormItem,
            ElInput,
            ElRow,
            ElIcon,
            Delete,
        },

        data() {
            return {
                urls: this.editor.commonUrlManager.urlList,
                selectedIndex: -1,
                rules: {
                    url: {
                        trigger: 'change',
                        validator: (rule, value, callback) => {
                            let v = (value || '').trim()
                            if (!v) {
                                return callback()
                            }

                            if (!protocolReg.test(value)) {
                                return callback(new Error('协议格式错误'))
                            }

                            return callback()
                        },
                    },
                },
            }
        },
        watch: {
            'attrs.url': function () {
                this.editor.commonUrlManager.filterUrls(this.attrs.url)
            },
        },

        methods: {
            onHide(cb) {
                this.$refs['updateApiUrlAttrs'].validate((valid) => {
                    if (!(this.attrs.url || '').trim() && !(this.attrs.path || '').trim()) {
                        return
                    }

                    if (valid && this.node && cb) {
                        cb(this.isCreate, this.node, { ...this.attrs })
                    }
                })
            },
            onSubmit() {
                this.$refs['updateApiUrlAttrs'].validate((valid) => {
                    if (!(this.attrs.url || '').trim() && !(this.attrs.path || '').trim()) {
                        this.close()
                        return
                    }

                    if (valid && this.node) {
                        $emit(this, this.isCreate ? 'on-create' : 'on-update-attr', this.node, { ...this.attrs })
                    }
                })
            },

            urlItemClass(index) {
                return {
                    selected: this.selectedIndex === index,
                }
            },

            onUrlItemClick(params) {
                this.attrs.url = params.url
            },

            onDeleteUrlBtnClick(e, id) {
                this.editor.commonUrlManager && this.editor.commonUrlManager.deleteUrlById(id)
            },

            handleKeyDown(event) {
                if (['Enter', 'ArrowUp', 'ArrowDown', 'Escape'].indexOf(event.key) !== -1) {
                    event.preventDefault()
                    event.stopPropagation()
                }

                if (event.key === 'Enter') {
                    //  用户进行了选择
                    if (this.selectedIndex !== -1) {
                        const item = this.urls[this.selectedIndex]
                        item && this.onUrlItemClick(item)
                        this.selectedIndex = -1
                        return
                    }

                    // 默认回车事件,提交
                    this.onSubmit()
                    return
                }

                if (event.key === 'ArrowUp' && this.urls.length) {
                    const prevIndex = this.selectedIndex - 1
                    const total = this.urls.length - 1

                    if (this.selectedIndex === 0) {
                        this.selectedIndex = total
                    } else {
                        this.selectedIndex = Math.max(0, prevIndex)
                    }
                    return
                }

                if (event.key === 'ArrowDown' && this.urls.length) {
                    const total = this.urls.length - 1
                    const nextIndex = this.selectedIndex + 1

                    if (this.selectedIndex === total) {
                        this.selectedIndex = 0
                    } else {
                        this.selectedIndex = Math.min(nextIndex, total)
                    }
                }

                if (event.key === 'Escape') {
                    this.close()
                }
            },

            focus() {
                this.$nextTick(() => {
                    this.$refs['url'].focus()
                })
            },
        },

        mounted() {
            window.addEventListener('keydown', this.handleKeyDown)
        },

        unmounted() {
            window.removeEventListener('keydown', this.handleKeyDown)
        },
    }
</script>
