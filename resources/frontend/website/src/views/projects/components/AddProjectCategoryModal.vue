<template>
    <el-dialog v-model="isShow" :width="340" :close-on-click-modal="false" :title="title" append-to-body custom-class="show-footer-line">
        <el-form ref="teamForm" :model="form" :rules="rules" label-position="top" style="margin-bottom: -19px" @submit.prevent="handleSubmit('teamForm')">
            <el-form-item label="分组名称" prop="name" class="hide_required">
                <el-input type="text" ref="input" v-model="form.name" placeholder="请输入分组名称" :maxlength="255" />
            </el-form-item>
        </el-form>

        <template #footer>
            <el-button @click="onCloseBtnClick()"> 取消 </el-button>
            <el-button :loading="isLoading" type="primary" @click="handleSubmit('teamForm')"> 确定 </el-button>
        </template>
    </el-dialog>
</template>
<script>
    import { mapActions } from 'pinia'
    import { useProjectsStore } from '@/stores/projects'

    export default {
        data() {
            return {
                title: '',
                isShow: false,
                isLoading: false,
                isEdit: false,
                form: {
                    name: '',
                },
                categoryInfo: null,
                rules: {
                    name: { required: true, message: '请输入分组名称', trigger: 'blur' },
                },
            }
        },

        watch: {
            isShow: function () {
                if (this.isShow) {
                    this.title = `${this.isEdit ? '编辑' : '添加'}分组`
                }
                !this.isShow && this.reset()
            },
        },
        methods: {
            ...mapActions(useProjectsStore, ['addProjectGroup', 'renameProjectGroup']),

            show(info) {
                this.isEdit = !!info
                if (info) {
                    this.categoryInfo = info
                    this.form = { ...info }
                } else {
                    this.form.name = ''
                }

                this.isShow = true
                this.focus()
            },
            hide() {
                this.isShow = false
            },

            onCloseBtnClick() {
                this.isShow = false
                this.reset()
            },

            reset() {
                this.$refs['teamForm'].resetFields()
            },

            handleSubmit(name) {
                this.$refs[name].validate((valid) => valid && this.submit())
            },

            submit() {
                this.isLoading = true
                let action = this.isEdit ? this.renameProjectGroup : this.addProjectGroup

                action(this.form)
                    .then((res) => {
                        this.$Message.success(res.msg || `${this.isEdit ? '编辑' : '添加'}分组成功!`)
                        this.onCloseBtnClick()
                    })
                    .finally(() => {
                        this.isLoading = false
                    })
                    .catch((e) => {
                        //
                    })
            },

            focus() {
                setTimeout(() => this.$refs.input.focus(), 200)
            },
        },
    }
</script>
