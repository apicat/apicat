<template>
    <el-dialog
        v-model="isShow"
        :width="400"
        :close-on-click-modal="false"
        append-to-body
        class="show-footer-line vertical-center-modal"
        title="修改成员密码"
    >
        <el-form @submit.prevent="handleSubmit" ref="formRef" :model="form" :rules="rules" label-position="top">
            <el-form-item label="新密码" prop="password" class="hide_required">
                <el-input type="password" v-model="form.password" placeholder="密码" maxlength="100" show-password />
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
    import { modifyMemberPassword } from '@/api/team'
    import { ElMessage as $Message } from 'element-plus'

    const initForm = {
        password: '',
    } as any

    export default defineComponent({
        setup() {
            const formRef = ref()
            let user_id: number | null = null

            const state = reactive({
                isEdit: false,
                isShow: false,
                isLoading: false,
                form: { ...initForm },
            })

            const rules = {
                password: { required: true, message: '请输入密码' },
            } as any

            const onCloseBtnClick = () => {
                state.isShow = false
            }

            const reset = () => {
                state.form = { ...initForm }
                user_id = null
                formRef.value?.resetFields()
            }

            const handleSubmit = () => {
                formRef.value?.validate((valid: boolean) => {
                    if (valid && user_id) {
                        state.isLoading = true
                        modifyMemberPassword({ ...state.form, user_id })
                            .then((res: any) => {
                                $Message.success(res.msg || '密码修改成功')
                                onCloseBtnClick()
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

            const show = (userId: number) => {
                state.isShow = true
                user_id = userId
                formRef.value?.clearValidate()
            }

            return {
                rules,
                ...toRefs(state),
                formRef,

                onCloseBtnClick,
                handleSubmit,
                show,
            }
        },
    })
</script>
