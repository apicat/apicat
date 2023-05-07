<template>
  <div class="flex min-h-screen">
    <main class="p-2 ac-login">
      <div class="shadow-xl ac-login__box">
        <div class="ac-login__logo">
          <router-link class="inline-flex items-center" to="/"> <AcLogo /> </router-link>
        </div>
        <h2 class="mb-4 text-xl font-medium text-zinc-800">登录</h2>

        <el-form label-position="top" size="large" :rules="rules" ref="authForm" :model="form" @keyup.enter="onLoginBtnClick(authForm)" @submit.prevent="onLoginBtnClick(authForm)">
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
            <el-button :loading="isLoading" @click="onLoginBtnClick(authForm)" class="w-full" type="primary">登&nbsp;录</el-button>
          </div>
        </el-form>

        <p class="mt-8 text-center"><router-link :to="REGISTER_PATH" class="mx-1 text-blue-600">注册账号</router-link></p>
      </div>
    </main>
  </div>
</template>
<script setup lang="ts">
import { REGISTER_PATH } from '@/router'
import { useApi } from '@/hooks/useApi'
import { useUserStore } from '@/store/user'
import { isEmail } from '@apicat/shared'
import { FormInstance } from 'element-plus'
import { UserInfo } from '@/typings/user'

const authForm = shallowRef()

const form: UserInfo = reactive({
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

const [isLoading, loginRequest] = useApi(useUserStore().login)

const onLoginBtnClick = async (formIns: FormInstance) => {
  try {
    await formIns.validate()
    await loginRequest(toRaw(form))
  } catch (error) {
    //
  }
}
</script>
