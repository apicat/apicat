<template>
    <el-form
        ref="updateHttpCodeAttrs"
        :model="attrs"
        :rules="rules"
        :show-message="false"
        label-width="0"
        size="small"
        class="update-attr-layer update-attr-http-code"
    >
        <el-form-item prop="intro">
            <el-input ref="intro" placeholder="Response Status Code:" v-model="attrs.intro" clearable />
        </el-form-item>

        <el-form-item prop="code">
            <el-select filterable clearable v-model="attrs.code" placeholder="状态码" :teleported="false">
                <el-option v-for="item in httpCodeList" :key="item.code" :label="item.code + ' ' + item.desc" :value="item.code">
                    {{ item.code }} {{ item.desc }}
                </el-option>
            </el-select>
        </el-form-item>
    </el-form>
</template>

<script>
    import { ElForm, ElFormItem, ElInput, ElSelect, ElOption } from 'element-plus'
    import setNode from './mixins/setNode'
    import { $emit } from '@ac/shared'
    import HttpCodeMap from '../../common/HttpCodeMap'

    export default {
        name: 'UpdateHttpCodeAttrs',
        mixins: [setNode],
        components: {
            ElForm,
            ElFormItem,
            ElInput,
            ElSelect,
            ElOption,
        },
        watch: {
            'attrs.code': function () {
                if (!this.attrs.intro) {
                    this.onSubmit()
                }
            },
        },
        data() {
            return {
                httpCodeList: HttpCodeMap.concat([]),
                popper: null,
                rules: {
                    code: { required: true, message: '', trigger: 'change' },
                    intro: { required: true, message: '', trigger: 'change' },
                    codeDesc: { required: true, message: '', trigger: 'change' },
                },
            }
        },

        methods: {
            onHide(cb) {
                this.$refs['updateHttpCodeAttrs'].validate((valid) => {
                    if (valid && this.node) {
                        const attrs = { ...this.attrs }
                        attrs.codeDesc = (this.httpCodeList.find((item) => item.code === this.attrs.code) || { desc: '' }).desc
                        cb(this.isCreate, this.node, { ...this.attrs })
                    }
                })
            },
            onSubmit() {
                this.$refs['updateHttpCodeAttrs'].validate((valid) => {
                    if (valid && this.node) {
                        const attrs = { ...this.attrs }
                        attrs.codeDesc = (this.httpCodeList.find((item) => item.code === this.attrs.code) || { desc: '' }).desc
                        $emit(this, this.isCreate ? 'on-create' : 'on-update-attr', this.node, attrs)
                    }
                })
            },

            focus() {
                this.$nextTick(() => {
                    this.$refs['intro'].focus()
                })
            },

            handleKeyDown(event) {
                if (['Enter', 'Escape'].indexOf(event.key) !== -1) {
                    event.preventDefault()
                    event.stopPropagation()
                }

                if (event.key === 'Enter') {
                    this.onSubmit()
                    return
                }

                if (event.key === 'Escape') {
                    this.close()
                }
            },
        },

        created() {
            window.addEventListener('keydown', this.handleKeyDown)
        },

        unmounted() {
            window.removeEventListener('keydown', this.handleKeyDown)
        },
    }
</script>
