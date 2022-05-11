<template>
    <div class="flex min-h-screen">
        <main class="ac-login p-2">
            <div class="ac-login__box shadow-xl">
                <div class="-ml-1 ac-login__logo">
                    <router-link class="inline-flex items-center" to="/">
                        <img src="@/assets/image/logo.svg" alt="ApiCat" /><span class="logo-text logo-apicat">ApiCat</span>
                    </router-link>
                </div>
                <h2 class="text-xl text-zinc-800 font-medium mb-4">登录</h2>

                <el-form
                    label-position="top"
                    size="large"
                    :rules="rules"
                    ref="authForm"
                    :model="form"
                    @keyup.enter="onLoginBtnClick"
                    @submit.prevent="onLoginBtnClick"
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
                        <el-input type="password" v-model="form.password" placeholder="请输入密码" autocomplete="on" show-password />
                    </el-form-item>

                    <div class="mt-7">
                        <el-button class="w-full" type="primary" :loading="emailLoginLoading" @click="onLoginBtnClick">登&nbsp;录</el-button>
                    </div>
                </el-form>

                <p class="text-center mt-8">还没有账号？马上<router-link :to="REGISTE_PATH" class="text-blue-600 mx-1">注册</router-link>吧</p>
            </div>
        </main>
    </div>
</template>
<script setup lang="ts">
    import { REGISTE_PATH } from '@/router/constant'
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

    const [emailLoginLoading, emailLoginRequest] = useApi(useUserStore().login)

    const onLoginBtnClick = () => {
        authForm?.value
            .validate(async (valid: boolean) => {
                valid && (await emailLoginRequest(form))
            })
            .catch(() => {
                //
            })
    }
</script>
