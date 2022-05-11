<template>
    <main class="ac-verification">
        <div class="ac-verification__main">
            <span class="logo large">
                <img src="@/assets/image/logo.svg" alt="ApiCat" />
                <span class="logo-text logo-apicat mt-0">ApiCat</span>
            </span>

            <el-input class="my-7 w-1/2" v-model="form.secret_key" placeholder="访问密码" maxlength="6" clearable />

            <el-button :loading="isLoading" type="primary" @click="onSubmitBtnClick">继续访问</el-button>

            <img src="@/assets/image/img_join.png" class="mt-5 w-full" />
        </div>
    </main>
</template>

<script>
    import { checkProjectSecretKey } from '@/api/preview'
    import { Storage } from '@ac/shared'
    import { inject } from 'vue'

    export default {
        name: 'ProjectVerification',
        inject: ['showHeader'],
        data() {
            return {
                isLoading: false,
                form: {
                    project_id: this.$route.params.project_id || '',
                    secret_key: '',
                },
            }
        },
        methods: {
            onSubmitBtnClick() {
                this.isLoading = true
                checkProjectSecretKey(this.form)
                    .then((res) => {
                        Storage.set(Storage.KEYS.SECRET_PROJECT_TOKEN + this.form.project_id, res.data || '', true)
                        location.href = `/app/${this.form.project_id}`
                    })
                    .finally(() => {
                        this.isLoading = false
                    })
            },
        },

        setup() {
            const showHeader = inject('showHeader')
            showHeader(false)
        },
    }
</script>
