<template>
  <el-form @submit.prevent="onSubmit" @keyup.enter="onSubmit" ref="formRef" :model="userInfo" :rules="rules" label-position="top" class="max-w-sm">
    <el-form-item label="用户名" prop="username">
      <el-input v-model="userInfo.username" maxlength="30" placeholder="用户名"></el-input>
    </el-form-item>

    <el-form-item label="邮箱" prop="email">
      <el-input v-model="userInfo.email" maxlength="30" placeholder="邮箱"></el-input>
    </el-form-item>

    <el-form-item>
      <el-button type="primary" @click="onSubmit" :loading="isLoading">保存</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/store/user'
import useApi from '@/hooks/useApi'
import { UserInfo } from '@/typings/user'
import { isEmail } from '@apicat/shared'

const formRef = ref()

const userStore = useUserStore()

const userInfo: Ref<UserInfo> = ref({ ...userStore.userInfo })

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
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
} as any

const [isLoading, modifyUserInfoRequest] = useApi(userStore.modifyUserInfo)

const onSubmit = async () => {
  try {
    const valid = await formRef.value?.validate()
    valid && modifyUserInfoRequest(toRaw(userInfo))
  } catch (error) {
    //
  }
}
</script>
