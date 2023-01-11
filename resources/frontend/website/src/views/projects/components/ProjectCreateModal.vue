<template>
    <el-dialog v-model="isShow" :width="440" :close-on-click-modal="false" title="新建项目" append-to-body class="show-footer-line">
        <el-form
            ref="projectForm"
            :model="form"
            :rules="rules"
            label-position="top"
            style="width: 400px"
            class="py-2 pl-2"
            @submit.prevent="handleSubmit('projectForm')"
        >
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

            <ImageCorp :img-url="imgUrl" :file="file" @on-ok="onProjectIconUploadSuccess" :handle-upload="handleUpload" />
        </el-form>

        <template v-slot:footer>
            <el-button @click="onCloseBtnClick()"> 取消 </el-button>
            <el-button :loading="isLoading" type="primary" @click="handleSubmit('projectForm')"> 确 定 </el-button>
        </template>
    </el-dialog>
</template>

<script>
    import { $emit } from '@natosoft/shared'
    import { DOCUMENT_TYPES, DOCUMENT_TYPE } from '@/common/constant'

    import { PROJECT_DEFAULT_ICON } from '@/common/constant'
    import { createProject, uploadProjectIcon } from '@/api/project'
    import { Lock, View, Upload } from '@element-plus/icons-vue'
    import { useProjectsStore } from '@/stores/projects'
    import { toRefs, reactive } from 'vue'
    import { mapState } from 'pinia'

    export default {
        setup() {
            const state = reactive({
                isShow: false,
                PROJECT_DEFAULT_ICON,
                isLoading: false,
                imgUrl: null,
                file: null,
                disableType: false,
                docTypes: DOCUMENT_TYPES,
                form: {
                    icon_link: PROJECT_DEFAULT_ICON,
                    name: '',
                    type: DOCUMENT_TYPE.API.value,
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

            return {
                ...toRefs(state),
            }
        },

        computed: {
            ...mapState(useProjectsStore, ['activeGroup']),
        },

        components: {
            Lock,
            View,
            Upload,
        },

        watch: {
            isShow: function () {
                !this.isShow && this.reset()
            },
        },
        methods: {
            handleUpload(data) {
                data.append('icon', this.file)
                return uploadProjectIcon(data)
            },
            show() {
                this.isShow = true
            },

            hide() {
                this.isShow = false
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

            onCloseBtnClick() {
                this.isShow = false
                this.reset()
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

                        this.form.group_id = this.activeGroup.id

                        createProject(this.form)
                            .then((res) => {
                                this.$message.success(res.msg || '项目创建成功！')
                                $emit(this, 'on-ok', res.data || {})
                            })
                            .catch((e) => {
                                //
                            })
                            .finally(() => {
                                this.isLoading = false
                                this.hide()
                            })
                    }
                })
            },
        },
        emits: ['on-ok'],
    }
</script>
