<template>
  <div class="flex min-h-screen">
    <main class="p-2 ac-login">
      <div class="shadow-xl ac-login__box">
        <div class="ac-login__logo">
          <router-link class="inline-flex items-center" to="/"> <AcLogo /> </router-link>
        </div>
        <h2 class="mb-4 text-xl font-medium text-zinc-800">欢迎使用</h2>

        <el-form
          @keyup.enter="onRegistBtnClick(authForm)"
          @submit.prevent="onRegistBtnClick(authForm)"
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

          <div class="mt-7">
            <el-button :loading="isLoading" @click="onRegistBtnClick(authForm)" class="w-full" type="primary">注&nbsp;册</el-button>
          </div>
        </el-form>

        <p class="mt-8 text-center"><router-link :to="LOGIN_PATH" class="mx-1 text-blue-600">登录</router-link></p>
      </div>
    </main>
  </div>
</template>
<script setup lang="ts">
import { LOGIN_PATH } from '@/router'
import { useApi } from '@/hooks/useApi'
import { useUserStore } from '@/store/user'
import { isEmail } from '@apicat/shared'
import { FormInstance } from 'element-plus'

const authForm = shallowRef()

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

const [isLoading, userRegisterRequest] = useApi(useUserStore().register)

const onRegistBtnClick = async (formIns: FormInstance) => {
  try {
    await formIns.validate()
    await userRegisterRequest(toRaw(form))
  } catch (error) {
    //
  }
}
</script>
