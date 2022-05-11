<template>
    <el-card shadow="never">
        <template #header>
            <span>修改密码</span>
        </template>

        <el-form
            @submit.prevent="onSubmitBtnClick"
            @keyup.enter="onSubmitBtnClick"
            ref="formRef"
            :model="form"
            :rules="rules"
            label-position="top"
            class="w-96"
        >
            <el-form-item label="旧密码" prop="password">
                <el-input type="password" v-model="form.password" placeholder="旧密码" autocomplete="on" />
            </el-form-item>

            <el-form-item label="新密码" prop="new_password">
                <el-input type="password" v-model="form.new_password" placeholder="新密码" autocomplete="on" />
            </el-form-item>

            <el-form-item label="确认新密码" prop="new_password_confirmation">
                <el-input type="password" v-model="form.new_password_confirmation" placeholder="新密码" autocomplete="on" />
            </el-form-item>

            <el-form-item>
                <el-button type="primary" @click="onSubmitBtnClick" :loading="isLoading">保存</el-button>
            </el-form-item>
        </el-form>
    </el-card>
</template>
<script setup lang="ts">
    import { modifyPassword } from '@/api/user'
    import { reactive, ref } from 'vue'

    const minLengthRule = { type: 'string', min: 8, message: '密码至少8位', trigger: 'blur' }

    const formRef = ref()

    const form = reactive({
        password: '',
        new_password: '',
        new_password_confirmation: '',
    })

    const rules = {
        password: [{ required: true, message: '请输入旧密码', trigger: 'blur' }, minLengthRule],
        new_password: [{ required: true, message: '请输入新密码', trigger: 'blur' }, minLengthRule],
        new_password_confirmation: [
            {
                required: true,
                message: '请输入确认新密码',
                trigger: 'blur',
            },
            minLengthRule,
            {
                validator: (rule: any, value: string, callback: any) => {
                    if (value === '') {
                        callback(new Error('请输入确认新密码'))
                    } else if (value !== form.new_password) {
                        callback(new Error('新密码不一致'))
                    } else {
                        callback()
                    }
                },
                trigger: 'blur',
            },
        ],
    } as any

    const [isLoading, execute] = modifyPassword()

    const onSubmitBtnClick = () => {
        formRef?.value.validate(async (valid: boolean) => {
            if (valid) {
                await execute(form)
            }
        })
    }
</script>
