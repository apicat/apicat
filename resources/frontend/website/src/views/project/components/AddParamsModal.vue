<template>
    <el-dialog
        v-model="isShow"
        :width="400"
        :close-on-click-modal="false"
        append-to-body
        custom-class="show-footer-line vertical-center-modal"
        :title="(isEdit ? '编辑' : '添加') + '参数'"
    >
        <el-form @keyup.enter="handleSubmit('form')" ref="form" :model="form" :rules="rules" label-position="top">
            <el-form-item label="参数名称" prop="name" class="hide_required">
                <el-input v-param @on-change="onChange" v-model="form.name" placeholder="参数名称" maxlength="100" />
            </el-form-item>

            <el-form-item label="参数类型" prop="type">
                <el-select v-model="form.type" placeholder="参数类型" class="w-full">
                    <el-option v-for="item in types" :value="item.value" :key="item.value" :label="item.text" />
                </el-select>
            </el-form-item>

            <el-form-item label="" prop="is_must">
                <el-checkbox v-model="form.is_must">是否必选</el-checkbox>
            </el-form-item>

            <el-form-item label="默认值" prop="default_value" class="hide_required">
                <el-input v-model="form.default_value" placeholder="默认值" maxlength="255" />
            </el-form-item>

            <el-form-item label="参数说明" prop="description" class="hide_required">
                <el-input v-model="form.description" placeholder="参数名称" type="textarea" :autosize="{ minRows: 4, maxRows: 4 }" maxlength="255" />
            </el-form-item>
        </el-form>

        <template #footer>
            <el-button @click="onCloseBtnClick()">取消</el-button>
            <el-button :loading="isLoading" type="primary" @click="handleSubmit('form')">确定</el-button>
        </template>
    </el-dialog>
</template>
<script>
    import { addApiParam, updateApiParam } from '@/api/params'
    import { PARAM_TYPES } from '@natosoft/shared'
    import { isEmpty } from 'lodash-es'
    import { mapState } from 'pinia'
    import { useProjectStore } from '@/stores/project'

    const initForm = {
        name: '',
        type: PARAM_TYPES.TYPES[0].value,
        is_must: true,
        default_value: '',
        description: '',
    }

    export default {
        emits: ['on-ok'],
        directives: {
            param: {
                bind: function (el) {
                    let input = el.querySelector('input')
                    input.removeEventListener('input', el.handler)

                    el.handler = function (e) {
                        e.preventDefault()
                        var pattern = new RegExp("[`~!#%^&*()=|{}':;',<>/?~！#￥……&*（）|{}【】‘；：”“'。，、？]")
                        var rs = ''
                        for (var i = 0; i < this.value.length; i++) {
                            rs = rs + this.value.substr(i, 1).replace(pattern, '')
                        }
                        this.value = rs.replace(/[^\w\.@$-\[\]]/gi, '')
                    }
                    input.addEventListener('input', el.handler)
                },
                unbind: function (el) {
                    let input = el.querySelector('input')
                    input.removeEventListener('input', el.handler)
                },
            },
        },
        data() {
            return {
                isShow: false,
                isEdit: false,
                isLoading: false,
                types: PARAM_TYPES.TYPES,
                form: { ...initForm },
                rules: {
                    name: { required: true, message: '请输入参数名称', trigger: 'blur' },
                },
            }
        },

        computed: {
            ...mapState(useProjectStore, {
                project: 'projectInfo',
            }),
        },

        watch: {
            isShow: function () {
                !this.isShow && this.reset()
            },
        },
        methods: {
            onChange(e) {
                this.form.name = e.target.value
            },
            show(param = {}) {
                !isEmpty(param) && (this.form = param)
                this.isEdit = !isEmpty(param)
                this.isShow = true
            },

            hide() {
                this.isShow = false
                this.reset()
                this.form.name = ''
            },

            onCloseBtnClick() {
                this.isShow = false
                this.reset()
            },

            reset() {
                this.form = { ...initForm }
                this.$refs['form'].resetFields()
            },

            handleSubmit(name) {
                this.$refs[name].validate((valid) => valid && this.submit())
            },

            submit() {
                this.isLoading = true
                let action = this.isEdit ? updateApiParam : addApiParam
                this.form.project_id = this.project.id
                this.form.param_id = this.form.id || ''
                action(this.form)
                    .then((res) => {
                        this.$Message.success(res.msg || (!this.isEdit ? '添加成功!' : '修改成功!'))
                        this.onCloseBtnClick()
                        this.$emit('on-ok')
                    })
                    .catch((e) => {
                        //
                    })
                    .finally(() => {
                        this.isLoading = false
                    })
            },
        },
    }
</script>
