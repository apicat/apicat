<template>
    <el-card shadow="never">
        <template #header>
            <span>个人信息</span>
        </template>

        <el-form @submit.prevent="onSubmit" @keyup.enter="onSubmit" ref="formRef" :model="userInfo" :rules="rules" label-position="top" class="w-96">
            <el-form-item label="">
                <el-avatar :size="60" class="mr-6">
                    <img v-if="userInfo.avatar" :src="userInfo.avatar" />
                    <span v-else>{{ lastName }}</span>
                </el-avatar>
                <ImagePreview @done="onImageUpload">
                    <el-button :icon="UploadFilled">选择新头像</el-button>
                </ImagePreview>
            </el-form-item>

            <el-form-item label="用户名" prop="name">
                <el-input v-model="userInfo.name" maxlength="30" placeholder="用户名"></el-input>
            </el-form-item>

            <el-form-item label="邮箱" prop="email">
                <el-input v-model="userInfo.email" maxlength="150" placeholder="邮箱"></el-input>
            </el-form-item>

            <el-form-item>
                <el-button type="primary" @click="onSubmit" :loading="isLoading">保存</el-button>
            </el-form-item>
        </el-form>
    </el-card>

    <ImageCorp :imgUrl="imgUrl" :file="file" @on-ok="onAvatarUpdateSuccess" />
</template>

<script setup lang="ts">
    import { ref } from 'vue'
    import { UploadFilled } from '@element-plus/icons-vue'
    import { useUserStore } from '@/stores/user'
    import { storeToRefs } from 'pinia'
    import useApi from '@/hooks/useApi'

    const imgUrl = ref<string>('')
    const file = ref<File | null>(null)
    const formRef = ref()

    const userStore = useUserStore()

    const userInfo = ref({ ...userStore.userInfo })
    const { lastName } = storeToRefs(userStore)

    const rules = {
        name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        email: [{ required: true, message: '请输入邮箱', type: 'email', trigger: 'blur' }],
    } as any

    const [isLoading, execute] = useApi(userStore.updateUserInfo, { msg: '修改成功' })

    const onSubmit = () => {
        formRef.value?.validate((valid: boolean) => {
            valid && execute({ name: userInfo.value.name, email: userInfo.value.email, avatar: userInfo.value.avatar })
        })
    }

    const onAvatarUpdateSuccess = (url: string) => {
        userInfo.value.avatar = url
    }

    const onImageUpload = (url: string, f: File) => {
        imgUrl.value = url
        file.value = f
    }
</script>
