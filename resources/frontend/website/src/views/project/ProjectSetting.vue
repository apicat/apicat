<template>
    <el-card shadow="never">
        <template #header>
            <span>项目设置</span>
        </template>

        <el-form ref="projectForm" :model="form" :rules="rules" label-position="top" class="py-2 pl-2 max-w-sm" @submit.prevent="handleSubmit('projectForm')">
            <el-form-item label="">
                <el-avatar class="mr-7" shape="square" :size="60" :src="form.icon_link ? form.icon_link : PROJECT_DEFAULT_ICON" />
                <ImagePreview @done="onImageUpload">
                    <el-button>
                        <el-icon><Upload /></el-icon>上传新图标
                    </el-button>
                </ImagePreview>
            </el-form-item>

            <el-form-item label="项目名称" prop="name" class="hide_required">
                <el-input v-model="form.name" placeholder="项目名称" maxlength="255" />
            </el-form-item>

            <el-form-item label="权限">
                <el-radio-group v-model="form.visibility">
                    <el-radio-button :label="0">
                        <el-icon class="mr-1"><Lock /></el-icon>私有
                    </el-radio-button>
                    <el-radio-button :label="1">
                        <el-icon class="mr-1"><View /></el-icon>公开
                    </el-radio-button>
                </el-radio-group>
            </el-form-item>

            <el-form-item label="项目描述">
                <el-input v-model="form.description" type="textarea" :autosize="{ minRows: 4, maxRows: 4 }" maxlength="255" />
            </el-form-item>

            <el-button @click="handleSubmit('projectForm')" :loading="isLoading">保存</el-button>
        </el-form>
    </el-card>

    <ImageCorp :img-url="imgUrl" :file="file" @on-ok="onProjectIconUploadSuccess" :handle-upload="handleUpload" />
</template>

<script>
    import { DOCUMENT_TYPE } from '@/common/constant'
    import { uploadProjectIcon } from '@/api/project'
    import { Lock, View, Upload } from '@element-plus/icons-vue'
    import PROJECT_DEFAULT_ICON from '@/assets/image/icon-project.png'
    import { useProjectStore } from '@/stores/project'
    import { toRefs, reactive, watch } from 'vue'
    import { mapActions } from 'pinia'

    export default {
        setup() {
            const store = useProjectStore()

            const state = reactive({
                PROJECT_DEFAULT_ICON,
                isLoading: false,
                imgUrl: null,
                file: null,
                disableType: false,
                form: {
                    icon_link: PROJECT_DEFAULT_ICON,
                    name: '',
                    visibility: 0,
                    description: '',
                },
                rules: {
                    name: [
                        { required: true, message: '请输入项目名称', trigger: 'blur' },
                        { min: 2, message: '项目名称不能少于两个字', trigger: 'blur' },
                    ],
                },
            })

            watch(
                () => store.projectInfo,
                () => {
                    const info = store.projectInfo || {}
                    state.form = { ...info, icon_link: info.icon || PROJECT_DEFAULT_ICON }
                },
                {
                    immediate: true,
                }
            )

            return {
                ...toRefs(state),
            }
        },

        components: {
            Lock,
            View,
            Upload,
        },

        methods: {
            ...mapActions(useProjectStore, ['updateProjectInfo']),
            handleUpload(data) {
                data.append('icon', this.file)
                return uploadProjectIcon(data)
            },

            reset() {
                this.$refs['projectForm'].resetFields()
                this.form = {
                    icon_link: PROJECT_DEFAULT_ICON,
                    name: '',
                    type: DOCUMENT_TYPE.API.value,
                    visibility: 0,
                    description: '',
                }
            },

            onImageUpload(url, file) {
                this.imgUrl = url
                this.file = file
            },

            onProjectIconUploadSuccess(url) {
                this.form.icon_link = url
            },

            handleSubmit(name) {
                this.$refs[name].validate((valid) => {
                    if (valid) {
                        this.isLoading = true

                        const data = {}
                        data.project_id = this.form.id
                        data.icon_link = this.form.icon_link
                        data.icon = this.form.icon_link
                        data.name = this.form.name
                        data.visibility = this.form.visibility
                        data.description = this.form.description

                        this.updateProjectInfo(data)
                            .then((res) => {
                                this.$message.success(res.msg || '项目修改成功！')
                            })
                            .finally(() => {
                                this.isLoading = false
                            })
                    }
                })
            },
        },
        emits: ['on-ok'],
    }
</script>
