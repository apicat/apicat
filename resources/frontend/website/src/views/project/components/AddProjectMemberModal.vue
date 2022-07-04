<template>
    <el-dialog
        v-model="isShow"
        :width="340"
        :close-on-click-modal="false"
        title="添加成员"
        append-to-body
        custom-class="show-footer-line vertical-center-modal"
    >
        <el-form @keyup.enter="handleSubmit('teamForm')" ref="teamForm" :model="form" :rules="rules" label-position="top" style="margin-bottom: -19px">
            <el-form-item label="选择成员" prop="user_ids" class="hide_required">
                <el-select v-model="form.user_ids" placeholder="选择成员" no-data-text="暂无成员" filterable multiple clearable class="w-full">
                    <el-option v-for="item in members" :value="item.user_id" :key="item.user_id" :label="item.name" />
                </el-select>
            </el-form-item>

            <el-form-item label="权限" prop="authority">
                <el-select v-model="form.authority" class="w-full">
                    <el-option v-for="item in roles" :value="item.value" :key="item.value" :label="item.text" />
                </el-select>
            </el-form-item>
        </el-form>

        <template #footer>
            <el-button @click="onCloseBtnClick()">取消</el-button>
            <el-button :loading="isLoading" type="primary" @click="handleSubmit('teamForm')">确定</el-button>
        </template>
    </el-dialog>
</template>
<script>
    import { addMember } from '@/api/project'
    import { PROJECT_ROLE_LIST } from '@/common/constant'
    export default {
        emits: ['on-ok'],
        props: {
            members: {
                type: Array,
                default: () => [],
            },
            project: {
                type: Object,
                default: () => ({}),
            },
        },
        data() {
            return {
                isShow: false,
                isLoading: false,
                roles: PROJECT_ROLE_LIST,
                form: {
                    user_ids: [],
                    authority: PROJECT_ROLE_LIST[0].value,
                },
                rules: {
                    user_ids: { required: true, message: '请选择需要添加的成员', trigger: 'change', type: 'array' },
                },
            }
        },

        watch: {
            isShow: function () {
                !this.isShow && this.reset()
            },
        },
        methods: {
            show(project) {
                this.form.project_id = project.id
                this.isShow = true
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
                addMember(this.form)
                    .then((res) => {
                        this.$Message.success(res.msg || '添加成员成功!')
                        this.hide()
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
