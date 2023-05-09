<template>
  <el-form @submit.prevent="onSubmit" @keyup.enter="onSubmit" ref="formRef" :model="form" :rules="rules" label-position="top" class="max-w-sm">
    <el-form-item label="旧密码" prop="password">
      <el-input type="password" v-model="form.password" placeholder="旧密码" autocomplete="on" />
    </el-form-item>

    <el-form-item label="新密码" prop="new_password">
      <el-input type="password" v-model="form.new_password" placeholder="新密码" autocomplete="on" />
    </el-form-item>

    <el-form-item label="确认新密码" prop="confirm_new_password">
      <el-input type="password" v-model="form.confirm_new_password" placeholder="新密码" autocomplete="on" />
    </el-form-item>

    <el-form-item>
      <el-button type="primary" @click="onSubmit" :loading="isLoading">{{ $t('app.common.save') }}</el-button>
    </el-form-item>
  </el-form>
</template>
<script setup lang="ts">
import useApi from '@/hooks/useApi'
import { useUserStore } from '@/store/user'
import { reactive, ref } from 'vue'

const userStore = useUserStore()
const formRef = ref()

const form = reactive({
  password: '',
  new_password: '',
  confirm_new_password: '',
})

const minLengthRule = { type: 'string', min: 8, message: '密码至少8位', trigger: 'blur' }

const rules = {
  password: [{ required: true, message: '请输入旧密码', trigger: 'blur' }, minLengthRule],
  new_password: [{ required: true, message: '请输入新密码', trigger: 'blur' }, minLengthRule],
  confirm_new_password: [
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

const [isLoading, modifyPasswordRequest] = useApi(userStore.modifyUserPassword)

const onSubmit = async () => {
  try {
    const valid = await formRef?.value.validate()
    valid && (await modifyPasswordRequest(form))
  } catch (error) {
    //
  }
}
</script>
