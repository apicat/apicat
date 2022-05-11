<template>
    <div class="flex min-h-screen">
        <main class="ac-login p-2">
            <div class="ac-login__box shadow-xl">
                <div class="-ml-1 ac-login__logo">
                    <router-link class="inline-flex items-center" to="/">
                        <img src="@/assets/image/logo.svg" alt="ApiCat" /><span class="logo-text logo-apicat">ApiCat</span>
                    </router-link>
                </div>
                <h2 class="text-xl text-zinc-800 font-medium mb-4">欢迎使用</h2>

                <el-form
                    @keyup.enter="onRegistBtnClick"
                    @submit.prevent="onRegistBtnClick"
                    label-position="top"
                    size="large"
                    :rules="rules"
                    ref="authForm"
                    :model="form"
                >
                    <el-form-item label="" prop="email">
                        <div class="ac-login__label">
                            <span>邮箱</span>
                        </div>
                        <el-input v-model="form.email" placeholder="请输入邮箱地址" autocomplete="on" />
                    </el-form-item>

                    <el-form-item label="" prop="password">
                        <div class="ac-login__label">
                            <span>密码</span>
                        </div>
                        <el-input type="password" v-model="form.password" placeholder="请输入密码" show-password autocomplete="on" />
                    </el-form-item>

                    <div class="">
                        <el-button :loading="userRegisterLoading" @click="onRegistBtnClick" class="w-full" type="primary">注&nbsp;册</el-button>
                    </div>
                </el-form>

                <p class="text-center mt-8">已有账号，去<router-link :to="LOGIN_PATH" class="text-blue-600 mx-1">登录</router-link></p>
            </div>
        </main>
    </div>
</template>
<script setup lang="ts">
    import { LOGIN_PATH } from '@/router/constant'
    import { reactive, ref } from 'vue'
    import { useApi } from '@/hooks/useApi'
    import { useUserStore } from '@/stores/user'
    import { isEmail } from '@ac/shared'

    const authForm = ref()

    const form = reactive({
        email: '',
        password: '',
    })

    const rules = {
        email: [
            { required: true, message: '请输入邮箱', trigger: 'blur' },
            {
                validator(rule: any, value: any, callback: any) {
                    if (!isEmail(value)) {
                        callback(new Error('请输入正确的邮箱'))
                    } else {
                        callback()
                    }
                },
                trigger: 'blur',
            },
        ],
        password: [{ required: true, message: '密码至少8位', min: 8, trigger: 'blur' }],
    } as any

    const [userRegisterLoading, userRegisterRequest] = useApi(useUserStore().register)

    const onRegistBtnClick = () => {
        authForm?.value
            .validate(async (valid: boolean) => {
                valid && (await userRegisterRequest(form))
            })
            .catch(() => {
                //
            })
    }
</script>
