<template>
    <el-form
        ref="updateHttpUrlAttrs"
        :model="attrs"
        :rules="rules"
        label-width="0"
        size="small"
        class="update-attr-layer update-http-url-attrs"
        :show-message="false"
    >
        <el-row type="flex" justify="space-between">
            <el-col :span="12">
                <el-form-item>
                    <el-select class="block" v-model="attrs.method" placeholder="请求方法" :teleported="false">
                        <el-option v-for="item in methods" :key="item.value" :label="item.text" :value="item.value"></el-option>
                    </el-select>
                </el-form-item>
            </el-col>
            <el-col :span="11">
                <el-form-item>
                    <el-select class="block" v-model="attrs.bodyDataType" :teleported="false" placeholder="请求格式">
                        <el-option v-for="item in requestTypes" :key="item.value" :label="item.text" :value="item.value"></el-option>
                    </el-select>
                </el-form-item>
            </el-col>
        </el-row>

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
            <div class="url-list-item" :class="urlItemClass(idx)" v-for="(item, idx) in urls" @click="onUrlItemClick(item)" :key="item.id">
                <div class="url-list-content" :title="item.value">{{ item.value }}</div>
                <el-icon @click.stop="onDeleteUrlBtnClick($event, item.id, idx)"><delete /></el-icon>
            </div>
        </el-form-item>
    </el-form>
</template>

<script>
    import { ElForm, ElFormItem, ElSelect, ElOption, ElCol, ElRow, ElInput, ElIcon } from 'element-plus'
    import { Delete } from '@element-plus/icons-vue'
    import { HTTP_METHODS, REQUEST_BODY_DATA_TYPES } from '../../common/constants'
    import setNode from './mixins/setNode'
    import { $emit } from '@natosoft/shared'

    export default {
        name: 'UpdateHttpUrlAttrs',
        mixins: [setNode],
        components: {
            ElForm,
            ElFormItem,
            ElSelect,
            ElOption,
            ElInput,
            ElCol,
            ElRow,
            ElIcon,
            Delete,
        },

        data() {
            return {
                urls: this.editor.commonUrlManager.urlList,
                selectedIndex: -1,
                methods: HTTP_METHODS.TYPES.concat([]),
                requestTypes: REQUEST_BODY_DATA_TYPES.TYPES.concat([]),

                rules: {
                    url: { message: '', trigger: 'change', type: 'url' },
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
                this.$refs['updateHttpUrlAttrs'].validate((valid) => {
                    if (!(this.attrs.url || '').trim() && !(this.attrs.path || '').trim()) {
                        return
                    }

                    if (valid && this.node && cb) {
                        cb(this.isCreate, this.node, { ...this.attrs })
                    }
                })
            },
            onSubmit() {
                this.$refs['updateHttpUrlAttrs'].validate((valid) => {
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
