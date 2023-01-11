<template>
    <el-dialog
        v-model="isShow"
        :width="400"
        :close-on-click-modal="false"
        append-to-body
        class="show-footer-line vertical-center-modal"
        :title="(isEdit ? '编辑' : '添加') + '成员'"
    >
        <el-form @keyup.enter="handleSubmit" ref="formRef" :model="form" :rules="rules" label-position="top">
            <el-form-item label="姓名" prop="name" class="hide_required">
                <el-input v-model="form.name" placeholder="姓名" maxlength="100" />
            </el-form-item>

            <el-form-item label="邮箱" prop="email" class="hide_required">
                <el-input v-model="form.email" placeholder="邮箱" maxlength="100" />
            </el-form-item>

            <el-form-item label="密码" prop="password" class="hide_required" v-if="!isEdit">
                <el-input type="password" v-model="form.password" placeholder="密码" maxlength="100" show-password />
            </el-form-item>

            <el-form-item label="权限" prop="authority" v-if="isAdmin">
                <el-select v-model="form.authority" placeholder="选择权限" class="w-full">
                    <el-option v-for="item in TEAM_ROLE_LIST" :value="item.value" :key="item.value" :label="item.text" />
                </el-select>
            </el-form-item>
        </el-form>

        <template #footer>
            <el-button @click="onCloseBtnClick">取消</el-button>
            <el-button :loading="isLoading" type="primary" @click="handleSubmit">确定</el-button>
        </template>
    </el-dialog>
</template>

<script lang="ts">
    import { defineComponent, reactive, toRefs, ref, watch } from 'vue'
    import { storeToRefs } from 'pinia'
    import { modifyMember, addMember } from '@/api/team'
    import { TEAM_ROLE, TEAM_ROLE_LIST } from '@/common/constant'
    import { ElMessage as $Message } from 'element-plus'
    import { useUserStore } from '@/stores/user'

    const initForm = {
        name: '',
        email: '',
        password: '',
        authority: TEAM_ROLE.NORMAL,
    } as any

    export default defineComponent({
        emits: ['on-ok'],
        setup(props, { emit }) {
            const formRef = ref()
            const userStore = useUserStore()
            const { isAdmin } = storeToRefs(userStore)

            const state = reactive({
                isEdit: false,
                isShow: false,
                isLoading: false,
                form: { ...initForm },
            })

            const rules = {
                name: { required: true, message: '请输入姓名' },
                password: { required: true, message: '请输入密码' },
                email: { required: true, type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
            } as any

            const onCloseBtnClick = () => {
                state.isShow = false
            }

            const reset = () => {
                state.form = { ...initForm }
                formRef.value?.resetFields()
            }

            const handleSubmit = () => {
                formRef.value?.validate((valid: boolean) => {
                    if (valid) {
                        state.isLoading = true
                        let action = state.isEdit ? modifyMember : addMember
                        const data = { ...state.form }

                        if (state.isEdit) {
                            delete data.password
                        }

                        action(data)
                            .then((res: any) => {
                                $Message.success(res.msg || (!state.isEdit ? '添加成功!' : '修改成功!'))
                                onCloseBtnClick()
                                emit('on-ok')
                            })
                            .catch((e) => {
                                //
                            })
                            .finally(() => {
                                state.isLoading = false
                            })
                    }
                })
            }

            watch(
                () => state.isShow,
                () => !state.isShow && reset()
            )

            const show = (member: any) => {
                state.isEdit = !(member === undefined)
                if (member) {
                    state.form = { ...member }
                }
                state.isShow = true
                formRef.value?.clearValidate()
            }

            return {
                isAdmin,
                rules,
                ...toRefs(state),

                TEAM_ROLE_LIST,
                formRef,
                onCloseBtnClick,
                handleSubmit,
                show,
            }
        },
    })
</script>
